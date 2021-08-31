package script

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTouchRemove(t *testing.T) {
	// dataset
	dataset := []string{"test1", "test2", "test3"}

	// table driven tests
	for _, v := range dataset {
		t.Run(v, func(t *testing.T) {
			// when
			cmds := []string{
				"touch " + v,
				"rm -f " + v,
			}
			err := ExecStdout(cmds...)

			// then
			assert.NoError(t, err)
		})
	}
}
