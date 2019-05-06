package parsers

import (
	"io"
	"io/ioutil"

	"github.com/michenriksen/aquatone/core"

	"github.com/lair-framework/go-nmap"
)

type NmapParser struct{}

func NewNmapParser() *NmapParser {
	return &NmapParser{}
}

func (p *NmapParser) Parse(r io.Reader) ([]string, error) {
	var targets []string
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return targets, err
	}
	scan, err := nmap.Parse(bytes)
	if err != nil {
		return targets, err
	}

	for _, host := range scan.Hosts {
		urls := p.hostToURLs(host)
		for _, url := range urls {
			targets = append(targets, url)
		}
	}

	return targets, nil
}

func (p *NmapParser) isHTTPPort(port int) bool {
	for _, p := range core.XLargePortList {
		if p == port {
			return true
		}
	}
	return false
}

func (p *NmapParser) hostToURLs(host nmap.Host) []string {
	var urls []string
	for _, port := range host.Ports {
		if port.State.State != "open" {
			continue
		}

		var protocol string
		if port.Service.Name == "ssl" {
			protocol = "https"
		} else if port.Service.Tunnel == "ssl" && (port.Service.Name != "smtp" && port.Service.Name != "imap" && port.Service.Name != "pop3") {
			protocol = "https"
		} else if port.Service.Name == "http" || port.Service.Name == "http-alt" {
			protocol = "http"
		} else {
			if !p.isHTTPPort(port.PortId) {
				continue
			}
		}

		if len(host.Hostnames) > 0 {
			for _, hostname := range host.Hostnames {
				urls = append(urls, core.HostAndPortToURL(hostname.Name, port.PortId, protocol))
			}
		} else {
			for _, address := range host.Addresses {
				if address.AddrType == "mac" {
					continue
				}
				urls = append(urls, core.HostAndPortToURL(address.Addr, port.PortId, protocol))
			}
		}
	}

	return urls
}
