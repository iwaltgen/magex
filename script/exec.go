package script

import (
	"os"

	"github.com/bitfield/script"
)

// Exec runs an external command and returns a pipe containing the output.
func Exec(cmds ...string) *Pipe {
	if len(cmds) == 0 {
		return nil
	}

	pipe := script.NewPipe().WithStderr(os.Stderr)
	for _, cmd := range cmds {
		pipe = pipe.Exec(os.ExpandEnv(cmd))
	}
	return pipe
}

// ExecStdout runs an external command and writes the contents of the pipe
// to the program's standard output.
func ExecStdout(cmds ...string) error {
	pipe := Exec(cmds...)
	if pipe == nil {
		return nil
	}

	_, err := pipe.Stdout()
	return err
}
