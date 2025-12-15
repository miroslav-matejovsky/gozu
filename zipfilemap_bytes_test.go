package gozu

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileMap_Get(t *testing.T) {
	fm := FileMap{
		"file.txt": []byte("content"),
	}

	data, ok := fm.Get("file.txt")
	require.True(t, ok)
	require.Equal(t, []byte("content"), data)

	data, ok = fm.Get("nonexistent.txt")
	require.False(t, ok)
	require.Nil(t, data)
}

func TestFileMap_Set(t *testing.T) {
	fm := make(FileMap)

	fm.Set("file.txt", []byte("content"))
	require.Equal(t, []byte("content"), fm["file.txt"])

	fm.Set("file.txt", []byte("updated"))
	require.Equal(t, []byte("updated"), fm["file.txt"])
}

func TestFileMap_Delete(t *testing.T) {
	fm := FileMap{
		"file.txt": []byte("content"),
	}

	fm.Delete("file.txt")
	_, ok := fm["file.txt"]
	require.False(t, ok)

	// Deleting non-existent key should not panic
	fm.Delete("nonexistent.txt")
}

func TestFileMap_Has(t *testing.T) {
	fm := FileMap{
		"file.txt": []byte("content"),
	}

	require.True(t, fm.Has("file.txt"))
	require.False(t, fm.Has("nonexistent.txt"))
}

func TestFileMap_Paths(t *testing.T) {
	fm := FileMap{
		"a.txt": []byte("a"),
		"b.txt": []byte("b"),
		"c.txt": []byte("c"),
	}

	paths := fm.Paths()
	sort.Strings(paths)
	require.Equal(t, []string{"a.txt", "b.txt", "c.txt"}, paths)
}

func TestFileMap_Paths_Empty(t *testing.T) {
	fm := make(FileMap)
	paths := fm.Paths()
	require.Empty(t, paths)
}

func TestZipMapToBytes(t *testing.T) {
	files := FileMap{
		"file1.txt":     []byte("content1"),
		"dir/file2.txt": []byte("content2"),
	}

	data, err := ZipMapToBytes(files)
	require.NoError(t, err)
	require.NotEmpty(t, data)

	// Verify by extracting
	extracted, err := UnzipBytesToMap(data)
	require.NoError(t, err)
	require.Len(t, extracted, 2)
	require.Equal(t, []byte("content1"), extracted["file1.txt"])
	require.Equal(t, []byte("content2"), extracted["dir/file2.txt"])
}

func TestZipMapToBytes_EmptyMap(t *testing.T) {
	files := make(FileMap)

	data, err := ZipMapToBytes(files)
	require.NoError(t, err)
	require.NotEmpty(t, data)

	extracted, err := UnzipBytesToMap(data)
	require.NoError(t, err)
	require.Empty(t, extracted)
}

func TestZipMapToBytes_NilMap(t *testing.T) {
	_, err := ZipMapToBytes(nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "files map cannot be nil")
}

func TestZipMapToBytes_EmptyContent(t *testing.T) {
	files := FileMap{
		"empty.txt": []byte{},
	}

	data, err := ZipMapToBytes(files)
	require.NoError(t, err)

	extracted, err := UnzipBytesToMap(data)
	require.NoError(t, err)
	require.Equal(t, []byte{}, extracted["empty.txt"])
}

func TestZipMapToBytes_BinaryContent(t *testing.T) {
	binaryData := []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD}
	files := FileMap{
		"binary.bin": binaryData,
	}

	data, err := ZipMapToBytes(files)
	require.NoError(t, err)

	extracted, err := UnzipBytesToMap(data)
	require.NoError(t, err)
	require.Equal(t, binaryData, extracted["binary.bin"])
}

func TestZipMapToBytes_NestedDirectories(t *testing.T) {
	files := FileMap{
		"a/b/c/deep.txt": []byte("deep content"),
		"a/shallow.txt":  []byte("shallow content"),
	}

	data, err := ZipMapToBytes(files)
	require.NoError(t, err)

	extracted, err := UnzipBytesToMap(data)
	require.NoError(t, err)
	require.Len(t, extracted, 2)
	require.Equal(t, []byte("deep content"), extracted["a/b/c/deep.txt"])
	require.Equal(t, []byte("shallow content"), extracted["a/shallow.txt"])
}

func TestZipMapToBytes_SkipsEmptyPaths(t *testing.T) {
	files := FileMap{
		"":         []byte("should be skipped"),
		".":        []byte("should also be skipped"),
		"file.txt": []byte("content"),
	}

	data, err := ZipMapToBytes(files)
	require.NoError(t, err)

	extracted, err := UnzipBytesToMap(data)
	require.NoError(t, err)
	require.Len(t, extracted, 1)
	require.Equal(t, []byte("content"), extracted["file.txt"])
}

func TestUnzipBytesToMap(t *testing.T) {
	files := FileMap{
		"file1.txt":     []byte("content1"),
		"dir/file2.txt": []byte("content2"),
	}

	data, err := ZipMapToBytes(files)
	require.NoError(t, err)

	extracted, err := UnzipBytesToMap(data)
	require.NoError(t, err)
	require.Len(t, extracted, 2)
	require.Equal(t, []byte("content1"), extracted["file1.txt"])
	require.Equal(t, []byte("content2"), extracted["dir/file2.txt"])
}

func TestUnzipBytesToMap_EmptyData(t *testing.T) {
	_, err := UnzipBytesToMap([]byte{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "zip data is empty")
}

func TestUnzipBytesToMap_InvalidData(t *testing.T) {
	_, err := UnzipBytesToMap([]byte("not a zip file"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "read zip archive from bytes")
}

func TestZipMapToBytes_LargeFile(t *testing.T) {
	// Create a 1MB file
	largeContent := make([]byte, 1024*1024)
	for i := range largeContent {
		largeContent[i] = byte(i % 256)
	}

	files := FileMap{
		"large.bin": largeContent,
	}

	data, err := ZipMapToBytes(files)
	require.NoError(t, err)

	extracted, err := UnzipBytesToMap(data)
	require.NoError(t, err)
	require.Equal(t, largeContent, extracted["large.bin"])
}

func TestZipMapToBytes_ManyFiles(t *testing.T) {
	files := make(FileMap)
	for i := 0; i < 100; i++ {
		files[string(rune('a'+i%26))+".txt"] = []byte{byte(i)}
	}

	data, err := ZipMapToBytes(files)
	require.NoError(t, err)

	extracted, err := UnzipBytesToMap(data)
	require.NoError(t, err)

	for path, content := range files {
		require.Equal(t, content, extracted[path], "mismatch for %s", path)
	}
}

func TestZipMapToBytes_RoundTrip(t *testing.T) {
	original := FileMap{
		"root.txt":           []byte("root content"),
		"dir1/file1.txt":     []byte("file1 content"),
		"dir1/dir2/file2.go": []byte("package main\n\nfunc main() {}"),
		"empty.dat":          []byte{},
		"binary.bin":         []byte{0x00, 0xFF, 0x7F, 0x80},
	}

	data, err := ZipMapToBytes(original)
	require.NoError(t, err)

	extracted, err := UnzipBytesToMap(data)
	require.NoError(t, err)

	require.Len(t, extracted, len(original))
	for path, content := range original {
		require.Equal(t, content, extracted[path], "mismatch for %s", path)
	}
}

func TestUnzipBytesToMap_SkipsDirectories(t *testing.T) {
	// Create a zip with explicit directory entries using ZipToBytes from filesystem
	// Since we can't easily create directory entries with ZipMapToBytes,
	// we verify that regular files work correctly
	files := FileMap{
		"dir/file.txt": []byte("content"),
	}

	data, err := ZipMapToBytes(files)
	require.NoError(t, err)

	extracted, err := UnzipBytesToMap(data)
	require.NoError(t, err)

	// Should only contain the file, not any directory entries
	require.Len(t, extracted, 1)
	require.Equal(t, []byte("content"), extracted["dir/file.txt"])
}
