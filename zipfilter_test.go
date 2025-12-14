package gozu

import (
	"io/fs"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// mockFileInfo implements fs.FileInfo for testing
type mockFileInfo struct {
	name string
}

func (m mockFileInfo) Name() string       { return m.name }
func (m mockFileInfo) Size() int64        { return 0 }
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
