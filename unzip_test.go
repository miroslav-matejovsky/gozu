package gozu

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnzip(t *testing.T) {
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

	data, err := ZipToBytes(srcDir, AllowAll)
	require.NoError(t, err)

	extractDir := filepath.Join(tempDir, "extract")
	err = UnzipFromBytes(data, extractDir, AllowAll)
	require.NoError(t, err)

	extracted1 := filepath.Join(extractDir, "file1.txt")
	content1, err := os.ReadFile(extracted1)
	require.NoError(t, err)
	require.Equal(t, "content1", string(content1))

	extracted2 := filepath.Join(extractDir, "sub", "file2.txt")
	content2, err := os.ReadFile(extracted2)
	require.NoError(t, err)
	require.Equal(t, "content2", string(content2))
}

func TestUnzipFromFile_NonExistent(t *testing.T) {
	tempDir := t.TempDir()
	err := Unzip(filepath.Join(tempDir, "nonexistent.zip"), tempDir)
	require.Error(t, err)
	require.Contains(t, err.Error(), "open zip archive")
}

func TestUnzipFromFile_NilFilter(t *testing.T) {
	tempDir := t.TempDir()

	// Create a valid zip first
	srcDir := filepath.Join(tempDir, "src")
	err := os.MkdirAll(srcDir, 0755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(srcDir, "file.txt"), []byte("test"), 0644)
	require.NoError(t, err)

	zipPath := filepath.Join(tempDir, "test.zip")
	err = Zip(srcDir, zipPath)
	require.NoError(t, err)

	err = UnzipFromFile(zipPath, tempDir, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "filter function cannot be nil")
}

func TestUnzipFromBytes_NilFilter(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	err := os.MkdirAll(srcDir, 0755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(srcDir, "file.txt"), []byte("test"), 0644)
	require.NoError(t, err)

	data, err := ZipToBytes(srcDir, AllowAll)
	require.NoError(t, err)

	err = UnzipFromBytes(data, tempDir, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "filter function cannot be nil")
}

func TestUnzipFromBytes_EmptyData(t *testing.T) {
	tempDir := t.TempDir()
	err := UnzipFromBytes([]byte{}, tempDir, AllowAll)
	require.Error(t, err)
	require.Contains(t, err.Error(), "zip data is empty")
}

func TestUnzipFromBytes_InvalidData(t *testing.T) {
	tempDir := t.TempDir()
	err := UnzipFromBytes([]byte("not a zip file"), tempDir, AllowAll)
	require.Error(t, err)
	require.Contains(t, err.Error(), "read zip archive from bytes")
}

func TestUnzipFromFile_WithFilter(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	err := os.MkdirAll(srcDir, 0755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(srcDir, "include.txt"), []byte("included"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(srcDir, "exclude.log"), []byte("excluded"), 0644)
	require.NoError(t, err)

	zipPath := filepath.Join(tempDir, "test.zip")
	err = Zip(srcDir, zipPath)
	require.NoError(t, err)

	// Filter to only extract .txt files
	txtFilter := func(path string, info fs.FileInfo) bool {
		return strings.HasSuffix(path, ".txt")
	}

	extractDir := filepath.Join(tempDir, "extract")
	err = UnzipFromFile(zipPath, extractDir, txtFilter)
	require.NoError(t, err)

	// .txt file should exist
	_, err = os.Stat(filepath.Join(extractDir, "include.txt"))
	require.NoError(t, err)

	// .log file should not exist
	_, err = os.Stat(filepath.Join(extractDir, "exclude.log"))
	require.True(t, os.IsNotExist(err))
}

func TestUnzip_PreservesDirectoryStructure(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")

	// Create nested structure
	nestedDir := filepath.Join(srcDir, "level1", "level2", "level3")
	err := os.MkdirAll(nestedDir, 0755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(nestedDir, "deep.txt"), []byte("deep"), 0644)
	require.NoError(t, err)

	zipPath := filepath.Join(tempDir, "nested.zip")
	err = Zip(srcDir, zipPath)
	require.NoError(t, err)

	extractDir := filepath.Join(tempDir, "extract")
	err = Unzip(zipPath, extractDir)
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(extractDir, "level1", "level2", "level3", "deep.txt"))
	require.NoError(t, err)
	require.Equal(t, "deep", string(content))
}

func TestUnzip_ToExistingDirectory(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	err := os.MkdirAll(srcDir, 0755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(srcDir, "new.txt"), []byte("new content"), 0644)
	require.NoError(t, err)

	zipPath := filepath.Join(tempDir, "test.zip")
	err = Zip(srcDir, zipPath)
	require.NoError(t, err)

	// Create extraction directory with existing file
	extractDir := filepath.Join(tempDir, "extract")
	err = os.MkdirAll(extractDir, 0755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(extractDir, "existing.txt"), []byte("existing"), 0644)
	require.NoError(t, err)

	err = Unzip(zipPath, extractDir)
	require.NoError(t, err)

	// Both files should exist
	_, err = os.Stat(filepath.Join(extractDir, "new.txt"))
	require.NoError(t, err)

	_, err = os.Stat(filepath.Join(extractDir, "existing.txt"))
	require.NoError(t, err)
}

func TestUnzip_LargeFile(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	err := os.MkdirAll(srcDir, 0755)
	require.NoError(t, err)

	// Create a larger file (1MB)
	largeContent := make([]byte, 1024*1024)
	for i := range largeContent {
		largeContent[i] = byte(i % 256)
	}

	err = os.WriteFile(filepath.Join(srcDir, "large.bin"), largeContent, 0644)
	require.NoError(t, err)

	zipPath := filepath.Join(tempDir, "large.zip")
	err = Zip(srcDir, zipPath)
	require.NoError(t, err)

	extractDir := filepath.Join(tempDir, "extract")
	err = Unzip(zipPath, extractDir)
	require.NoError(t, err)

	extracted, err := os.ReadFile(filepath.Join(extractDir, "large.bin"))
	require.NoError(t, err)
	require.Equal(t, largeContent, extracted)
}
