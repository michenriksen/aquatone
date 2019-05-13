package core

import (
	"crypto/sha1"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type Header struct {
	Name              string `json:"name"`
	Value             string `json:"value"`
	DecreasesSecurity bool   `json:"decreasesSecurity"`
	IncreasesSecurity bool   `json:"increasesSecurity"`
}

func (h *Header) SetSecurityFlags() {
	if h.decreasesSecurity() {
		h.DecreasesSecurity = true
		h.IncreasesSecurity = false
	} else if h.increasesSecurity() {
		h.DecreasesSecurity = false
		h.IncreasesSecurity = true
	} else {
		h.DecreasesSecurity = false
		h.IncreasesSecurity = false
	}
}

func (h Header) decreasesSecurity() bool {
	switch strings.ToLower(h.Name) {
	case "server", "wpe-Backend", "x-powered-by", "x-cf-powered-by", "x-pingback":
		return true
	case "access-control-allow-origin":
		if h.Value == "*" {
			return true
		}
	case "x-xss-protection":
		if !strings.HasPrefix(h.Value, "1") {
			return true
		}
	}
	return false
}

func (h Header) increasesSecurity() bool {
	switch strings.ToLower(h.Name) {
	case "content-security-policy", "content-security-policy-report-only":
		return true
	case "strict-transport-security":
		return true
	case "x-frame-options":
		return true
	case "referrer-policy":
		return true
	case "public-key-pins":
		return true
	case "x-permitted-cross-domain-policies":
		if strings.ToLower(h.Value) == "master-only" {
			return true
		}
	case "x-content-type-options":
		if strings.ToLower(h.Value) == "nosniff" {
			return true
		}
	case "x-xss-protection":
		if strings.HasPrefix(h.Value, "1") {
			return true
		}
	}
	return false
}

type Tag struct {
	Text string `json:"text"`
	Type string `json:"type"`
	Link string `json:"link"`
	Hash string `json:"hash"`
}

func (t Tag) HasLink() bool {
	if t.Link != "" {
		return true
	}
	return false
}

type Note struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type Page struct {
	sync.Mutex
	UUID           string   `json:"uuid"`
	URL            string   `json:"url"`
	Hostname       string   `json:"hostname"`
	Addrs          []string `json:"addrs"`
	Status         string   `json:"status"`
	PageTitle      string   `json:"pageTitle"`
	PageStructure  []string `json:"-"`
	HeadersPath    string   `json:"headersPath"`
	BodyPath       string   `json:"bodyPath"`
	ScreenshotPath string   `json:"screenshotPath"`
	HasScreenshot  bool     `json:"hasScreenshot"`
	Headers        []Header `json:"headers"`
	Tags           []Tag    `json:"tags"`
	Notes          []Note   `json:"notes"`
}

func (p *Page) AddHeader(name string, value string) {
	p.Lock()
	defer p.Unlock()
	header := Header{
		Name:  name,
		Value: value,
	}
	header.SetSecurityFlags()
	p.Headers = append(p.Headers, header)
}

func (p *Page) AddTag(text string, tagType string, link string) {
	p.Lock()
	defer p.Unlock()

	h := sha1.New()
	io.WriteString(h, text)
	io.WriteString(h, tagType)
	io.WriteString(h, link)

	p.Tags = append(p.Tags, Tag{
		Text: text,
		Type: tagType,
		Link: link,
		Hash: fmt.Sprintf("%x", h.Sum(nil)),
	})
}

func (p *Page) AddNote(text string, noteType string) {
	p.Lock()
	defer p.Unlock()
	p.Notes = append(p.Notes, Note{
		Text: text,
		Type: noteType,
	})
}

func (p *Page) BaseFilename() string {
	u := p.ParsedURL()
	h := sha1.New()
	io.WriteString(h, u.Path)
	io.WriteString(h, u.Fragment)

	pathHash := fmt.Sprintf("%x", h.Sum(nil))[0:16]
	host := strings.Replace(u.Host, ":", "__", 1)
	filename := fmt.Sprintf("%s__%s__%s", u.Scheme, strings.Replace(host, ".", "_", -1), pathHash)
	return strings.ToLower(filename)
}

func (p *Page) ParsedURL() *url.URL {
	parsedURL, _ := url.Parse(p.URL)
	return parsedURL
}

func (p *Page) IsIPHost() bool {
	return net.ParseIP(p.ParsedURL().Hostname()) != nil
}

func NewPage(pageURL string) (*Page, error) {
	u, err := url.Parse(pageURL)
	if err != nil {
		return nil, err
	}

	return &Page{
		UUID:     uuid.New().String(),
		URL:      pageURL,
		Hostname: u.Hostname(),
	}, nil
}
