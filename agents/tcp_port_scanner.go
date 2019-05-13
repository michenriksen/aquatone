package agents

import (
	"fmt"
	"net"
	"time"

	"github.com/michenriksen/aquatone/core"
)

type TCPPortScanner struct {
	session *core.Session
}

func NewTCPPortScanner() *TCPPortScanner {
	return &TCPPortScanner{}
}

func (d *TCPPortScanner) ID() string {
	return "agent:tcp_port_scanner"
}

func (a *TCPPortScanner) Register(s *core.Session) error {
	s.EventBus.SubscribeAsync(core.Host, a.OnHost, false)
	a.session = s
	return nil
}

func (a *TCPPortScanner) OnHost(host string) {
	a.session.Out.Debug("[%s] Received new host: %s\n", a.ID(), host)
	for _, port := range a.session.Ports {
		a.session.WaitGroup.Add()
		go func(port int, host string) {
			defer a.session.WaitGroup.Done()
			if a.scanPort(port, host) {
				a.session.Stats.IncrementPortOpen()
				a.session.Out.Info("%s: port %s %s\n", host, Green(fmt.Sprintf("%d", port)), Green("open"))
				a.session.EventBus.Publish(core.TCPPort, port, host)
			} else {
				a.session.Stats.IncrementPortClosed()
				a.session.Out.Debug("[%s] Port %d is closed on %s\n", a.ID(), port, host)
			}
		}(port, host)
	}
}

func (a *TCPPortScanner) scanPort(port int, host string) bool {
	conn, _ := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), time.Duration(*a.session.Options.ScanTimeout)*time.Millisecond)
	if conn != nil {
		conn.Close()
		return true
	}
	return false
}
