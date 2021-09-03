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

	pipe := script.NewPipe()
	for i, cmd := range cmds {
		pipe = pipe.Exec(os.ExpandEnv(cmd))
		if i < (len(cmds) - 1) {
			pipe.SetError(nil)
		}
	}
	return pipe
}

// ExecOutput runs an external command and returns the contents of the pipe.
func ExecOutput(cmds ...string) (string, error) {
	pipe := Exec(cmds...)
	if pipe == nil {
		return "", nil
	}

	return pipe.String()
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
