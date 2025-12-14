package gozu

import (
	"io/fs"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// mockFileInfo implements fs.FileInfo for testing
type mockFileInfo struct {
	name string
	size int64
}

func (m mockFileInfo) Name() string       { return m.name }
func (m mockFileInfo) Size() int64        { return m.size }
func (m mockFileInfo) Mode() fs.FileMode  { return 0 }
func (m mockFileInfo) ModTime() time.Time { return time.Now() }
func (m mockFileInfo) IsDir() bool        { return false }
func (m mockFileInfo) Sys() any           { return nil }

func TestAllowAll(t *testing.T) {
	filter := AllowAll
	info := mockFileInfo{name: "test.txt"}
	require.True(t, filter("path/to/test.txt", info), "AllowAll should return true for any file")
}

func TestCombineFilters(t *testing.T) {
	allowTxt := func(path string, info fs.FileInfo) bool {
		return info.Name() == "test.txt"
	}
	allowSmall := func(path string, info fs.FileInfo) bool {
		return info.Size() < 100
	}
	combined := CombineFilters(allowTxt, allowSmall)

	// Test case: matches both filters
	info := mockFileInfo{name: "test.txt"}
	require.True(t, combined("path/to/test.txt", info), "Combined filter should return true when both filters pass")

	// Test case: fails one filter
	infoFail := mockFileInfo{name: "test.jpg"}
	require.False(t, combined("path/to/test.jpg", infoFail), "Combined filter should return false when any filter fails")
}

func TestNewFilterFuncFromPathOnly(t *testing.T) {
	pathOnly := func(path string) bool {
		return strings.HasSuffix(path, ".go")
	}
	filter := NewFilterFunc(pathOnly)
	info := mockFileInfo{name: "example.go", size: 1}
	require.True(t, filter("example.go", info))
	require.False(t, filter("example.txt", info))
}

func TestNewFilterFuncFromInfoOnly(t *testing.T) {
	infoOnly := func(info fs.FileInfo) bool {
		return info.Size() == 0
	}
	filter := NewFilterFunc(infoOnly)
	info := mockFileInfo{name: "empty", size: 0}
	require.True(t, filter("empty", info))
	infoFull := mockFileInfo{name: "data", size: 10}
	require.False(t, filter("data", infoFull))
}

func TestNewFilterFuncFromFilterFunc(t *testing.T) {
	original := FilterFunc(func(path string, info fs.FileInfo) bool {
		return strings.HasPrefix(path, "docs/") && info.Size() > 0
	})
	filter := NewFilterFunc(original)
	info := mockFileInfo{name: "doc", size: 5}
	require.True(t, filter("docs/readme", info))
	require.False(t, filter("other/readme", info))
}
