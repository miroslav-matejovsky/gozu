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

// pathOnlyFilter represents functions that filter based only on the file path.
// These functions take a string (the relative path) and return a bool indicating
// whether to include the file.
type pathOnlyFilter interface {
	~func(string) bool
}

// infoOnlyFilter represents functions that filter based only on the file info.
// These functions take an fs.FileInfo and return a bool indicating
// whether to include the file.
type infoOnlyFilter interface {
	~func(fs.FileInfo) bool
}

// filterFuncConstraint represents the full FilterFunc signature.
// This is used to allow passing the exact FilterFunc type.
type filterFuncConstraint interface {
	~func(string, fs.FileInfo) bool
}

// filterConstructor is a union type that allows functions with different signatures
// to be passed to NewFilterFunc. It includes path-only, info-only, and full FilterFunc signatures.
type filterConstructor interface {
	pathOnlyFilter | infoOnlyFilter | filterFuncConstraint
}

// NewFilterFunc builds a FilterFunc from a path-only, info-only, or full predicate.
// This generic constructor allows creating FilterFunc instances from simpler function types,
// making it easier to define filters without always needing to handle both path and info.
//
// Parameters:
//   - fn: A function that matches one of the supported signatures.
//
// Returns:
//   - A FilterFunc that adapts the input function to the full signature.
//
// Examples:
//   - Path-only: NewFilterFunc(func(path string) bool { return strings.HasSuffix(path, ".go") })
//   - Info-only: NewFilterFunc(func(info fs.FileInfo) bool { return info.Size() > 0 })
//   - Full: NewFilterFunc(func(path string, info fs.FileInfo) bool { return true })
//
// The function uses type switching to determine the input function's signature and wraps it accordingly.
// For path-only functions, the info parameter is ignored.
// For info-only functions, the path parameter is ignored.
// For full functions, it's returned as-is.
func NewFilterFunc[Fn filterConstructor](fn Fn) FilterFunc {
	switch wrapped := any(fn).(type) {
	case FilterFunc:
		// If it's already a FilterFunc, return it directly.
		return wrapped
	case func(string, fs.FileInfo) bool:
		// If it's the full signature, cast to FilterFunc.
		return FilterFunc(wrapped)
	case func(string) bool:
		// For path-only functions, create a wrapper that ignores the info.
		return func(path string, _ fs.FileInfo) bool {
			return wrapped(path)
		}
	case func(fs.FileInfo) bool:
		// For info-only functions, create a wrapper that ignores the path.
		return func(_ string, info fs.FileInfo) bool {
			return wrapped(info)
		}
	default:
		// This should not happen due to the generic constraint, but panic for safety.
		panic("unsupported filter constructor")
	}
}
