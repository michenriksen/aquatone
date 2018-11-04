package agents

import (
	"github.com/michenriksen/aquatone/core"
)

type URLPublisher struct {
	session *core.Session
}

func NewURLPublisher() *URLPublisher {
	return &URLPublisher{}
}

func (d *URLPublisher) ID() string {
	return "agent:url_publisher"
}

func (a *URLPublisher) Register(s *core.Session) error {
	s.EventBus.SubscribeAsync(core.TCPPort, a.OnTCPPort, false)
	a.session = s
	return nil
}

func (a *URLPublisher) OnTCPPort(port int, host string) {
	a.session.Out.Debug("[%s] Received new open port on %s: %d\n", a.ID(), host, port)
	url := HostAndPortToURL(host, port, "")
	a.session.EventBus.Publish(core.URL, url)
}
