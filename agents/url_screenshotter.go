package agents

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/michenriksen/aquatone/core"
)

type URLScreenshotter struct {
	session         *core.Session
	chromePath      string
	tempUserDirPath string
}

func NewURLScreenshotter() *URLScreenshotter {
	return &URLScreenshotter{}
}

func (a *URLScreenshotter) ID() string {
	return "agent:url_screenshotter"
}

func (a *URLScreenshotter) Register(s *core.Session) error {
	s.EventBus.SubscribeAsync(core.URLResponsive, a.OnURLResponsive, false)
	s.EventBus.SubscribeAsync(core.SessionEnd, a.OnSessionEnd, false)
	a.session = s
	a.createTempUserDir()

	return nil
}

func (a *URLScreenshotter) OnURLResponsive(url string) {
	a.session.Out.Debug("[%s] Received new responsive URL %s\n", a.ID(), url)
	page := a.session.GetPage(url)
	if page == nil {
		a.session.Out.Error("Unable to find page for URL: %s\n", url)
		return
	}

	a.session.WaitGroup.Add()
	go func(page *core.Page) {
		defer a.session.WaitGroup.Done()
		a.screenshotPage(page)
	}(page)
}

func (a *URLScreenshotter) OnSessionEnd() {
	a.session.Out.Debug("[%s] Received SessionEnd event\n", a.ID())
	os.RemoveAll(a.tempUserDirPath)
	a.session.Out.Debug("[%s] Deleted temporary user directory at: %s\n", a.ID(), a.tempUserDirPath)
}

func (a *URLScreenshotter) createTempUserDir() {
	dir, err := ioutil.TempDir("", "aquatone-chrome")
	if err != nil {
		a.session.Out.Fatal("Unable to create temporary user directory for Chrome/Chromium browser\n")
		os.Exit(1)
	}
	a.session.Out.Debug("[%s] Created temporary user directory at: %s\n", a.ID(), dir)
	a.tempUserDirPath = dir
}

func (a URLScreenshotter) getOpts() (options []chromedp.ExecAllocatorOption) {
	if *a.session.Options.Proxy != "" {
		options = append(options, chromedp.ProxyServer(*a.session.Options.Proxy))
	}

	if *a.session.Options.ChromePath != "" {
		options = append(options, chromedp.ExecPath(*a.session.Options.ChromePath))
	}

	return
}

// execAllocator turns a.getOpts() (the chrome instance allocator options) into a derivative context.Context
func (a URLScreenshotter) execAllocator(parent context.Context) (context.Context, context.CancelFunc) {
	return chromedp.NewExecAllocator(parent, a.getOpts()...)
}

func (a *URLScreenshotter) screenshotPage(page *core.Page) {
	filePath := fmt.Sprintf("screenshots/%s.png", page.BaseFilename())

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*a.session.Options.ScreenshotTimeout)*time.Millisecond)
	ctx, cancel = a.execAllocator(ctx)
	ctx, cancel = chromedp.NewContext(ctx)

	defer cancel()

	var pic []byte
	if err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(page.URL),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Screenshot("body", &pic, chromedp.NodeVisible, chromedp.ByQuery),
	}); err != nil {
		a.session.Out.Debug("%s Error: %v\n", a.ID, err)
		a.session.Stats.IncrementScreenshotFailed()
		a.session.Out.Error("%s: screenshot failed: %s\n", page.URL, err)
		return
	}

	if err := ioutil.WriteFile(filePath, pic, 0700); err != nil {
		a.session.Out.Debug("%s Error: %v\n", a.ID(), err)
		a.session.Stats.IncrementScreenshotFailed()
		a.session.Out.Error("%s: screenshot failed: %s\n", page.URL, err)
		return
	}

	a.session.Stats.IncrementScreenshotSuccessful()
	a.session.Out.Info("%s: %s\n", page.URL, Green("screenshot successful"))
	page.ScreenshotPath = filePath
	page.HasScreenshot = true
}

func (a *URLScreenshotter) killChromeProcessIfRunning(cmd *exec.Cmd) {
	if cmd.Process == nil {
		return
	}
	cmd.Process.Release()
	cmd.Process.Kill()
}
