package gzu

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestZip(t *testing.T) {
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

	_, err = os.Stat(dstZip)
	require.NoError(t, err)

	extractDir := filepath.Join(tempDir, "extract")
	err = Unzip(dstZip, extractDir)
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

func TestZipToFile_NilFilter(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	err := os.MkdirAll(srcDir, 0755)
	require.NoError(t, err)

	dstZip := filepath.Join(tempDir, "test.zip")
	err = ZipToFile(srcDir, dstZip, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "filter function cannot be nil")
}

func TestZipToFile_NonExistentSource(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "nonexistent")
	dstZip := filepath.Join(tempDir, "test.zip")

	err := Zip(srcDir, dstZip)
	require.Error(t, err)
	require.Contains(t, err.Error(), "walk directory")
}

func TestZipToFile_InvalidDestination(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	err := os.MkdirAll(srcDir, 0755)
	require.NoError(t, err)

	// Try to create zip in non-existent directory
	dstZip := filepath.Join(tempDir, "nonexistent", "subdir", "test.zip")
	err = Zip(srcDir, dstZip)
	require.Error(t, err)
	require.Contains(t, err.Error(), "create zip file")
}

func TestZipToBytes(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	err := os.MkdirAll(srcDir, 0755)
	require.NoError(t, err)

	file1 := filepath.Join(srcDir, "file1.txt")
	err = os.WriteFile(file1, []byte("test content"), 0644)
	require.NoError(t, err)

	data, err := ZipToBytes(srcDir, AllowAll)
	require.NoError(t, err)
	require.NotEmpty(t, data)

	// Verify by extracting
	extractDir := filepath.Join(tempDir, "extract")
	err = UnzipFromBytes(data, extractDir, AllowAll)
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(extractDir, "file1.txt"))
	require.NoError(t, err)
	require.Equal(t, "test content", string(content))
}

func TestZipToBytes_NilFilter(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	err := os.MkdirAll(srcDir, 0755)
	require.NoError(t, err)

	_, err = ZipToBytes(srcDir, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "filter function cannot be nil")
}

func TestZipToFile_WithFilter(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	err := os.MkdirAll(srcDir, 0755)
	require.NoError(t, err)

	// Create multiple files with different extensions
	err = os.WriteFile(filepath.Join(srcDir, "include.txt"), []byte("included"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(srcDir, "exclude.log"), []byte("excluded"), 0644)
	require.NoError(t, err)

	// Filter to only include .txt files
	txtFilter := func(path string, info fs.FileInfo) bool {
		return strings.HasSuffix(path, ".txt")
	}

	dstZip := filepath.Join(tempDir, "filtered.zip")
	err = ZipToFile(srcDir, dstZip, txtFilter)
	require.NoError(t, err)

	extractDir := filepath.Join(tempDir, "extract")
	err = Unzip(dstZip, extractDir)
	require.NoError(t, err)

	// .txt file should exist
	_, err = os.Stat(filepath.Join(extractDir, "include.txt"))
	require.NoError(t, err)

	// .log file should not exist
	_, err = os.Stat(filepath.Join(extractDir, "exclude.log"))
	require.True(t, os.IsNotExist(err))
}

func TestZip_EmptyDirectory(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "empty")
	err := os.MkdirAll(srcDir, 0755)
	require.NoError(t, err)

	dstZip := filepath.Join(tempDir, "empty.zip")
	err = Zip(srcDir, dstZip)
	require.NoError(t, err)

	_, err = os.Stat(dstZip)
	require.NoError(t, err)
}

func TestZip_NestedDirectories(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")

	// Create deeply nested structure
	deepDir := filepath.Join(srcDir, "a", "b", "c", "d")
	err := os.MkdirAll(deepDir, 0755)
	require.NoError(t, err)

	deepFile := filepath.Join(deepDir, "deep.txt")
	err = os.WriteFile(deepFile, []byte("deep content"), 0644)
	require.NoError(t, err)

	dstZip := filepath.Join(tempDir, "nested.zip")
	err = Zip(srcDir, dstZip)
	require.NoError(t, err)

	extractDir := filepath.Join(tempDir, "extract")
	err = Unzip(dstZip, extractDir)
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(extractDir, "a", "b", "c", "d", "deep.txt"))
	require.NoError(t, err)
	require.Equal(t, "deep content", string(content))
}
