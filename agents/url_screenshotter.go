package agents

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/security"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/runner"
	"github.com/michenriksen/aquatone/core"
)

type URLScreenshotter struct {
	session    *core.Session
	chromePath string
	ctxt       context.Context
	pool       *chromedp.Pool
}

func NewURLScreenshotter() *URLScreenshotter {
	return &URLScreenshotter{}
}

func (d *URLScreenshotter) ID() string {
	return "agent:url_screenshotter"
}

func (a *URLScreenshotter) Register(s *core.Session) error {
	s.EventBus.SubscribeAsync(core.URLResponsive, a.OnURLResponsive, false)
	a.session = s
	a.locateChrome()
	a.ctxt, _ = context.WithCancel(context.Background())
	a.pool, _ = chromedp.NewPool()

	return nil
}

func (a *URLScreenshotter) OnURLResponsive(url string) {
	a.session.Out.Debug("[%s] Received new responsive URL %s\n", a.ID(), url)
	a.session.WaitGroup.Add()
	go func(url string) {
		defer a.session.WaitGroup.Done()
		a.screenshotURL(url)
	}(url)
}

func (a *URLScreenshotter) locateChrome() {
	if *a.session.Options.ChromePath != "" {
		a.chromePath = *a.session.Options.ChromePath
		return
	}

	paths := []string{
		"/usr/bin/google-chrome",
		"/usr/bin/google-chrome-beta",
		"/usr/bin/google-chrome-unstable",
		"/usr/bin/chromium-browser",
		"/usr/bin/chromium",
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		"/Applications/Google Chrome Canary.app/Contents/MacOS/Google Chrome Canary",
		"/Applications/Chromium.app/Contents/MacOS/Chromium",
		"C:/Program Files (x86)/Google/Chrome/Application/chrome.exe",
	}

	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}
		a.chromePath = path
	}

	if a.chromePath == "" {
		a.session.Out.Fatal("Unable to locate a valid installation of Chrome. Install Google Chrome or try specifying a valid location with the -chrome-path option.\n")
		os.Exit(1)
	}

	if strings.Contains(strings.ToLower(a.chromePath), "chrome") {
		a.session.Out.Warn("Using unreliable Google Chrome for screenshots. Install Chromium for better results.\n\n")
	}

	a.session.Out.Debug("[%s] Located Chrome/Chromium binary at %s\n", a.ID(), a.chromePath)
}

func (a *URLScreenshotter) screenshotURL(s string) {
	filePath := a.session.GetFilePath(fmt.Sprintf("screenshots/%s.png", BaseFilenameFromURL(s)))

	// allocate a chrome headless instance
	c, err := a.pool.Allocate(a.ctxt,
		runner.DisableGPU,
		runner.Headless,
		runner.ExecPath(a.chromePath),
		runner.Flag("ignore-certificate-errors", true),
		runner.Flag("disable-crash-reporter", true),
		runner.Flag("disable-notifications", true),
		runner.Flag("hide-scrollbars", true),
		runner.Flag("window-size", *a.session.Options.Resolution),
		runner.Flag("user-agent", RandomUserAgent()),
		runner.Flag("mute-audio", true),
		runner.Flag("incognito", true),
	)
	defer c.Release()

	// screenshot buffer
	var picbuf []byte

	// set headers
	ip := RandomIPv4Address()
	if *a.session.Options.Waf {
		ip = "127.0.0.1"
	}
	headers := map[string]interface{}{
		"X-Client-IP":     ip,
		"X-Remote-IP":     ip,
		"X-Remote-Addr":   ip,
		"X-Forwarded-For": ip,
		"X-OriginatingIP": ip,
		"Via":             "1.1 " + ip,
		"Forwarded":       "for=" + ip + ";proto=http;by=" + ip,
	}

	// create a new tab and snap it
	t := chromedp.Tasks{
		network.Enable(),
		network.SetExtraHTTPHeaders(network.Headers(headers)),
		security.SetIgnoreCertificateErrors(true),
		chromedp.Navigate(s),
		chromedp.Sleep(time.Second * 5),
		chromedp.ActionFunc(func(ctxt context.Context, h cdp.Executor) error {
			picbuf, err = page.CaptureScreenshot().Do(ctxt, h)
			if err != nil {
				a.session.Out.Debug("[%s] Error: %v\n", a.ID(), err)
				a.session.Stats.IncrementScreenshotFailed()
				a.session.Out.Error("%s: screenshot failed: %s\n", s, err)
				return err
			}
			return nil
		}),
	}

	// run tasks
	err = c.Run(a.ctxt, t)
	if err != nil {
		a.session.Out.Debug("[%s] Error: %v\n", a.ID(), err)
		a.session.Stats.IncrementScreenshotFailed()
		a.session.Out.Error("%s: screenshot failed: %s\n", s, err)
		return
	}

	// write to disk
	err = ioutil.WriteFile(filePath, picbuf, 0644)
	if err != nil {
		a.session.Out.Debug("[%s] Error: %v\n", a.ID(), err)
		a.session.Stats.IncrementScreenshotFailed()
		a.session.Out.Error("%s: screenshot failed: %s\n", s, err)
		return
	}

	/*
		if os.Geteuid() == 0 {
			chromeArguments = append(chromeArguments, "--no-sandbox")
		}

		if *a.session.Options.Proxy != "" {
			chromeArguments = append(chromeArguments, "--proxy-server="+*a.session.Options.Proxy)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*a.session.Options.ScreenshotTimeout)*time.Millisecond)
	*/

	a.session.Stats.IncrementScreenshotSuccessful()
	a.session.Out.Info("%s: %s\n", s, Green("screenshot successful"))
}
