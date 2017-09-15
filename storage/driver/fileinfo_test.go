package driver

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFileInfo(t *testing.T) {
	var fif FileInfoFields
	fif.Path = "test"
	fif.Size = 1024
	fif.ModTime = time.Now()
	fif.IsDir = true

	var fii FileInfoInternal
	fii.FileInfoFields = fif

	assert.Equal(t, fif.Path, fii.Path())
	assert.Equal(t, fif.Size, fii.Size())
	assert.Equal(t, fif.ModTime, fii.ModTime())
	assert.Equal(t, fif.IsDir, fii.IsDir())
}
