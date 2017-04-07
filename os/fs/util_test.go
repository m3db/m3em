package fs

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testingPrefixDir = "fs-testing-"
)

func newTempDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", testingPrefixDir)
	require.NoError(t, err)
	return dir
}

func newTempFile(t *testing.T, dir string, content []byte) *os.File {
	tmpfile, err := ioutil.TempFile(dir, "example")
	require.NoError(t, err)
	n, err := tmpfile.Write(content)
	require.NoError(t, err)
	require.Equal(t, len(content), n)
	require.NoError(t, tmpfile.Close())
	return tmpfile
}

func TestRemoveContents(t *testing.T) {
	var (
		content        = []byte("temporary file content")
		tmpdir         = newTempDir(t)
		numFilesToTest = 3
	)

	// create a random files
	for i := 0; i < numFilesToTest; i++ {
		tmpfile := newTempFile(t, tmpdir, content)
		tmpfile.Close()
	}
	files, err := ioutil.ReadDir(tmpdir)
	require.NoError(t, err)
	require.Equal(t, numFilesToTest, len(files))

	// remove contents of dir, list to make sure nothing is in there
	require.NoError(t, RemoveContents(tmpdir))
	fi, err := os.Stat(tmpdir)
	require.NoError(t, err)
	require.True(t, fi.IsDir())
	files, err = ioutil.ReadDir(tmpdir)
	require.NoError(t, err)
	require.Equal(t, 0, len(files))
}
