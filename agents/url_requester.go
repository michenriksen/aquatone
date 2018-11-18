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
		ip := RandomIPv4Address()
		if *a.session.Options.Waf {
			ip = "127.0.0.1"
		}
		resp, _, errs := http.Get(url).
			Set("User-Agent", RandomUserAgent()).
			Set("X-Client-IP", ip).
			Set("X-Remote-IP", ip).
			Set("X-Remote-Addr", ip).
			Set("X-Forwarded-For", ip).
			Set("X-OriginatingIP", ip).
			Set("Via", "1.1 "+ip).
			Set("Forwarded", "for="+ip+";proto=http;by="+ip).
			End()
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

		a.writeHeaders(url, resp)
		if *a.session.Options.SaveBody {
			a.writeBody(url, resp)
		}

		a.session.EventBus.Publish(core.URLResponsive, url)
	}(url)
}

func (a *URLRequester) writeHeaders(url string, resp gorequest.Response) {
	filepath := a.session.GetFilePath(fmt.Sprintf("headers/%s.txt", BaseFilenameFromURL(url)))
	headers := fmt.Sprintf("%s\n", resp.Status)
	for name, value := range resp.Header {
		headers += fmt.Sprintf("%v: %v\n", name, strings.Join(value, " "))
	}
	if err := ioutil.WriteFile(filepath, []byte(headers), 0644); err != nil {
		a.session.Out.Debug("[%s] Error: %v\n", a.ID(), err)
		a.session.Out.Error("Failed to write HTTP response headers for %s to %s\n", url, filepath)
	}
}

func (a *URLRequester) writeBody(url string, resp gorequest.Response) {
	filepath := a.session.GetFilePath(fmt.Sprintf("html/%s.html", BaseFilenameFromURL(url)))
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		a.session.Out.Debug("[%s] Error: %v\n", a.ID(), err)
		a.session.Out.Error("Failed to read response body for %s\n", url)
		return
	}

	if err := ioutil.WriteFile(filepath, body, 0644); err != nil {
		a.session.Out.Debug("[%s] Error: %v\n", a.ID(), err)
		a.session.Out.Error("Failed to write HTTP response body for %s to %s\n", url, filepath)
	}
}
