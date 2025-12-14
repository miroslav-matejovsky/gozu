package gozu

import (
	"io/fs"
)

// FilterFunc is a function type used to filter files during zip operations.
// It takes a file path (relative to the zip root) and file info, returning true
// to include the file or false to exclude it. This allows selective inclusion
// of files based on custom criteria like path patterns, file types, or directory depth.
//
// Example usage:
//   - To include only .txt files: func(path string, info fs.FileInfo) bool { return filepath.Ext(path) == ".txt" }
//   - To exclude hidden files: func(path string, info fs.FileInfo) bool { return !strings.HasPrefix(filepath.Base(path), ".") }
//   - Combine with predefined filters for common cases.
type FilterFunc func(path string, info fs.FileInfo) bool

// AllowAll is a predefined FilterFunc that allows all files.
// It can be used as a default filter when no filtering is needed.
var AllowAll FilterFunc = func(path string, info fs.FileInfo) bool {
	return true
}

// CombineFilters combines multiple FilterFunc functions into a single FilterFunc.
// The combined filter returns true only if all individual filters return true for a given file.
// This allows for flexible and reusable filtering logic by composing simple filters.
func CombineFilters(filters ...FilterFunc) FilterFunc {
	return func(path string, info fs.FileInfo) bool {
		for _, filter := range filters {
			if !filter(path, info) {
				return false
			}
		}
		return true
	}
}
