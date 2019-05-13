package agents

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/michenriksen/aquatone/core"
	"github.com/parnurzeal/gorequest"
)

type URLRequester struct {
	session *core.Session
}

func NewURLRequester() *URLRequester {
	return &URLRequester{}
}

func (d *URLRequester) ID() string {
	return "agent:url_requester"
}

func (a *URLRequester) Register(s *core.Session) error {
	s.EventBus.SubscribeAsync(core.URL, a.OnURL, false)
	a.session = s
	return nil
}

func (a *URLRequester) OnURL(url string) {
	a.session.Out.Debug("[%s] Received new URL %s\n", a.ID(), url)
	a.session.WaitGroup.Add()
	go func(url string) {
		defer a.session.WaitGroup.Done()
		http := Gorequest(a.session.Options)
		resp, _, errs := http.Get(url).
			Set("User-Agent", RandomUserAgent()).
			Set("X-Forwarded-For", RandomIPv4Address()).
			Set("Via", fmt.Sprintf("1.1 %s", RandomIPv4Address())).
			Set("Forwarded", fmt.Sprintf("for=%s;proto=http;by=%s", RandomIPv4Address(), RandomIPv4Address())).End()
		var status string
		if errs != nil {
			a.session.Stats.IncrementRequestFailed()
			for _, err := range errs {
				a.session.Out.Debug("[%s] Error: %v\n", a.ID(), err)
				if os.IsTimeout(err) {
					a.session.Out.Error("%s: request timeout\n", url)
					return
				}
			}
			a.session.Out.Debug("%s: failed\n", url)
			return
		}

		a.session.Stats.IncrementRequestSuccessful()
		if resp.StatusCode >= 500 {
			a.session.Stats.IncrementResponseCode5xx()
			status = Red(resp.Status)
		} else if resp.StatusCode >= 400 {
			a.session.Stats.IncrementResponseCode4xx()
			status = Yellow(resp.Status)
		} else if resp.StatusCode >= 300 {
			a.session.Stats.IncrementResponseCode3xx()
			status = Green(resp.Status)
		} else {
			a.session.Stats.IncrementResponseCode2xx()
			status = Green(resp.Status)
		}
		a.session.Out.Info("%s: %s\n", url, status)

		page, err := a.createPageFromResponse(url, resp)
		if err != nil {
			a.session.Out.Debug("[%s] Error: %v\n", a.ID(), err)
			a.session.Out.Error("Failed to create page for URL: %s\n", url)
			return
		}

		a.writeHeaders(page)
		if *a.session.Options.SaveBody {
			a.writeBody(page, resp)
		}

		a.session.EventBus.Publish(core.URLResponsive, url)
	}(url)
}

func (a *URLRequester) createPageFromResponse(url string, resp gorequest.Response) (*core.Page, error) {
	page, err := a.session.AddPage(url)
	if err != nil {
		return nil, err
	}

	page.Status = resp.Status
	for name, value := range resp.Header {
		page.AddHeader(name, strings.Join(value, " "))
	}

	return page, nil
}

func (a *URLRequester) writeHeaders(page *core.Page) {
	filepath := fmt.Sprintf("headers/%s.txt", page.BaseFilename())
	headers := fmt.Sprintf("%s\n", page.Status)
	for _, header := range page.Headers {
		headers += fmt.Sprintf("%v: %v\n", header.Name, header.Value)
	}
	if err := ioutil.WriteFile(a.session.GetFilePath(filepath), []byte(headers), 0644); err != nil {
		a.session.Out.Debug("[%s] Error: %v\n", a.ID(), err)
		a.session.Out.Error("Failed to write HTTP response headers for %s to %s\n", page.URL, a.session.GetFilePath(filepath))
	}
	page.HeadersPath = filepath
}

func (a *URLRequester) writeBody(page *core.Page, resp gorequest.Response) {
	filepath := fmt.Sprintf("html/%s.html", page.BaseFilename())
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		a.session.Out.Debug("[%s] Error: %v\n", a.ID(), err)
		a.session.Out.Error("Failed to read response body for %s\n", page.URL)
		return
	}

	if err := ioutil.WriteFile(a.session.GetFilePath(filepath), body, 0644); err != nil {
		a.session.Out.Debug("[%s] Error: %v\n", a.ID(), err)
		a.session.Out.Error("Failed to write HTTP response body for %s to %s\n", page.URL, a.session.GetFilePath(filepath))
	}
	page.BodyPath = filepath
}
