package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/michenriksen/aquatone/agents"
	"github.com/michenriksen/aquatone/core"
	"github.com/michenriksen/aquatone/parsers"
)

var (
	sess *core.Session
	err  error
)

func isURL(s string) bool {
	u, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}
	if u.Scheme == "" {
		return false
	}
	return true
}

func hasSupportedScheme(s string) bool {
	u, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}
	if u.Scheme == "http" || u.Scheme == "https" {
		return true
	}
	return false
}

func main() {
	if sess, err = core.NewSession(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fi, err := os.Stat(*sess.Options.OutDir)

	if os.IsNotExist(err) {
		sess.Out.Fatal("Output destination %s does not exist\n", *sess.Options.OutDir)
		os.Exit(1)
	}

	if !fi.IsDir() {
		sess.Out.Fatal("Output destination must be a directory\n")
		os.Exit(1)
	}

	sess.Out.Important("%s v%s started at %s\n\n", core.Name, core.Version, sess.Stats.StartedAt.Format(time.RFC3339))

	agents.NewTCPPortScanner().Register(sess)
	agents.NewURLPublisher().Register(sess)
	agents.NewURLRequester().Register(sess)
	agents.NewURLLogger().Register(sess)
	agents.NewURLScreenshotter().Register(sess)
	agents.NewURLTechnologyFingerprinter().Register(sess)
	agents.NewURLTakeoverDetector().Register(sess)

	reader := bufio.NewReader(os.Stdin)
	var targets []string

	if *sess.Options.Nmap {
		parser := parsers.NewNmapParser()
		targets, err = parser.Parse(reader)
		if err != nil {
			sess.Out.Fatal("Unable to parse input as Nmap/Masscan XML: %s\n", err)
			os.Exit(1)
		}
	} else {
		parser := parsers.NewRegexParser()
		targets, err = parser.Parse(reader)
		if err != nil {
			sess.Out.Fatal("Unable to parse input.\n")
			os.Exit(1)
		}
	}

	if len(targets) == 0 {
		sess.Out.Fatal("No targets found in input.\n")
		os.Exit(1)
	}

	sess.Out.Important("Targets    : %d\n", len(targets))
	sess.Out.Important("Threads    : %d\n", *sess.Options.Threads)
	sess.Out.Important("Ports      : %s\n", strings.Trim(strings.Replace(fmt.Sprint(sess.Ports), " ", ", ", -1), "[]"))
	sess.Out.Important("Output dir : %s\n\n", *sess.Options.OutDir)

	for _, target := range targets {
		if isURL(target) {
			if hasSupportedScheme(target) {
				sess.EventBus.Publish(core.URL, target)
			}
		} else {
			sess.EventBus.Publish(core.Host, target)
		}
	}

	time.Sleep(1 * time.Second)
	sess.EventBus.WaitAsync()
	sess.WaitGroup.Wait()

	sess.Out.Important("\nClustering similar sites...")
	pageStructures := make(map[string][]string)
	var pageClusters [][]*core.ResponsiveURL

	f, _ := os.OpenFile(sess.GetFilePath("aquatone_urls.txt"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	for _, responsiveURL := range sess.ResponsiveURLs {
		filename := sess.GetFilePath(fmt.Sprintf("html/%s.html", agents.BaseFilenameFromURL(responsiveURL.URL)))
		body, err := os.Open(filename)
		if err != nil {
			continue
		}
		structure, _ := core.GetPageStructure(body)
		pageStructures[responsiveURL.URL] = structure
		f.WriteString(responsiveURL.URL + "\n")
	}
	f.Close()

	// Loop over URL and page structure pairs
	for url, structure := range pageStructures {
		foundCluster := false
		// Loop over existing page clusters
		for i, cluster := range pageClusters {
			addToCluster := true
			// Loop over pages in cluster and check if similarity for all are 0.80 or above
			for _, url2 := range cluster {
				if core.GetSimilarity(structure, pageStructures[url2.URL]) < 0.80 {
					addToCluster = false
				}
			}
			// Add to cluster if similarity between all pages are 0.80 or above
			if addToCluster {
				foundCluster = true
				pageClusters[i] = append(pageClusters[i], sess.ResponsiveURLs[url])
				break
			}
		}
		// If a cluster was not found for the page, create a new cluster for the page
		if !foundCluster {
			pageClusters = append(pageClusters, []*core.ResponsiveURL{sess.ResponsiveURLs[url]})
		}
	}

	sess.Out.Important(" done\n")
	sess.Out.Important("Generating HTML report...")

	reportData := core.ReportData{
		Session: sess,
	}

	for _, urls := range pageClusters {
		cluster, err := core.NewCluster(urls, sess)
		if err != nil {
			sess.Out.Fatal("Error during report generation: %s\n", err)
			os.Exit(1)
		}
		reportData.Clusters = append(reportData.Clusters, cluster)
	}

	report := core.NewReport(reportData)
	f, err = os.OpenFile(sess.GetFilePath("aquatone_report.html"), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		sess.Out.Fatal("Error during report generation: %s\n", err)
		os.Exit(1)
	}
	err = report.Render(f)
	if err != nil {
		sess.Out.Fatal("Error during report generation: %s\n", err)
		os.Exit(1)
	}

	sess.Out.Important(" done\n\n")

	sess.End()

	sess.Out.Important("Time:\n")
	sess.Out.Info(" - Started at  : %v\n", sess.Stats.StartedAt.Format(time.RFC3339))
	sess.Out.Info(" - Finished at : %v\n", sess.Stats.FinishedAt.Format(time.RFC3339))
	sess.Out.Info(" - Duration    : %v\n\n", sess.Stats.Duration().Round(time.Second))

	sess.Out.Important("Requests:\n")
	sess.Out.Info(" - Successful : %v\n", sess.Stats.RequestSuccessful)
	sess.Out.Info(" - Failed     : %v\n\n", sess.Stats.RequestFailed)

	sess.Out.Info(" - 2xx : %v\n", sess.Stats.ResponseCode2xx)
	sess.Out.Info(" - 3xx : %v\n", sess.Stats.ResponseCode3xx)
	sess.Out.Info(" - 4xx : %v\n", sess.Stats.ResponseCode4xx)
	sess.Out.Info(" - 5xx : %v\n\n", sess.Stats.ResponseCode5xx)

	sess.Out.Important("Screenshots:\n")
	sess.Out.Info(" - Successful : %v\n", sess.Stats.ScreenshotSuccessful)
	sess.Out.Info(" - Failed     : %v\n\n", sess.Stats.ScreenshotFailed)

	sess.Out.Important("Wrote HTML report to: %s\n\n", sess.GetFilePath("aquatone_report.html"))
}
