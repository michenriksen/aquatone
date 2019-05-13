package agents

import (
	"fmt"
	"net"

	"github.com/michenriksen/aquatone/core"
)

type URLHostnameResolver struct {
	session *core.Session
}

func NewURLHostnameResolver() *URLHostnameResolver {
	return &URLHostnameResolver{}
}

func (a *URLHostnameResolver) ID() string {
	return "agent:url_hostname_resolver"
}

func (a *URLHostnameResolver) Register(s *core.Session) error {
	s.EventBus.SubscribeAsync(core.URLResponsive, a.OnURLResponsive, false)
	a.session = s

	return nil
}

func (a *URLHostnameResolver) OnURLResponsive(url string) {
	a.session.Out.Debug("[%s] Received new responsive URL %s\n", a.ID(), url)
	page := a.session.GetPage(url)
	if page == nil {
		a.session.Out.Error("Unable to find page for URL: %s\n", url)
		return
	}

	if page.IsIPHost() {
		a.session.Out.Debug("[%s] Skipping hostname resolving on IP host: %s\n", a.ID(), url)
		page.Addrs = []string{page.ParsedURL().Hostname()}
		return
	}

	a.session.WaitGroup.Add()
	go func(page *core.Page) {
		defer a.session.WaitGroup.Done()
		addrs, err := net.LookupHost(fmt.Sprintf("%s.", page.ParsedURL().Hostname()))
		if err != nil {
			a.session.Out.Debug("[%s] Error: %v\n", a.ID(), err)
			a.session.Out.Error("Failed to resolve hostname for %s\n", page.URL)
			return
		}

		page.Addrs = addrs
	}(page)
}
