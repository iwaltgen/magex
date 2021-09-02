package archive

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bitfield/script"
	"github.com/stretchr/testify/assert"
)

func TestUnarchive(t *testing.T) {
	files, err := script.ListFiles("testdata").Slice()
	assert.NoError(t, err)

	for _, v := range files {
		src := v
		t.Run(filepath.Base(src), func(t *testing.T) {
			// given
			dest, err := os.MkdirTemp("", "magex-test-*")
			assert.NoError(t, err)
			fmt.Println(dest)
			defer os.RemoveAll(dest)

			// when
			err = Unarchive(src, dest)

			// then
			if strings.HasSuffix(src, "7z") {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			n, err := script.FindFiles(dest).CountLines()
			assert.NoError(t, err)
			assert.Equal(t, 5, n)
		})
	}
}
