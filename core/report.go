package core

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"
)

type Header struct {
	Name  string
	Value string
}

func (h Header) IsInsecure() bool {
	switch strings.ToLower(h.Name) {
	case "server", "x-powered-by":
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

func (h Header) IsSecure() bool {
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

type Page struct {
	URL            string
	Status         string
	Headers        []Header
	HeadersPath    string
	BodyPath       string
	ScreenshotPath string
	HasScreenshot  bool
	Tags           []Tag
	Notes          []Note
}

const (
	Template = `
<!doctype html>
<html lang="en">
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
		<link rel="stylesheet" href="https://bootswatch.com/4/darkly/bootstrap.min.css" crossorigin="anonymous">
		<title>AQUATONE REPORT</title>
		<style type="text/css">
			a {
				color: rgb(52, 152, 219) !important;
			}

			a.badge {
				color: inherit !important;
			}

			footer {
				margin-top: 20px;
				padding: 20px;
				text-align: center;
				font-size: 12px;
				color: rgb(68, 68, 68);
			}

			footer a {
				color: inherit !important;
				text-decoration: underline;
			}

			.logo {
				padding-left: 30px;
				padding-top: 10px;
				font-weight: bold;
				font-size: 16px;
			}

			.cluster {
				border-bottom: 1px solid rgb(68, 68, 68);
				padding: 30px 20px 20px 20px;
				overflow-x: auto;
				white-space: nowrap;
			}

			.cluster:nth-child(even) {
				background-color: rgba(0, 0, 0, 0.075);
				box-shadow: inset 0px 6px 8px rgb(24, 24, 24);
			}

			.page {
				display: inline-block;
				margin: 10px;
				width: 470px;
				overflow: hidden;
				box-shadow: 10px 10px 8px rgb(24, 24, 24);
			}

			.page .card-text {
        white-space: normal;
      }

			.page .screenshot-container {
			  position: relative;
			  width: 470px;
			  height: 293px;
			  overflow: hidden;
			}

			.page .screenshot-container img.screenshot {
			  position: absolute;
			  top: 0;
			  left: 0;
			  width: 100%;
			  height: 100%;
			  background-repeat: no-repeat;
			  background-position: center;
			  background-size: cover;
			  transition: transform .5s ease-out;
			}

			.page .response-headers-container {
				display: none;
			}

			table.response-headers td {
				font-family: Anonymous Pro,Consolas,Menlo,Monaco,Lucida Console,Liberation Mono,DejaVu Sans Mono,Bitstream Vera Sans Mono,Courier New,monospace,serif;
			}

			table.response-headers tr.insecure td {
				color: #E74C3C;
				font-weight: bold;
			}

			table.response-headers tr.secure td {
				color: rgb(0, 188, 140);
				font-weight: bold;
			}

			img.screenshot {
				transition: transform .2s ease-out;
			}

			img.no-screenshot {
				cursor: not-allowed;
			}
		</style>
	</head>
	<body>
		<pre class="logo">
                          __
.---.-.-----.--.--.---.-.|  |_.-----.-----.-----.
|  _  |  _  |  |  |  _  ||   _|  _  |     |  -__|
|___._|__   |_____|___._||____|_____|__|__|_____|
         |__| report from {{.Session.Stats.StartedAt.Format "2006-01-02T15:04:05Z07:00"}}
		</pre>
		{{range .Clusters}}
			<div class="cluster">
				{{range .}}
					<div class="page card mb-3">
						<div class="card-body">
							<h5 class="card-title">{{.URL}}</h5>
							<h6 class="card-subtitle text-muted">{{.Status}}</h6>
							<p class="card-text">
								{{range .Tags}}
									{{if .HasLink}}
										<a href="{{.Link}}" target="_blank" class="badge badge-pill badge-{{.Type}}">{{.Text}}</a>
									{{else}}
										<span class="badge badge-pill badge-{{.Type}}">{{.Text}}</span>
									{{end}}
								{{end}}
							</p>
						</div>
						{{if .HasScreenshot}}
							<div class="screenshot-container">
								<a href="{{.ScreenshotPath}}" target="_blank"><img src="{{.ScreenshotPath}}" alt="" width="470" class="screenshot" /></a>
							</div>
						{{else}}
							<img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAdYAAAElBAMAAACrMt7eAAAAG1BMVEXMzMyWlpajo6PFxcW3t7ecnJyqqqq+vr6xsbGXmO98AAAACXBIWXMAAA7EAAAOxAGVKw4bAAADu0lEQVR4nO3XsY/aSBTH8WcMeEsmYNgSvMpdSvuUja6E4npY5M2WrE6RroRISU1Oyv3d994bO0BCGa1TfD8Slhke1vyYwTMWAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA+BmSub5KST/kR3//bvNYinzMfxepQnjyYxineghlU+slY60dW/mT+Icz+3L/lciXvX5HC1f6qqy8aU9CmPpXmku+cE6TTLz/92GTW7fSTbHZSj/UYSlVUSyOdizu0ld6LGNtLIlZ000etmlR5MWdXawXRA5HzbkUWYxOWa09eSg2uyarXbKLrLln3SxltdO3w61kj5Js5X6sPZWbkR31JxidamNJzDqYSH8qzXDph3qZZC3yrAnzySmrt+/sajHry8eM3XvWbpSZTsG+9WOgM/AvWemkvrUupeOLrF4bS2LWg55rcZt1ddjJUM9rvdpsespq7ZpVr9Zt1k9r7VBvp6dv7K12/k8/+8e7dJnVa2NJzPqbHnvHb1k/9NbSn0h2qxNk/fqU1dota9fj+va2GSv5V18D/x+9l6ZL2eUcjrVeErNa4c38W9bHm5mkU9GjDHZ6k2qzWrtmHa47zlrWNlbWpU/6ysKbfdt17dIfc6nquj6meV0/NbWxRBvqOCWzdfuFbGIj+lp0FOWw1JtUk9Xbk7wO+yarXbKTrKtjUvovfbDD55Af4yj60rD3486WlWlT6yUSfCmyQh37JuvNWGzRKXUI7Thvs3q7rjnheFpzdi+e1LL2xmfjql3/OmnHta7D/GIOe+2p5PtxHcxlUdqILkUepDdus3q7/gAfO5/D2fTs/2q+im0M/rMuDS//r14bS679X228lprsS+n7i0mb1ds1a9b1vamURRXvw3ZTrUrrim6XZHplzfHaWHJ2H961WVd1rduF4ajQJUf/z9M2q7cnu3YidJm1ej6tr8lS5F6HzZaOH9ccr40l19bXhU7XtWQzvUH1RrbKNlm9/VdYc0oZhrN903tJC6nu5LC+mlVrY8m1fZOe2XKT6zu/Pe2brN5u+6Z551nTcLYfDg+Lse+H996l7cV+2Gtjydl+uM2qU8G3EYt13BEfjpXtgou33q5ni73kthPubD9c2vrw43PO3/Hnf/Iby8wfZfZN7ffPOW1W2x3aomMTxHbEybzyleazt+uZPjv5E1HVPhcBAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAANCJ/wEokZhguJz55gAAAABJRU5ErkJggg==" class="no-screenshot" alt="" width="470" />
						{{end}}
						<div class="card-footer text-muted">
							<a href="#" class="card-link page-details-link">Details</a>
							<a href="{{.URL}}" target="_blank" class="card-link page-visit-link">Visit</a>
						</div>
						<div class="response-headers-container">
							<table class="table table-responsive table-striped table-hover table-sm response-headers">
								<thead class="thead-dark">
									<tr>
										<th scope="col">Header</th>
										<th scope="col">Value</th>
									</tr>
								</thead>
								<tbody>
									{{range .Headers}}
										{{if .IsInsecure}}
											<tr class="insecure">
										{{else if .IsSecure}}
											<tr class="secure">
										{{else}}
											<tr>
										{{end}}
											<td>{{.Name}}</td>
											<td>{{.Value}}</td>
										</tr>
									{{end}}
								</tbody>
							</table>
						</div>
					</div>
				{{end}}
			</div>
		{{end}}
		<footer>
			<p>AQUATONE v{{.Session.Version}} &middot; made with <span style="color:red;font-weight:bold">&#65533;</span> by <a href="https://michenriksen.com" target="_blank">Michael Henriksen</a></p>
			<p><a href="https://www.buymeacoffee.com/michenriksen" target="_blank"><img src=" data:image/png;base64,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" alt="Buy Me A Coffee" style="height: auto !important;width: auto !important;"></a></p>
		</footer>
		<div class="modal fade" tabindex="-1" role="dialog" aria-hidden="true" id="details_modal">
			<div class="modal-dialog modal-lg">
				<div class="modal-content">
					<div class="modal-header">
        		<h5 class="modal-title"></h5>
        		<button type="button" class="close" data-dismiss="modal" aria-label="Close">
          		<span aria-hidden="true">&times;</span>
        		</button>
      		</div>
      		<div class="modal-body"></div>
      		<div class="modal-footer">
        		<button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button>
      		</div>
				</div>
			</div>
		</div>
		<script src="https://code.jquery.com/jquery-3.3.1.slim.min.js" integrity="sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo" crossorigin="anonymous"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.3/umd/popper.min.js" integrity="sha384-ZMP7rVo3mIykV+2+9J3UJ46jBk0WLaUAdn689aCwoqbBJiSnjAK/l8WvCWPIPm49" crossorigin="anonymous"></script>
		<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.1.3/js/bootstrap.min.js" integrity="sha384-ChfqqxuZUCnJSK3+MXmPNIyE6ZbWh2IMqE241rYiqJxyMiZ6OW/JmZQ5stwEULTy" crossorigin="anonymous"></script>
		<script type="text/javascript">
			$(function() {
				$(".screenshot-container").on("mouseover", function() {
				  $(this).find("img.screenshot").css({"transform": "scale(2)"});
				}).on("mouseout", function() {
				  $(this).find("img.screenshot").css({"transform": "scale(1)"})
				}).on("mousemove", function(e) {
				  $(this).find("img.screenshot").css({'transform-origin': ((e.pageX - $(this).offset().left) / $(this).width()) * 100 + '% ' + ((e.pageY - $(this).offset().top) / $(this).height()) * 100 +'%'});
				});

				$(".page-details-link").on("click", function(e) {
					e.preventDefault();
					var page    = $(this).closest(".page");
					var url     = page.find("h5.card-title").text();
					var headers = page.find(".response-headers-container").html();
					$("#details_modal .modal-header h5").text(url);
					$("#details_modal .modal-body").html(headers);
					$("#details_modal").modal();
				});
			});
		</script>
	</body>
</html>`
)

func (p *Page) AddHeader(name string, value string) {
	p.Headers = append(p.Headers, Header{
		Name:  name,
		Value: value,
	})
}

type ReportData struct {
	Session  *Session
	Clusters [][]Page
}

type Report struct {
	Data ReportData
}

func (r *Report) Render(dest io.Writer) error {
	tmpl, err := template.New("Aquatone Report").Parse(Template)
	if err != nil {
		return err
	}
	err = tmpl.Execute(dest, r.Data)
	if err != nil {
		return err
	}
	return nil
}

func NewReport(data ReportData) *Report {
	return &Report{
		Data: data,
	}
}

func NewCluster(urls []*ResponsiveURL, session *Session) ([]Page, error) {
	var cluster []Page
	for _, url := range urls {
		page, err := NewPage(url, session)
		if err != nil {
			continue
		}
		cluster = append(cluster, page)
	}

	return cluster, nil
}

func NewPage(url *ResponsiveURL, session *Session) (Page, error) {
	baseFilename := session.BaseFilenameFromURL(url.URL)
	page := Page{
		URL:            url.URL,
		HeadersPath:    fmt.Sprintf("headers/%s.txt", baseFilename),
		BodyPath:       fmt.Sprintf("html/%s.html", baseFilename),
		ScreenshotPath: fmt.Sprintf("screenshots/%s.png", baseFilename),
		Tags:           url.Tags,
		Notes:          url.Notes,
	}
	contents, err := session.ReadFile(fmt.Sprintf("headers/%s.txt", baseFilename))
	if err != nil {
		return page, err
	}

	if _, err := os.Stat(session.GetFilePath(page.ScreenshotPath)); os.IsNotExist(err) {
		page.HasScreenshot = false
	} else {
		page.HasScreenshot = true
	}

	lines := strings.Split(string(contents), "\n")
	status, headers := lines[0], lines[1:]
	page.Status = status
	for _, header := range headers {
		h := strings.Split(header, ": ")
		if len(h) < 2 {
			continue
		}
		name, value := h[0], strings.Join(h[1:], ": ")
		page.AddHeader(name, value)
	}
	return page, nil
}
