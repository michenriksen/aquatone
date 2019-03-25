package agents

import (
	"github.com/bberastegui/aquatone/core"
)

type URLLogger struct {
	session *core.Session
}

func NewURLLogger() *URLLogger {
	return &URLLogger{}
}

func (d *URLLogger) ID() string {
	return "agent:url_logger"
}

func (a *URLLogger) Register(s *core.Session) error {
	s.EventBus.SubscribeAsync(core.URLResponsive, a.OnURLResponsive, false)
	a.session = s
	return nil
}

func (a *URLLogger) OnURLResponsive(url string) {
	a.session.Out.Debug("[%s] Received new url: %s\n", a.ID(), url)
	a.session.AddResponsiveURL(url)
}
