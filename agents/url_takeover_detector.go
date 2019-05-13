package agents

import (
	"fmt"
	"net"
	"strings"

	"github.com/michenriksen/aquatone/core"
)

type URLTakeoverDetector struct {
	session *core.Session
}

func NewURLTakeoverDetector() *URLTakeoverDetector {
	return &URLTakeoverDetector{}
}

func (d *URLTakeoverDetector) ID() string {
	return "agent:url_takeover_detector"
}

func (a *URLTakeoverDetector) Register(s *core.Session) error {
	s.EventBus.SubscribeAsync(core.URLResponsive, a.OnURLResponsive, false)
	a.session = s
	return nil
}

func (a *URLTakeoverDetector) OnURLResponsive(u string) {
	a.session.Out.Debug("[%s] Received new url: %s\n", a.ID(), u)
	page := a.session.GetPage(u)
	if page == nil {
		a.session.Out.Error("Unable to find page for URL: %s\n", u)
		return
	}

	if page.IsIPHost() {
		a.session.Out.Debug("[%s] Skipping takeover detection on IP URL %s\n", a.ID(), u)
		return
	}

	a.session.WaitGroup.Add()
	go func(p *core.Page) {
		defer a.session.WaitGroup.Done()
		a.runDetectorFunctions(p)
	}(page)
}

func (a *URLTakeoverDetector) runDetectorFunctions(page *core.Page) {
	hostname := page.ParsedURL().Hostname()
	addrs, err := net.LookupHost(fmt.Sprintf("%s.", hostname))
	if err != nil {
		a.session.Out.Error("Unable to resolve %s to IP addresses: %s\n", hostname, err)
		return
	}
	cname, err := net.LookupCNAME(fmt.Sprintf("%s.", hostname))
	if err != nil {
		a.session.Out.Error("Unable to resolve %s to CNAME: %s\n", hostname, err)
		return
	}

	a.session.Out.Debug("[%s] IP addresses for %s: %v\n", a.ID(), hostname, addrs)
	a.session.Out.Debug("[%s] CNAME for %s: %s\n", a.ID(), hostname, cname)

	body, err := a.session.ReadFile(fmt.Sprintf("html/%s.html", page.BaseFilename()))
	if err != nil {
		a.session.Out.Debug("[%s] Error reading HTML body file for %s: %s\n", a.ID(), page.URL, err)
		return
	}

	if a.detectGithubPages(page, addrs, cname, string(body)) {
		return
	}

	if a.detectAmazonS3(page, addrs, cname, string(body)) {
		return
	}

	if a.detectCampaignMonitor(page, addrs, cname, string(body)) {
		return
	}

	if a.detectCargoCollective(page, addrs, cname, string(body)) {
		return
	}

	if a.detectFeedPress(page, addrs, cname, string(body)) {
		return
	}

	if a.detectGhost(page, addrs, cname, string(body)) {
		return
	}

	if a.detectHelpjuice(page, addrs, cname, string(body)) {
		return
	}

	if a.detectHelpScout(page, addrs, cname, string(body)) {
		return
	}

	if a.detectHeroku(page, addrs, cname, string(body)) {
		return
	}

	if a.detectJetBrains(page, addrs, cname, string(body)) {
		return
	}

	if a.detectMicrosoftAzure(page, addrs, cname, string(body)) {
		return
	}

	if a.detectReadme(page, addrs, cname, string(body)) {
		return
	}

	if a.detectSurge(page, addrs, cname, string(body)) {
		return
	}

	if a.detectTumblr(page, addrs, cname, string(body)) {
		return
	}

	if a.detectUserVoice(page, addrs, cname, string(body)) {
		return
	}

	if a.detectWordpress(page, addrs, cname, string(body)) {
		return
	}

	if a.detectSmugMug(page, addrs, cname, string(body)) {
		return
	}

	if a.detectStrikingly(page, addrs, cname, string(body)) {
		return
	}

	if a.detectUptimeRobot(page, addrs, cname, string(body)) {
		return
	}

	if a.detectPantheon(page, addrs, cname, string(body)) {
		return
	}
}

func (a *URLTakeoverDetector) detectGithubPages(p *core.Page, addrs []string, cname string, body string) bool {
	githubAddrs := [...]string{"185.199.108.153", "185.199.109.153", "185.199.110.153", "185.199.111.153"}
	fingerprints := [...]string{"There isn't a GitHub Pages site here.", "For root URLs (like http://example.com/) you must provide an index.html file"}
	for _, githubAddr := range githubAddrs {
		for _, addr := range addrs {
			if addr == githubAddr {
				for _, fingerprint := range fingerprints {
					if strings.Contains(body, fingerprint) {
						p.AddTag("Domain Takeover", "danger", "https://help.github.com/articles/using-a-custom-domain-with-github-pages/")
						a.session.Out.Warn("%s: vulnerable to takeover on Github Pages\n", p.URL)
						return true
					}
				}
				return true
			}
		}
	}
	return false
}

func (a *URLTakeoverDetector) detectAmazonS3(p *core.Page, addrs []string, cname string, body string) bool {
	fingerprints := [...]string{"NoSuchBucket", "The specified bucket does not exist"}
	if !strings.HasSuffix(cname, ".amazonaws.com.") {
		return false
	}
	for _, fingerprint := range fingerprints {
		if strings.Contains(body, fingerprint) {
			p.AddTag("Domain Takeover", "danger", "https://docs.aws.amazon.com/AmazonS3/latest/dev/website-hosting-custom-domain-walkthrough.html")
			a.session.Out.Warn("%s: vulnerable to takeover on Amazon S3\n", p.URL)
			return true
		}
	}
	return true
}

func (a *URLTakeoverDetector) detectCampaignMonitor(p *core.Page, addrs []string, cname string, body string) bool {
	if cname != "cname.createsend.com." {
		return false
	}
	p.AddTag("Campaign Monitor", "info", "https://www.campaignmonitor.com/")
	if strings.Contains(body, "Double check the URL or ") {
		p.AddTag("Domain Takeover", "danger", "https://help.campaignmonitor.com/custom-domain-names")
		a.session.Out.Warn("%s: vulnerable to takeover on Campaign Monitor\n", p.URL)
		return true
	}
	return true
}

func (a *URLTakeoverDetector) detectCargoCollective(p *core.Page, addrs []string, cname string, body string) bool {
	if cname != "subdomain.cargocollective.com." {
		return false
	}
	p.AddTag("Cargo Collective", "info", "https://cargocollective.com/")
	if strings.Contains(body, "404 Not Found") {
		p.AddTag("Domain Takeover", "danger", "https://support.2.cargocollective.com/Using-a-Third-Party-Domain")
		a.session.Out.Warn("%s: vulnerable to takeover on Cargo Collective\n", p.URL)
		return true
	}
	return true
}

func (a *URLTakeoverDetector) detectFeedPress(p *core.Page, addrs []string, cname string, body string) bool {
	if cname != "redirect.feedpress.me." {
		return false
	}
	p.AddTag("FeedPress", "info", "https://feed.press/")
	if strings.Contains(body, "The feed has not been found.") {
		p.AddTag("Domain Takeover", "danger", "https://support.feed.press/article/61-how-to-create-a-custom-hostname")
		a.session.Out.Warn("%s: vulnerable to takeover on FeedPress\n", p.URL)
		return true
	}
	return true
}

func (a *URLTakeoverDetector) detectGhost(p *core.Page, addrs []string, cname string, body string) bool {
	if !strings.HasSuffix(cname, ".ghost.io.") {
		return false
	}
	if strings.Contains(body, "The thing you were looking for is no longer here, or never was") {
		p.AddTag("Domain Takeover", "danger", "https://docs.ghost.org/faq/using-custom-domains/")
		a.session.Out.Warn("%s: vulnerable to takeover on Ghost\n", p.URL)
		return true
	}
	return true
}

func (a *URLTakeoverDetector) detectHelpjuice(p *core.Page, addrs []string, cname string, body string) bool {
	if !strings.HasSuffix(cname, ".helpjuice.com.") {
		return false
	}
	p.AddTag("Helpjuice", "info", "https://helpjuice.com/")
	if strings.Contains(body, "We could not find what you're looking for.") {
		p.AddTag("Domain Takeover", "danger", "https://help.helpjuice.com/34339-getting-started/custom-domain")
		a.session.Out.Warn("%s: vulnerable to takeover on Helpjuice\n", p.URL)
		return true
	}
	return false
}

func (a *URLTakeoverDetector) detectHelpScout(p *core.Page, addrs []string, cname string, body string) bool {
	if !strings.HasSuffix(cname, ".helpscoutdocs.com.") {
		return false
	}
	p.AddTag("HelpScout", "info", "https://www.helpscout.net/")
	if strings.Contains(body, "No settings were found for this company:") {
		p.AddTag("Domain Takeover", "danger", "https://docs.helpscout.net/article/42-setup-custom-domain")
		a.session.Out.Warn("%s: vulnerable to takeover on HelpScout\n", p.URL)
		return true
	}
	return true
}

func (a *URLTakeoverDetector) detectHeroku(p *core.Page, addrs []string, cname string, body string) bool {
	herokuCnames := [...]string{".herokudns.com.", ".herokuapp.com.", ".herokussl.com."}
	for _, herokuCname := range herokuCnames {
		if strings.HasSuffix(cname, herokuCname) {
			p.AddTag("Heroku", "info", "https://www.heroku.com/")
			if strings.Contains(body, "No such app") {
				p.AddTag("Domain Takeover", "danger", "https://devcenter.heroku.com/articles/custom-domains")
				a.session.Out.Warn("%s: vulnerable to takeover on Heroku\n", p.URL)
				return true
			}
			return true
		}
	}
	return false
}

func (a *URLTakeoverDetector) detectJetBrains(p *core.Page, addrs []string, cname string, body string) bool {
	if !strings.HasSuffix(cname, ".myjetbrains.com.") {
		return false
	}
	p.AddTag("JetBrains", "info", "https://www.jetbrains.com/")
	if strings.Contains(body, "is not a registered InCloud YouTrack") {
		p.AddTag("Domain Takeover", "danger", "https://www.jetbrains.com/help/youtrack/incloud/Domain-Settings.html#use-custom-domain-name")
		a.session.Out.Warn("%s: vulnerable to takeover on JetBrains\n", p.URL)
		return true
	}
	return true
}

func (a *URLTakeoverDetector) detectMicrosoftAzure(p *core.Page, addrs []string, cname string, body string) bool {
	if !strings.HasSuffix(cname, ".azurewebsites.net.") {
		return false
	}
	p.AddTag("Microsoft Azure", "info", "https://azure.microsoft.com/")
	if strings.Contains(body, "404 Web Site not found") {
		p.AddTag("Domain Takeover", "danger", "https://docs.microsoft.com/en-us/azure/app-service/app-service-web-tutorial-custom-domain")
		a.session.Out.Warn("%s: vulnerable to takeover on Microsoft Azure\n", p.URL)
		return true
	}
	return true
}

func (a *URLTakeoverDetector) detectReadme(p *core.Page, addrs []string, cname string, body string) bool {
	readmeCnames := [...]string{".readme.io.", ".readmessl.com."}
	for _, readmeCname := range readmeCnames {
		if strings.HasSuffix(cname, readmeCname) {
			p.AddTag("Readme", "info", "https://readme.io/")
			if strings.Contains(body, "Project doesnt exist... yet!") {
				p.AddTag("Domain Takeover", "danger", "https://readme.readme.io/docs/setting-up-custom-domain")
				a.session.Out.Warn("%s: vulnerable to takeover on Readme\n", p.URL)
				return true
			}
			return true
		}
	}
	return false
}

func (a *URLTakeoverDetector) detectSurge(p *core.Page, addrs []string, cname string, body string) bool {
	detected := false
	for _, addr := range addrs {
		if addr == "45.55.110.124" {
			detected = true
			break
		}
	}
	if cname == "na-west1.surge.sh." {
		detected = true
	}
	if detected {
		p.AddTag("Surge", "info", "https://surge.sh/")
		if strings.Contains(body, "project not found") {
			p.AddTag("Domain Takeover", "danger", "https://surge.sh/help/adding-a-custom-domain")
			a.session.Out.Warn("%s: vulnerable to takeover on Surge\n", p.URL)
		}
		return true
	}
	return false
}

func (a *URLTakeoverDetector) detectTumblr(p *core.Page, addrs []string, cname string, body string) bool {
	detected := false
	for _, addr := range addrs {
		if addr == "66.6.44.4" {
			detected = true
			break
		}
	}
	if cname == "domains.tumblr.com." {
		detected = true
	}
	if detected {
		if strings.Contains(body, "Whatever you were looking for doesn't currently exist at this address") {
			p.AddTag("Domain Takeover", "danger", "https://tumblr.zendesk.com/hc/en-us/articles/231256548-Custom-domains")
			a.session.Out.Warn("%s: vulnerable to takeover on Tumblr\n", p.URL)
		}
		return true
	}
	return false
}

func (a *URLTakeoverDetector) detectUserVoice(p *core.Page, addrs []string, cname string, body string) bool {
	if !strings.HasSuffix(cname, ".uservoice.com.") {
		return false
	}
	p.AddTag("UserVoice", "info", "https://www.uservoice.com/")
	if strings.Contains(body, "This UserVoice subdomain is currently available!") {
		p.AddTag("Domain Takeover", "danger", "https://developer.uservoice.com/docs/site/domain-aliasing/")
		a.session.Out.Warn("%s: vulnerable to takeover on UserVoice\n", p.URL)
	}
	return true
}

func (a *URLTakeoverDetector) detectWordpress(p *core.Page, addrs []string, cname string, body string) bool {
	if !strings.HasSuffix(cname, ".wordpress.com.") {
		return false
	}
	if strings.Contains(body, "Do you want to register") {
		p.AddTag("Domain Takeover", "danger", "https://en.support.wordpress.com/domains/map-subdomain/")
		a.session.Out.Warn("%s: vulnerable to takeover on Wordpress\n", p.URL)
	}
	return true
}

func (a *URLTakeoverDetector) detectSmugMug(p *core.Page, addrs []string, cname string, body string) bool {
	if cname != "domains.smugmug.com." {
		return false
	}
	p.AddTag("SmugMug", "info", "https://www.smugmug.com/")
	if body == "" {
		p.AddTag("Domain Takeover", "danger", "https://help.smugmug.com/use-a-custom-domain-BymMexwJVHG")
		a.session.Out.Warn("%s: vulnerable to takeover on SmugMug\n", p.URL)
	}
	return true
}

func (a *URLTakeoverDetector) detectStrikingly(p *core.Page, addrs []string, cname string, body string) bool {
	detected := false
	for _, addr := range addrs {
		if addr == "54.183.102.22" {
			detected = true
			break
		}
	}
	if strings.HasSuffix(cname, ".s.strikinglydns.com.") {
		detected = true
	}
	if detected {
		p.AddTag("Strikingly", "info", "https://www.strikingly.com/")
		if strings.Contains(body, "But if you're looking to build your own website,") {
			p.AddTag("Domain Takeover", "danger", "https://support.strikingly.com/hc/en-us/articles/215046947-Connect-Custom-Domain")
			a.session.Out.Warn("%s: vulnerable to takeover on Strikingly\n", p.URL)
		}
		return true
	}
	return false
}

func (a *URLTakeoverDetector) detectUptimeRobot(p *core.Page, addrs []string, cname string, body string) bool {
	if cname != "stats.uptimerobot.com." {
		return false
	}
	p.AddTag("UptimeRobot", "info", "https://uptimerobot.com/")
	if strings.Contains(body, "This public status page <b>does not seem to exist</b>.") {
		p.AddTag("Domain Takeover", "danger", "https://blog.uptimerobot.com/introducing-public-status-pages-yay/")
		a.session.Out.Warn("%s: vulnerable to takeover on UptimeRobot\n", p.URL)
	}
	return true
}

func (a *URLTakeoverDetector) detectPantheon(p *core.Page, addrs []string, cname string, body string) bool {
	if !strings.HasSuffix(cname, ".pantheonsite.io.") {
		return false
	}
	p.AddTag("Pantheon", "info", "https://pantheon.io/")
	if strings.Contains(body, "The gods are wise") {
		p.AddTag("Domain Takeover", "danger", "https://pantheon.io/docs/domains/")
		a.session.Out.Warn("%s: vulnerable to takeover on Pantheon\n", p.URL)
	}
	return true
}
