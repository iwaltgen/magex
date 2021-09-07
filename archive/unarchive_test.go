package archive

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bitfield/script"
	"github.com/stretchr/testify/assert"
)

func TestUnarchive(t *testing.T) {
	files, err := script.ListFiles("testdata").EachLine(func(s string, b *strings.Builder) {
		if strings.Contains(s, "invalid") {
			return
		}
		b.WriteString(s + "\n")
	}).Slice()
	assert.NoError(t, err)

	for _, v := range files {
		src := v
		t.Run(filepath.Base(src), func(t *testing.T) {
			// given
			dest, err := os.MkdirTemp("", "magex-test-*")
			assert.NoError(t, err)
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

func TestUnarchivePermissionError(t *testing.T) {
	// given
	files, err := script.ListFiles("testdata").EachLine(func(s string, b *strings.Builder) {
		if strings.Contains(s, "invalid") {
			return
		}
		b.WriteString(s + "\n")
	}).Last(2).Slice()
	assert.NoError(t, err)

	for _, v := range files {
		src := v
		t.Run(filepath.Base(src), func(t *testing.T) {
			// given
			dest := "/home/unknown/magex"

			// when
			err = Unarchive(src, dest)

			// then
			assert.Error(t, err)
		})
	}
}

func TestUnarchiveOpenFileError(t *testing.T) {
	files := []string{
		"testdata/archive-nf.tgz",
		"testdata/archive-nf.zip",
	}
	for _, v := range files {
		src := v
		t.Run(filepath.Base(src), func(t *testing.T) {
			// given
			dest, err := os.MkdirTemp("", "magex-test-*")
			assert.NoError(t, err)
			defer os.RemoveAll(dest)

			// when
			err = Unarchive(src, dest)

			// then
			assert.Error(t, err)
		})
	}
}

func TestUnarchiveOpenReaderError(t *testing.T) {
	files, err := script.ListFiles("testdata").EachLine(func(s string, b *strings.Builder) {
		if !strings.Contains(s, "invalid") {
			return
		}
		b.WriteString(s + "\n")
	}).Slice()
	assert.NoError(t, err)

	for _, v := range files {
		src := v
		t.Run(filepath.Base(src), func(t *testing.T) {
			// given
			dest, err := os.MkdirTemp("", "magex-test-*")
			assert.NoError(t, err)
			defer os.RemoveAll(dest)

			// when
			err = Unarchive(src, dest)

			// then
			assert.Error(t, err)
		})
	}
}
