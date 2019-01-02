package agents

import (
	"fmt"
	"net"
	"net/url"
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
	parsedURL, err := url.Parse(u)
	if err != nil {
		a.session.Out.Debug("[%s] Unable to parse %s as an URL\n", a.ID(), u)
		return
	}
	if !a.isDomainURL(parsedURL) {
		a.session.Out.Debug("[%s] Skipping takeover detection on IP URL %s\n", a.ID(), u)
		return
	}
	a.session.WaitGroup.Add()
	go func(u *url.URL) {
		defer a.session.WaitGroup.Done()
		a.runDetectorFunctions(u)
	}(parsedURL)
}

func (a *URLTakeoverDetector) isDomainURL(u *url.URL) bool {
	return net.ParseIP(u.Hostname()) == nil
}

func (a *URLTakeoverDetector) runDetectorFunctions(u *url.URL) {
	hostname := u.Hostname()
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

	baseFileName := BaseFilenameFromURL(u.String())
	body, err := a.session.ReadFile(fmt.Sprintf("html/%s.html", baseFileName))
	if err != nil {
		a.session.Out.Debug("[%s] Error reading HTML body file for %s: %s\n", a.ID(), u.String(), err)
		return
	}

	if a.detectGithubPages(u.String(), addrs, cname, string(body)) {
		return
	}

	if a.detectAmazonS3(u.String(), addrs, cname, string(body)) {
		return
	}

	if a.detectCampaignMonitor(u.String(), addrs, cname, string(body)) {
		return
	}

	if a.detectCargoCollective(u.String(), addrs, cname, string(body)) {
		return
	}

	if a.detectFeedPress(u.String(), addrs, cname, string(body)) {
		return
	}

	if a.detectGhost(u.String(), addrs, cname, string(body)) {
		return
	}

	if a.detectHelpjuice(u.String(), addrs, cname, string(body)) {
		return
	}

	if a.detectHelpScout(u.String(), addrs, cname, string(body)) {
		return
	}

	if a.detectHeroku(u.String(), addrs, cname, string(body)) {
		return
	}

	if a.detectJetBrains(u.String(), addrs, cname, string(body)) {
		return
	}

	if a.detectMicrosoftAzure(u.String(), addrs, cname, string(body)) {
		return
	}

	if a.detectReadme(u.String(), addrs, cname, string(body)) {
		return
	}

	if a.detectSurge(u.String(), addrs, cname, string(body)) {
		return
	}

	if a.detectTumblr(u.String(), addrs, cname, string(body)) {
		return
	}

	if a.detectUserVoice(u.String(), addrs, cname, string(body)) {
		return
	}

	if a.detectWordpress(u.String(), addrs, cname, string(body)) {
		return
	}

	if a.detectSmugMug(u.String(), addrs, cname, string(body)) {
		return
	}

	if a.detectStrikingly(u.String(), addrs, cname, string(body)) {
		return
	}

	if a.detectUptimeRobot(u.String(), addrs, cname, string(body)) {
		return
	}

	if a.detectPantheon(u.String(), addrs, cname, string(body)) {
		return
	}

	if a.detectEdgecastCDN(u.String(), addrs, cname, string(body)) {
		return
	}

	if a.detectSquarespace(u.String(), addrs, cname, string(body)) {
		return
	}
}

func (a *URLTakeoverDetector) detectGithubPages(u string, addrs []string, cname string, body string) bool {
	githubAddrs := [...]string{"185.199.108.153", "185.199.109.153", "185.199.110.153", "185.199.111.153"}
	fingerprints := [...]string{"There isn't a GitHub Pages site here.", "For root URLs (like http://example.com/) you must provide an index.html file"}
	for _, githubAddr := range githubAddrs {
		for _, addr := range addrs {
			if addr == githubAddr {
				for _, fingerprint := range fingerprints {
					if strings.Contains(body, fingerprint) {
						a.session.AddTagToResponsiveURL(u, "Domain Takeover", "danger", "https://help.github.com/articles/using-a-custom-domain-with-github-pages/")
						return true
					}
				}
				return true
			}
		}
	}
	return false
}

func (a *URLTakeoverDetector) detectAmazonS3(u string, addrs []string, cname string, body string) bool {
	fingerprints := [...]string{"NoSuchBucket", "The specified bucket does not exist"}
	if !strings.HasSuffix(cname, ".amazonaws.com.") {
		return false
	}
	for _, fingerprint := range fingerprints {
		if strings.Contains(body, fingerprint) {
			a.session.AddTagToResponsiveURL(u, "Domain Takeover", "danger", "https://docs.aws.amazon.com/AmazonS3/latest/dev/website-hosting-custom-domain-walkthrough.html")
			return true
		}
	}
	return true
}

func (a *URLTakeoverDetector) detectCampaignMonitor(u string, addrs []string, cname string, body string) bool {
	if cname != "cname.createsend.com." {
		return false
	}
	a.session.AddTagToResponsiveURL(u, "Campaign Monitor", "info", "https://www.campaignmonitor.com/")
	if strings.Contains(body, "Double check the URL or ") {
		a.session.AddTagToResponsiveURL(u, "Domain Takeover", "danger", "https://help.campaignmonitor.com/custom-domain-names")
		return true
	}
	return true
}

func (a *URLTakeoverDetector) detectCargoCollective(u string, addrs []string, cname string, body string) bool {
	if cname != "subdomain.cargocollective.com." {
		return false
	}
	a.session.AddTagToResponsiveURL(u, "Cargo Collective", "info", "https://cargocollective.com/")
	if strings.Contains(body, "404 Not Found") {
		a.session.AddTagToResponsiveURL(u, "Domain Takeover", "danger", "https://support.2.cargocollective.com/Using-a-Third-Party-Domain")
		return true
	}
	return true
}

func (a *URLTakeoverDetector) detectFeedPress(u string, addrs []string, cname string, body string) bool {
	if cname != "redirect.feedpress.me." {
		return false
	}
	a.session.AddTagToResponsiveURL(u, "FeedPress", "info", "https://feed.press/")
	if strings.Contains(body, "The feed has not been found.") {
		a.session.AddTagToResponsiveURL(u, "Domain Takeover", "danger", "https://support.feed.press/article/61-how-to-create-a-custom-hostname")
		return true
	}
	return true
}

func (a *URLTakeoverDetector) detectGhost(u string, addrs []string, cname string, body string) bool {
	if !strings.HasSuffix(cname, ".ghost.io.") {
		return false
	}
	if strings.Contains(body, "The thing you were looking for is no longer here, or never was") {
		a.session.AddTagToResponsiveURL(u, "Domain Takeover", "danger", "https://docs.ghost.org/faq/using-custom-domains/")
		return true
	}
	return true
}

func (a *URLTakeoverDetector) detectHelpjuice(u string, addrs []string, cname string, body string) bool {
	if !strings.HasSuffix(cname, ".helpjuice.com.") {
		return false
	}
	a.session.AddTagToResponsiveURL(u, "Helpjuice", "info", "https://helpjuice.com/")
	if strings.Contains(body, "We could not find what you're looking for.") {
		a.session.AddTagToResponsiveURL(u, "Domain Takeover", "danger", "https://help.helpjuice.com/34339-getting-started/custom-domain")
		return true
	}
	return false
}

func (a *URLTakeoverDetector) detectHelpScout(u string, addrs []string, cname string, body string) bool {
	if !strings.HasSuffix(cname, ".helpscoutdocs.com.") {
		return false
	}
	a.session.AddTagToResponsiveURL(u, "HelpScout", "info", "https://www.helpscout.net/")
	if strings.Contains(body, "No settings were found for this company:") {
		a.session.AddTagToResponsiveURL(u, "Domain Takeover", "danger", "https://docs.helpscout.net/article/42-setup-custom-domain")
		return true
	}
	return true
}

func (a *URLTakeoverDetector) detectHeroku(u string, addrs []string, cname string, body string) bool {
	herokuCnames := [...]string{".herokudns.com.", ".herokuapp.com.", ".herokussl.com."}
	for _, herokuCname := range herokuCnames {
		if strings.HasSuffix(cname, herokuCname) {
			a.session.AddTagToResponsiveURL(u, "Heroku", "info", "https://www.heroku.com/")
			if strings.Contains(body, "No such app") {
				a.session.AddTagToResponsiveURL(u, "Domain Takeover", "danger", "https://devcenter.heroku.com/articles/custom-domains")
				return true
			}
			return true
		}
	}
	return false
}

func (a *URLTakeoverDetector) detectJetBrains(u string, addrs []string, cname string, body string) bool {
	if !strings.HasSuffix(cname, ".myjetbrains.com.") {
		return false
	}
	a.session.AddTagToResponsiveURL(u, "JetBrains", "info", "https://www.jetbrains.com/")
	if strings.Contains(body, "is not a registered InCloud YouTrack") {
		a.session.AddTagToResponsiveURL(u, "Domain Takeover", "danger", "https://www.jetbrains.com/help/youtrack/incloud/Domain-Settings.html#use-custom-domain-name")
		return true
	}
	return true
}

func (a *URLTakeoverDetector) detectMicrosoftAzure(u string, addrs []string, cname string, body string) bool {
	if !strings.HasSuffix(cname, ".azurewebsites.net.") {
		return false
	}
	a.session.AddTagToResponsiveURL(u, "Microsoft Azure", "info", "https://azure.microsoft.com/")
	if strings.Contains(body, "404 Web Site not found") {
		a.session.AddTagToResponsiveURL(u, "Domain Takeover", "danger", "https://docs.microsoft.com/en-us/azure/app-service/app-service-web-tutorial-custom-domain")
		return true
	}
	return true
}

func (a *URLTakeoverDetector) detectReadme(u string, addrs []string, cname string, body string) bool {
	readmeCnames := [...]string{".readme.io.", ".readmessl.com."}
	for _, readmeCname := range readmeCnames {
		if strings.HasSuffix(cname, readmeCname) {
			a.session.AddTagToResponsiveURL(u, "Readme", "info", "https://readme.io/")
			if strings.Contains(body, "Project doesnt exist... yet!") {
				a.session.AddTagToResponsiveURL(u, "Domain Takeover", "danger", "https://readme.readme.io/docs/setting-up-custom-domain")
				return true
			}
			return true
		}
	}
	return false
}

func (a *URLTakeoverDetector) detectSurge(u string, addrs []string, cname string, body string) bool {
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
		a.session.AddTagToResponsiveURL(u, "Surge", "info", "https://surge.sh/")
		if strings.Contains(body, "project not found") {
			a.session.AddTagToResponsiveURL(u, "Domain Takeover", "danger", "https://surge.sh/help/adding-a-custom-domain")
		}
		return true
	}
	return false
}

func (a *URLTakeoverDetector) detectTumblr(u string, addrs []string, cname string, body string) bool {
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
			a.session.AddTagToResponsiveURL(u, "Domain Takeover", "danger", "https://tumblr.zendesk.com/hc/en-us/articles/231256548-Custom-domains")
		}
		return true
	}
	return false
}

func (a *URLTakeoverDetector) detectUserVoice(u string, addrs []string, cname string, body string) bool {
	if !strings.HasSuffix(cname, ".uservoice.com.") {
		return false
	}
	a.session.AddTagToResponsiveURL(u, "UserVoice", "info", "https://www.uservoice.com/")
	if strings.Contains(body, "This UserVoice subdomain is currently available!") {
		a.session.AddTagToResponsiveURL(u, "Domain Takeover", "danger", "https://developer.uservoice.com/docs/site/domain-aliasing/")
	}
	return true
}

func (a *URLTakeoverDetector) detectWordpress(u string, addrs []string, cname string, body string) bool {
	if !strings.HasSuffix(cname, ".wordpress.com.") {
		return false
	}
	if strings.Contains(body, "Do you want to register") {
		a.session.AddTagToResponsiveURL(u, "Domain Takeover", "danger", "https://en.support.wordpress.com/domains/map-subdomain/")
	}
	return true
}

func (a *URLTakeoverDetector) detectSmugMug(u string, addrs []string, cname string, body string) bool {
	if cname != "domains.smugmug.com." {
		return false
	}
	a.session.AddTagToResponsiveURL(u, "SmugMug", "info", "https://www.smugmug.com/")
	if body == "" {
		a.session.AddTagToResponsiveURL(u, "Domain Takeover", "danger", "https://help.smugmug.com/use-a-custom-domain-BymMexwJVHG")
	}
	return true
}

func (a *URLTakeoverDetector) detectStrikingly(u string, addrs []string, cname string, body string) bool {
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
		a.session.AddTagToResponsiveURL(u, "Strikingly", "info", "https://www.strikingly.com/")
		if strings.Contains(body, "But if you're looking to build your own website,") {
			a.session.AddTagToResponsiveURL(u, "Domain Takeover", "danger", "https://support.strikingly.com/hc/en-us/articles/215046947-Connect-Custom-Domain")
		}
		return true
	}
	return false
}

func (a *URLTakeoverDetector) detectUptimeRobot(u string, addrs []string, cname string, body string) bool {
	if cname != "stats.uptimerobot.com." {
		return false
	}
	a.session.AddTagToResponsiveURL(u, "UptimeRobot", "info", "https://uptimerobot.com/")
	if strings.Contains(body, "This public status page <b>does not seem to exist</b>.") {
		a.session.AddTagToResponsiveURL(u, "Domain Takeover", "danger", "https://blog.uptimerobot.com/introducing-public-status-pages-yay/")
	}
	return true
}

func (a *URLTakeoverDetector) detectPantheon(u string, addrs []string, cname string, body string) bool {
	if !strings.HasSuffix(cname, ".pantheonsite.io.") {
		return false
	}
	a.session.AddTagToResponsiveURL(u, "Pantheon", "info", "https://pantheon.io/")
	if strings.Contains(body, "The gods are wise") {
		a.session.AddTagToResponsiveURL(u, "Domain Takeover", "danger", "https://pantheon.io/docs/domains/")
	}
	return true
}

func (a *URLTakeoverDetector) detectEdgecastCDN(u string, addrs []string, cname string, body string) bool {
	if !strings.HasSuffix(cname, ".edgecastcdn.net") {
		return false
	}
	a.session.AddTagToResponsiveURL(u, "EdgecastCDN", "info", "https://www.verizondigitalmedia.com")
	if strings.Contains(body, "404 - Not Found") {
		a.session.AddTagToResponsiveURL(u, "Domain Takeover", "danger", "https://www.verizondigitalmedia.com/platform/edgecast-cdn/")
	}
	return true
}

func (a *URLTakeoverDetector) detectSquarespace(u string, addrs []string, cname string, body string) bool {
	detected := false
	_, subnet1, _ := net.ParseCIDR("198.49.23.0/24")
	_, subnet2, _ := net.ParseCIDR("198.185.159.0/24")

	for _, addr := range addrs {
		ip := net.ParseIP(addr)
		if ip != nil && (subnet1.Contains(ip) || subnet2.Contains(ip)) {
			detected = true
			break
		}
	}
	if strings.HasSuffix(cname, ".squarespace.com.") {
		detected = true
	}
	if detected {
		a.session.AddTagToResponsiveURL(u, "Squarespace", "info", "https://www.squarespace.com/")
		if strings.Contains(body, "Squarespace - Claim This Domain") {
			a.session.AddTagToResponsiveURL(u, "Domain Takeover", "danger", "https://support.squarespace.com/hc/en-us/articles/115002755267")
		}
		return true
	}
	return false
}
