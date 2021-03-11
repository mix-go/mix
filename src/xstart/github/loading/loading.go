package loading

import (
    "errors"
    "fmt"
    "io"
    "sync"
    "time"

    "github.com/fatih/color"
)

const (
	// 100ms per frame
	DEFAULT_FRAME_RATE = time.Millisecond * 100
)

var errColorNotFound = errors.New("color not found")

// DefaultLoaders is the
// default styles of loading
var DefaultLoaders = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}

// allowableColors is a collection of allowed colors.
var allowableColors = map[string]bool{
	"red":   true,
	"green": true,
	"blue":  true,
	"white": true, // The default color
	"black": true,
}

// foregroundColorAttribute is a collection
// of foreground attributes
var foregroundColorAttribute = map[string]color.Attribute{
	"red":   color.FgRed,
	"green": color.FgGreen,
	"blue":  color.FgBlue,
	"white": color.FgWhite,
	"black": color.FgBlack,
}

// IsColorAllowed returns boolean if a color is allowed.
func IsColorAllowed(color string) bool {
	allowed := false
	if allowableColors[color] {
		allowed = true
	}
	return allowed
}

// Loading is our blueprint for loading.
type Loading struct {
	sync.Mutex
	Title     string
	Loaders   []string
	FrameRate time.Duration
	runChan   chan struct{}
	stopOnce  sync.Once
	Output    io.Writer
	NoTty     bool
	Color     func(a ...interface{}) string
}

// NewLoading returns an instance of Loading
func NewLoading(title string) *Loading {
	loading := &Loading{
		Title:     title,
		Loaders:   DefaultLoaders,
		FrameRate: DEFAULT_FRAME_RATE,
		runChan:   make(chan struct{}),
		Color:     color.New(color.FgWhite).SprintFunc(),
	}
	return loading
}

// StartNew starts a new loader
func StartNew(title string) *Loading {
	return NewLoading(title).Start()
}

// Start starts the loading
func (loading *Loading) Start() *Loading {
	go loading.writer()
	return loading
}

// SetColor sets the color of the loading
func (loading *Loading) SetColor(c string) error {
	if !IsColorAllowed(c) {
		return errColorNotFound
	}
	loading.Color = color.New(foregroundColorAttribute[c]).SprintFunc()
	return nil
}

// SetSpeed sets the speed of our loader
func (loading *Loading) SetSpeed(rate time.Duration) *Loading {
	loading.Lock()
	loading.FrameRate = rate
	loading.Unlock()
	return loading
}

// SetLoaders set the character loader of the loading
func (loading *Loading) SetLoaders(chars []string) *Loading {
	loading.Lock()
	loading.Loaders = chars
	loading.Unlock()
	return loading
}

// Stop stops and clears the loading
func (loading *Loading) Stop() {
	loading.stopOnce.Do(func() {
		close(loading.runChan)
		loading.clear()
	})
}

// Restart restarts our loading
func (loading *Loading) Restart() {
	loading.Stop()
	loading.Start()
}

// animates our loader
func (loading *Loading) animate() {
	var out string
	for i := 0; i < len(loading.Loaders); i++ {
		out = " " + loading.Color(loading.Loaders[i]) + " " + loading.Title
		switch {
		case loading.Output != nil:
			fmt.Fprint(loading.Output, out)
		case !loading.NoTty:
			fmt.Print(out)
		}
		time.Sleep(loading.FrameRate)
		loading.clear()
	}
}

// Writes to our channel
func (loading *Loading) writer() {
	loading.animate()
	for {
		select {
		case <-loading.runChan:
			return
		default:
			loading.animate()
		}
	}
}

// Clears the screen
func (loading *Loading) clear() {
	if !loading.NoTty {
		fmt.Printf("\033[2K")
		fmt.Println()
		fmt.Printf("\033[1A")
	}
}
