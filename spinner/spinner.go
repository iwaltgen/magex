package spinner

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

// Start run spinner with the supplied options.
func Start(d time.Duration) func() {
	s := spinner.New(spinner.CharSets[11], d)
	s.Start()
	return func() {
		s.Stop()
		fmt.Fprintf(s.Writer, "\r%s", "")
	}
}
