package agents

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

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
	var url string
	if a.isTLS(port, host) {
		url = HostAndPortToURL(host, port, "https")
	} else {
		url = HostAndPortToURL(host, port, "http")
	}
	a.session.EventBus.Publish(core.URL, url)
}

func (a *URLPublisher) isTLS(port int, host string) bool {
	if port == 80 {
		return false
	}

	if port == 443 {
		return true
	}

	dialer := &net.Dialer{Timeout: time.Duration(*a.session.Options.HTTPTimeout) * time.Millisecond}
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := tls.DialWithDialer(dialer, "tcp", fmt.Sprintf("%s:%d", host, port), conf)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
