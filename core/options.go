package core

import (
	"flag"
	"fmt"
	"strings"
)

type Options struct {
	Threads           *int
	OutDir            *string
	Proxy             *string
	ChromePath        *string
	Resolution        *string
	Ports             *string
	ScanTimeout       *int
	HTTPTimeout       *int
	ScreenshotTimeout *int
	Nmap              *bool
	SaveBody          *bool
	Waf               *bool
	Silent            *bool
	Debug             *bool
}

func ParseOptions() (Options, error) {
	options := Options{
		Threads:           flag.Int("threads", 0, "Number of concurrent threads (default number of logical CPUs)"),
		OutDir:            flag.String("out", ".", "Directory to write files to"),
		Proxy:             flag.String("proxy", "", "Proxy to use for HTTP requests"),
		ChromePath:        flag.String("chrome-path", "", "Full path to the Chrome/Chromium executable to use. By default, aquatone will search for Chrome or Chromium"),
		Resolution:        flag.String("resolution", "1440,900", "screenshot resolution"),
		Ports:             flag.String("ports", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(MediumPortList)), ","), "[]"), "Ports to scan on hosts. Supported list aliases: small, medium, large, xlarge"),
		ScanTimeout:       flag.Int("scan-timeout", 100, "Timeout in miliseconds for port scans"),
		HTTPTimeout:       flag.Int("http-timeout", 3*1000, "Timeout in miliseconds for HTTP requests"),
		ScreenshotTimeout: flag.Int("screenshot-timeout", 30*1000, "Timeout in miliseconds for screenshots"),
		Nmap:              flag.Bool("nmap", false, "Parse input as Nmap/Masscan XML"),
		SaveBody:          flag.Bool("save-body", true, "Save response bodies to files"),
		Waf:               flag.Bool("waf", false, "Attempt to bypass weak ACLs"),
		Silent:            flag.Bool("silent", false, "Suppress all output except for errors"),
		Debug:             flag.Bool("debug", false, "Print debugging information"),
	}

	flag.Parse()

	return options, nil
}
