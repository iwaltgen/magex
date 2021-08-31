package spinner

import (
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
)

// Start run spinner with the supplied options.
func Start(d time.Duration) func() {
	// TODO(iwaltgen): add options
	s := spinner.New(spinner.CharSets[11], d, spinner.WithWriter(os.Stderr))
	s.Start()
	return func() {
		s.Stop()
		fmt.Fprintf(s.Writer, "\r%s", "")
	}
}
