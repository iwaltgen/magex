package pipe

import "github.com/bitfield/script"

// Pipe represents a pipe object with an associated ReadAutoCloser.
type Pipe = script.Pipe

// Exec runs an external command and returns a pipe containing the output.
func Exec(cmds ...string) *Pipe {
	if len(cmds) == 0 {
		return nil
	}

	pipe := script.NewPipe()
	for _, cmd := range cmds {
		pipe = pipe.Exec(cmd)
		pipe.SetError(nil)
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
