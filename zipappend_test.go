package gozu

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestZipAppend_AppendsContentAtSpecifiedPath(t *testing.T) {
	tempDir := t.TempDir()
	baseDir := filepath.Join(tempDir, "base")
	require.NoError(t, os.MkdirAll(baseDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(baseDir, "base.txt"), []byte("base"), 0644))

	baseZip, err := ZipToBytes(baseDir, AllowAll)
	require.NoError(t, err)

	payloadDir := filepath.Join(tempDir, "payload")
	require.NoError(t, os.MkdirAll(filepath.Join(payloadDir, "nested"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(payloadDir, "new.txt"), []byte("new"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(payloadDir, "nested", "deep.txt"), []byte("deep"), 0644))

	payloadZip, err := ZipToBytes(payloadDir, AllowAll)
	require.NoError(t, err)

	merged, err := ZipAppend(baseZip, "assets/content", payloadZip)
	require.NoError(t, err)

	extractDir := filepath.Join(tempDir, "extract")
	require.NoError(t, UnzipFromBytes(merged, extractDir, AllowAll))

	content, err := os.ReadFile(filepath.Join(extractDir, "base.txt"))
	require.NoError(t, err)
	require.Equal(t, "base", string(content))

	content, err = os.ReadFile(filepath.Join(extractDir, "assets", "content", "new.txt"))
	require.NoError(t, err)
	require.Equal(t, "new", string(content))

	content, err = os.ReadFile(filepath.Join(extractDir, "assets", "content", "nested", "deep.txt"))
	require.NoError(t, err)
	require.Equal(t, "deep", string(content))
}

func TestZipAppend_CreatesArchiveWhenSourceEmpty(t *testing.T) {
	tempDir := t.TempDir()
	payloadDir := filepath.Join(tempDir, "payload")
	require.NoError(t, os.MkdirAll(payloadDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(payloadDir, "only.txt"), []byte("only"), 0644))

	payloadZip, err := ZipToBytes(payloadDir, AllowAll)
	require.NoError(t, err)

	merged, err := ZipAppend(nil, "/assets//nested/../final/", payloadZip)
	require.NoError(t, err)

	extractDir := filepath.Join(tempDir, "extract")
	require.NoError(t, UnzipFromBytes(merged, extractDir, AllowAll))

	content, err := os.ReadFile(filepath.Join(extractDir, "assets", "final", "only.txt"))
	require.NoError(t, err)
	require.Equal(t, "only", string(content))
}

func TestZipAppend_InvalidZippedContent(t *testing.T) {
	tempDir := t.TempDir()
	baseDir := filepath.Join(tempDir, "base")
	require.NoError(t, os.MkdirAll(baseDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(baseDir, "base.txt"), []byte("base"), 0644))

	baseZip, err := ZipToBytes(baseDir, AllowAll)
	require.NoError(t, err)

	_, err = ZipAppend(baseZip, "", []byte("not a zip"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "append new archive entries")
}

func TestZipAppend_InvalidTargetPath(t *testing.T) {
	tempDir := t.TempDir()
	payloadDir := filepath.Join(tempDir, "payload")
	require.NoError(t, os.MkdirAll(payloadDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(payloadDir, "file.txt"), []byte("data"), 0644))

	payloadZip, err := ZipToBytes(payloadDir, AllowAll)
	require.NoError(t, err)

	_, err = ZipAppend(nil, "../invalid", payloadZip)
	require.Error(t, err)
	require.Contains(t, err.Error(), "sanitize target path")
}

func TestZipAppend_RejectsTraversalEntries(t *testing.T) {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)
	_, err := zipWriter.Create("../evil.txt")
	require.NoError(t, err)
	require.NoError(t, zipWriter.Close())

	_, err = ZipAppend(nil, "", buf.Bytes())
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid archive entry path")
}

func TestZipAppend_EmptyZippedContent(t *testing.T) {
	_, err := ZipAppend(nil, "", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "zipped content is empty")
}
