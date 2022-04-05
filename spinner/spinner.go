package spinner

import (
	"io"
	"time"

	"github.com/briandowns/spinner"
)

// Option is a function that takes a spinner and applies a given configuration.
type Option func(*options)

// Start run spinner with the supplied options.
func Start(d time.Duration, opts ...Option) func() {
	opt := newOptions(opts...)
	s := spinner.New(opt.charSet(), d, opt.opts...)
	s.Start()
	return s.Stop
}

type options struct {
	predefinedSet int // Default: 11
	customSet     []string
	opts          []spinner.Option
}

func (o *options) charSet() []string {
	if 0 < len(o.customSet) {
		return o.customSet
	}

	return spinner.CharSets[o.predefinedSet]
}

// WithPredefinedSet adds the given predefined charsets the spinner.
// https://github.com/briandowns/spinner#available-character-sets
func WithPredefinedSet(n int) Option {
	return func(o *options) {
		o.predefinedSet = n
	}
}

// WithCustomSet adds the given custom charsets the spinner.
func WithCustomSet(set []string) Option {
	return func(o *options) {
		o.customSet = set
	}
}

// WithColor adds the given color to the spinner.
func WithColor(color string) Option {
	return func(o *options) {
		o.opts = append(o.opts, spinner.WithColor(color))
	}
}

// WithSuffix adds the given string to the spinner as the suffix.
func WithSuffix(suffix string) Option {
	return func(o *options) {
		o.opts = append(o.opts, spinner.WithSuffix(suffix))
	}
}

// WithFinalMSG adds the given string ot the spinner as the final message to be written.
func WithFinalMSG(msg string) Option {
	return func(o *options) {
		o.opts = append(o.opts, spinner.WithFinalMSG(msg))
	}
}

// WithHiddenCursor hides the cursor if hideCursor = true given.
func WithHiddenCursor(hide bool) Option {
	return func(o *options) {
		o.opts = append(o.opts, spinner.WithHiddenCursor(hide))
	}
}

// WithWriter adds the given writer to the spinner.
// This function should be favored over directly assigning to the struct value.
func WithWriter(w io.Writer) Option {
	return func(o *options) {
		o.opts = append(o.opts, spinner.WithWriter(w))
	}
}

func newOptions(opts ...Option) *options {
	ret := &options{
		predefinedSet: 11,
	}
	for _, v := range opts {
		v(ret)
	}
	return ret
}
