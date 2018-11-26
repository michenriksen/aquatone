package core

import (
	"fmt"
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
	StartedAt            time.Time
	FinishedAt           time.Time
	PortOpen             uint32
	PortClosed           uint32
	RequestSuccessful    uint32
	RequestFailed        uint32
	ResponseCode2xx      uint32
	ResponseCode3xx      uint32
	ResponseCode4xx      uint32
	ResponseCode5xx      uint32
	ScreenshotSuccessful uint32
	ScreenshotFailed     uint32
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

type Tag struct {
	Text string
	Type string
	Link string
}

func (t *Tag) HasLink() bool {
	if t.Link != "" {
		return true
	}
	return false
}

type Note struct {
	Text string
	Type string
}

type ResponsiveURL struct {
	URL   string
	Tags  []Tag
	Notes []Note
}

func (u *ResponsiveURL) AddTag(t Tag) {
	u.Tags = append(u.Tags, t)
}

func (u *ResponsiveURL) AddNote(n Note) {
	u.Notes = append(u.Notes, n)
}

type Session struct {
	sync.Mutex
	Version        string
	Options        Options `json:"-"`
	Out            *Logger `json:"-"`
	Stats          *Stats
	ResponsiveURLs map[string]*ResponsiveURL
	Ports          []int
	EventBus       EventBus.Bus                  `json:"-"`
	WaitGroup      sizedwaitgroup.SizedWaitGroup `json:"-"`
	WaitGroup2     sizedwaitgroup.SizedWaitGroup `json:"-"`
}

func (s *Session) Start() {
	s.ResponsiveURLs = make(map[string]*ResponsiveURL)
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

func (s *Session) AddResponsiveURL(url string) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.ResponsiveURLs[url]; ok {
		return
	}
	s.ResponsiveURLs[url] = &ResponsiveURL{URL: url}
}

func (s *Session) AddTagToResponsiveURL(url string, t string, tagType string, link string) {
	s.Lock()
	defer s.Unlock()
	if u, ok := s.ResponsiveURLs[url]; ok {
		u.AddTag(Tag{
			Text: t,
			Type: tagType,
			Link: link,
		})
	}
}

func (s *Session) AddNoteToResponsiveURL(url string, text string, noteType string) {
	s.Lock()
	defer s.Unlock()
	if u, ok := s.ResponsiveURLs[url]; ok {
		u.AddNote(Note{
			Text: text,
			Type: noteType,
		})
	}
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
	s.WaitGroup2 = sizedwaitgroup.New(*s.Options.Threads)
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
	host := strings.Replace(u.Host, ":", "__", 1)
	filename := fmt.Sprintf("%s__%s", u.Scheme, strings.Replace(host, ".", "_", -1))
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

	outdir := filepath.Clean(*session.Options.OutDir)
	session.Options.OutDir = &outdir

	session.Version = Version
	session.Start()

	return &session, nil
}
