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
			// expect
			err := ExecStdout("touch " + v)
			assert.NoError(t, err)

			err = ExecStdout("rm -f " + v)
			assert.NoError(t, err)
		})
	}
}

func TestListFile(t *testing.T) {
	// when
	ret, err := Exec("ls -l", "wc -l").String()

	// then
	assert.NoError(t, err)
	assert.NotEmpty(t, ret)
}

func TestStdoutEmptyCmd(t *testing.T) {
	// when
	err := ExecStdout()

	// then
	assert.NoError(t, err)
}
