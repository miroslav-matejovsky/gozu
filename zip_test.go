package gzu

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestZip(t *testing.T) {
	// Create temp dir with files
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	err := os.MkdirAll(srcDir, 0755)
	require.NoError(t, err)

	file1 := filepath.Join(srcDir, "file1.txt")
	err = os.WriteFile(file1, []byte("content1"), 0644)
	require.NoError(t, err)

	subDir := filepath.Join(srcDir, "sub")
	err = os.MkdirAll(subDir, 0755)
	require.NoError(t, err)

	file2 := filepath.Join(subDir, "file2.txt")
	err = os.WriteFile(file2, []byte("content2"), 0644)
	require.NoError(t, err)

	dstZip := filepath.Join(tempDir, "test.zip")
	err = Zip(srcDir, dstZip)
	require.NoError(t, err)

	// Check zip exists
	_, err = os.Stat(dstZip)
	require.NoError(t, err)

	// Unzip to check
	extractDir := filepath.Join(tempDir, "extract")
	err = Unzip(dstZip, extractDir)
	require.NoError(t, err)

	// Check files
	extracted1 := filepath.Join(extractDir, "file1.txt")
	content1, err := os.ReadFile(extracted1)
	require.NoError(t, err)
	require.Equal(t, "content1", string(content1))

	extracted2 := filepath.Join(extractDir, "sub", "file2.txt")
	content2, err := os.ReadFile(extracted2)
	require.NoError(t, err)
	require.Equal(t, "content2", string(content2))
}
