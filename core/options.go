package core

import (
	"flag"
	"fmt"
	"strings"
)

type Options struct {
	Threads           *int
	OutDir            *string
	SessionPath       *string
	TemplatePath      *string
	Proxy             *string
	ChromePath        *string
	Resolution        *string
	Ports             *string
	ScanTimeout       *int
	HTTPTimeout       *int
	ScreenshotTimeout *int
	ScreenshotDelay   *int
	Nmap              *bool
	SaveBody          *bool
	Silent            *bool
	Debug             *bool
	Version           *bool
}

func ParseOptions() (Options, error) {
	options := Options{
		Threads:           flag.Int("threads", 0, "Number of concurrent threads (default number of logical CPUs)"),
		OutDir:            flag.String("out", ".", "Directory to write files to"),
		SessionPath:       flag.String("session", "", "Load Aquatone session file and generate HTML report"),
		TemplatePath:      flag.String("template-path", "", "Path to HTML template to use for report"),
		Proxy:             flag.String("proxy", "", "Proxy to use for HTTP requests"),
		ChromePath:        flag.String("chrome-path", "", "Full path to the Chrome/Chromium executable to use. By default, aquatone will search for Chrome or Chromium"),
		Resolution:        flag.String("resolution", "1440,900", "screenshot resolution"),
		Ports:             flag.String("ports", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(MediumPortList)), ","), "[]"), "Ports to scan on hosts. Supported list aliases: small, medium, large, xlarge"),
		ScanTimeout:       flag.Int("scan-timeout", 100, "Timeout in miliseconds for port scans"),
		HTTPTimeout:       flag.Int("http-timeout", 3*1000, "Timeout in miliseconds for HTTP requests"),
		ScreenshotTimeout: flag.Int("screenshot-timeout", 30*1000, "Timeout in miliseconds for screenshots"),
		ScreenshotDelay:   flag.Int("screenshot-delay", 0, "The delay before taking screenshots"),
		Nmap:              flag.Bool("nmap", false, "Parse input as Nmap/Masscan XML"),
		SaveBody:          flag.Bool("save-body", true, "Save response bodies to files"),
		Silent:            flag.Bool("silent", false, "Suppress all output except for errors"),
		Debug:             flag.Bool("debug", false, "Print debugging information"),
		Version:           flag.Bool("version", false, "Print current Aquatone version"),
	}

	flag.Parse()

	return options, nil
}
