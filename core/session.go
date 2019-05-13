package core

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/asaskevich/EventBus"
	"github.com/remeh/sizedwaitgroup"
)

type Stats struct {
	StartedAt            time.Time `json:"startedAt"`
	FinishedAt           time.Time `json:"finishedAt"`
	PortOpen             uint32    `json:"portOpen"`
	PortClosed           uint32    `json:"portClosed"`
	RequestSuccessful    uint32    `json:"requestSuccessful"`
	RequestFailed        uint32    `json:"requestFailed"`
	ResponseCode2xx      uint32    `json:"responseCode2xx"`
	ResponseCode3xx      uint32    `json:"responseCode3xx"`
	ResponseCode4xx      uint32    `json:"responseCode4xx"`
	ResponseCode5xx      uint32    `json:"responseCode5xx"`
	ScreenshotSuccessful uint32    `json:"screenshotSuccessful"`
	ScreenshotFailed     uint32    `json:"screenshotFailed"`
}

func (s *Stats) Duration() time.Duration {
	return s.FinishedAt.Sub(s.StartedAt)
}

func (s *Stats) IncrementPortOpen() {
	atomic.AddUint32(&s.PortOpen, 1)
}

func (s *Stats) IncrementPortClosed() {
	atomic.AddUint32(&s.PortClosed, 1)
}

func (s *Stats) IncrementRequestSuccessful() {
	atomic.AddUint32(&s.RequestSuccessful, 1)
}

func (s *Stats) IncrementRequestFailed() {
	atomic.AddUint32(&s.RequestFailed, 1)
}

func (s *Stats) IncrementResponseCode2xx() {
	atomic.AddUint32(&s.ResponseCode2xx, 1)
}

func (s *Stats) IncrementResponseCode3xx() {
	atomic.AddUint32(&s.ResponseCode3xx, 1)
}

func (s *Stats) IncrementResponseCode4xx() {
	atomic.AddUint32(&s.ResponseCode4xx, 1)
}

func (s *Stats) IncrementResponseCode5xx() {
	atomic.AddUint32(&s.ResponseCode5xx, 1)
}

func (s *Stats) IncrementScreenshotSuccessful() {
	atomic.AddUint32(&s.ScreenshotSuccessful, 1)
}

func (s *Stats) IncrementScreenshotFailed() {
	atomic.AddUint32(&s.ScreenshotFailed, 1)
}

type Session struct {
	sync.Mutex
	Version                string                        `json:"version"`
	Options                Options                       `json:"-"`
	Out                    *Logger                       `json:"-"`
	Stats                  *Stats                        `json:"stats"`
	Pages                  map[string]*Page              `json:"pages"`
	PageSimilarityClusters map[string][]string           `json:"pageSimilarityClusters"`
	Ports                  []int                         `json:"-"`
	EventBus               EventBus.Bus                  `json:"-"`
	WaitGroup              sizedwaitgroup.SizedWaitGroup `json:"-"`
}

func (s *Session) Start() {
	s.Pages = make(map[string]*Page)
	s.PageSimilarityClusters = make(map[string][]string)
	s.initStats()
	s.initLogger()
	s.initPorts()
	s.initThreads()
	s.initEventBus()
	s.initWaitGroup()
	s.initDirectories()
}

func (s *Session) End() {
	s.Stats.FinishedAt = time.Now()
}

func (s *Session) AddPage(url string) (*Page, error) {
	s.Lock()
	defer s.Unlock()
	if page, ok := s.Pages[url]; ok {
		return page, nil
	}

	page, err := NewPage(url)
	if err != nil {
		return nil, err
	}

	s.Pages[url] = page
	return page, nil
}

func (s *Session) GetPage(url string) *Page {
	if page, ok := s.Pages[url]; ok {
		return page
	}
	return nil
}

func (s *Session) GetPageByUUID(id string) *Page {
	for _, page := range s.Pages {
		if page.UUID == id {
			return page
		}
	}
	return nil
}

func (s *Session) initStats() {
	if s.Stats != nil {
		return
	}
	s.Stats = &Stats{
		StartedAt: time.Now(),
	}
}

func (s *Session) initPorts() {
	var ports []int
	switch *s.Options.Ports {
	case "small":
		ports = SmallPortList
	case "", "medium", "default":
		ports = MediumPortList
	case "large":
		ports = LargePortList
	case "xlarge", "huge":
		ports = XLargePortList
	default:
		for _, p := range strings.Split(*s.Options.Ports, ",") {
			port, err := strconv.Atoi(strings.TrimSpace(p))
			if err != nil {
				s.Out.Fatal("Invalid port range given\n")
				os.Exit(1)
			}
			if port < 1 || port > 65535 {
				s.Out.Fatal("Invalid port given: %v\n", port)
				os.Exit(1)
			}
			ports = append(ports, port)
		}
	}
	s.Ports = ports
}

func (s *Session) initLogger() {
	s.Out = &Logger{}
	s.Out.SetDebug(*s.Options.Debug)
	s.Out.SetSilent(*s.Options.Silent)
}

func (s *Session) initThreads() {
	if *s.Options.Threads == 0 {
		numCPUs := runtime.NumCPU()
		s.Options.Threads = &numCPUs
	}
}

func (s *Session) initEventBus() {
	s.EventBus = EventBus.New()
}

func (s *Session) initWaitGroup() {
	s.WaitGroup = sizedwaitgroup.New(*s.Options.Threads)
}

func (s *Session) initDirectories() {
	for _, d := range []string{"headers", "html", "screenshots"} {
		d = s.GetFilePath(d)
		if _, err := os.Stat(d); os.IsNotExist(err) {
			err = os.MkdirAll(d, 0755)
			if err != nil {
				s.Out.Fatal("Failed to create required directory %s\n", d)
				os.Exit(1)
			}
		}
	}
}

func (s *Session) BaseFilenameFromURL(stru string) string {
	u, err := url.Parse(stru)
	if err != nil {
		return ""
	}

	h := sha1.New()
	io.WriteString(h, u.Path)
	io.WriteString(h, u.Fragment)

	pathHash := fmt.Sprintf("%x", h.Sum(nil))[0:16]
	host := strings.Replace(u.Host, ":", "__", 1)
	filename := fmt.Sprintf("%s__%s__%s", u.Scheme, strings.Replace(host, ".", "_", -1), pathHash)
	return strings.ToLower(filename)
}

func (s *Session) GetFilePath(p string) string {
	return path.Join(*s.Options.OutDir, p)
}

func (s *Session) ReadFile(p string) ([]byte, error) {
	content, err := ioutil.ReadFile(s.GetFilePath(p))
	if err != nil {
		return content, err
	}
	return content, nil
}

func (s *Session) ToJSON() string {
	sessionJSON, _ := json.Marshal(s)
	return string(sessionJSON)
}

func (s *Session) SaveToFile(filename string) error {
	path := s.GetFilePath(filename)
	err := ioutil.WriteFile(path, []byte(s.ToJSON()), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (s *Session) Asset(name string) ([]byte, error) {
	return Asset(name)
}

func NewSession() (*Session, error) {
	var err error
	var session Session

	session.Version = Version

	if session.Options, err = ParseOptions(); err != nil {
		return nil, err
	}

	if *session.Options.ChromePath != "" {
		if _, err := os.Stat(*session.Options.ChromePath); os.IsNotExist(err) {
			return nil, fmt.Errorf("Chrome path %s does not exist", *session.Options.ChromePath)
		}
	}

	if *session.Options.SessionPath != "" {
		if _, err := os.Stat(*session.Options.SessionPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("Session path %s does not exist", *session.Options.SessionPath)
		}
	}

	if *session.Options.TemplatePath != "" {
		if _, err := os.Stat(*session.Options.TemplatePath); os.IsNotExist(err) {
			return nil, fmt.Errorf("Template path %s does not exist", *session.Options.TemplatePath)
		}
	}

	envOutPath := os.Getenv("AQUATONE_OUT_PATH")
	if *session.Options.OutDir == "." && envOutPath != "" {
		session.Options.OutDir = &envOutPath
	}

	outdir := filepath.Clean(*session.Options.OutDir)
	session.Options.OutDir = &outdir

	session.Version = Version
	session.Start()

	return &session, nil
}
