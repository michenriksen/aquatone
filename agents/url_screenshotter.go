package agents

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/michenriksen/aquatone/core"
)

type URLScreenshotter struct {
	session    *core.Session
	chromePath string
}

func NewURLScreenshotter() *URLScreenshotter {
	return &URLScreenshotter{}
}

func (a *URLScreenshotter) ID() string {
	return "agent:url_screenshotter"
}

func (a *URLScreenshotter) Register(s *core.Session) error {
	s.EventBus.SubscribeAsync(core.URLResponsive, a.OnURLResponsive, false)
	a.session = s
	a.locateChrome()

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
	} else {
		out, err := exec.Command(a.chromePath, "--version").Output()
		if err != nil {
			a.session.Out.Warn("An error occurred while trying to determine version of Chromium.\n\n")
			return
		}
		version := string(out)
		re := regexp.MustCompile(`(\d+)\.`)
		match := re.FindStringSubmatch(version)
		if len(match) <= 0 {
			a.session.Out.Warn("Unable to determine version of Chromium. Screenshotting might be unreliable.\n\n")
			return
		}
		majorVersion, _ := strconv.Atoi(match[1])
		if majorVersion < 72 {
			a.session.Out.Warn("An older version of Chromium is installed. Screenshotting of HTTPS URLs might be unreliable.\n\n")
		}
	}

	a.session.Out.Debug("[%s] Located Chrome/Chromium binary at %s\n", a.ID(), a.chromePath)
}

func (a *URLScreenshotter) screenshotURL(s string) {
	filePath := a.session.GetFilePath(fmt.Sprintf("screenshots/%s.png", BaseFilenameFromURL(s)))
	var chromeArguments = []string{
		"--headless", "--disable-gpu", "--hide-scrollbars", "--mute-audio", "--disable-notifications",
		"--disable-crash-reporter",
		"--ignore-certificate-errors",
		"--user-agent=" + RandomUserAgent(),
		"--window-size=" + *a.session.Options.Resolution,
		"--screenshot=" + filePath,
	}

	if os.Geteuid() == 0 {
		chromeArguments = append(chromeArguments, "--no-sandbox")
	}

	if *a.session.Options.Proxy != "" {
		chromeArguments = append(chromeArguments, "--proxy-server="+*a.session.Options.Proxy)
	}

	chromeArguments = append(chromeArguments, s)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*a.session.Options.ScreenshotTimeout)*time.Millisecond)
	defer cancel()

	cmd := exec.CommandContext(ctx, a.chromePath, chromeArguments...)
	if err := cmd.Start(); err != nil {
		a.session.Out.Debug("[%s] Error: %v\n", a.ID(), err)
		a.session.Stats.IncrementScreenshotFailed()
		a.session.Out.Error("%s: screenshot failed: %s\n", s, err)
		a.killChromeProcessIfRunning(cmd)
		return
	}

	if err := cmd.Wait(); err != nil {
		a.session.Stats.IncrementScreenshotFailed()
		a.session.Out.Debug("[%s] Error: %v\n", a.ID(), err)
		if ctx.Err() == context.DeadlineExceeded {
			a.session.Out.Error("%s: screenshot timed out\n", s)
			a.killChromeProcessIfRunning(cmd)
			return
		}

		a.session.Out.Error("%s: screenshot failed: %s\n", s, err)
		a.killChromeProcessIfRunning(cmd)
		return
	}

	a.session.Stats.IncrementScreenshotSuccessful()
	a.session.Out.Info("%s: %s\n", s, Green("screenshot successful"))
	a.killChromeProcessIfRunning(cmd)
}

func (a *URLScreenshotter) killChromeProcessIfRunning(cmd *exec.Cmd) {
	if cmd.Process == nil {
		return
	}
	cmd.Process.Release()
	cmd.Process.Kill()
}
