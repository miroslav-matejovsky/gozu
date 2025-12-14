// Package gozu provides utilities for creating and extracting zip archives in Go.
//
// gozu offers a comprehensive set of functions for working with zip files, supporting
// both filesystem-based and in-memory operations with flexible filtering capabilities.
//
// # Features
//
//   - Create zip archives from directories
//   - Extract zip archives to directories
//   - In-memory zip operations (no filesystem required)
//   - Flexible file filtering during compression and extraction
//   - Append content to existing zip archives
//   - FileMap type for pure in-memory file manipulation
//   - Cross-platform path handling (uses forward slashes)
//   - Protection against zip slip attacks
//
// # Basic Operations
//
// The simplest way to create and extract zip archives:
//
//	// Zip a directory
//	err := gozu.Zip("path/to/dir", "output.zip")
//
//	// Unzip to a directory
//	err := gozu.Unzip("input.zip", "path/to/extract")
//
// # Filtered Operations
//
// Use FilterFunc to selectively include or exclude files:
//
//	// Zip only .go files
//	filter := func(path string, info fs.FileInfo) bool {
//	    return strings.HasSuffix(path, ".go")
//	}
//	err := gozu.ZipToFile("src", "source.zip", filter)
//
//	// Extract only files smaller than 1MB
//	sizeFilter := func(path string, info fs.FileInfo) bool {
//	    return info.Size() < 1024*1024
//	}
//	err := gozu.UnzipFromFile("archive.zip", "output", sizeFilter)
//
// Use [NewFilterFunc] to create filters from simpler predicates:
//
//	// Path-only filter
//	pathFilter := gozu.NewFilterFunc(func(path string) bool {
//	    return strings.HasPrefix(path, "docs/")
//	})
//
//	// FileInfo-only filter
//	infoFilter := gozu.NewFilterFunc(func(info fs.FileInfo) bool {
//	    return !info.IsDir()
//	})
//
// Combine multiple filters with [CombineFilters]:
//
//	combined := gozu.CombineFilters(filter1, filter2, filter3)
//
// # In-Memory Operations
//
// Work with zip archives entirely in memory:
//
//	// Compress directory to bytes
//	data, err := gozu.ZipToBytes("path/to/dir", gozu.AllowAll)
//
//	// Extract bytes to directory
//	err := gozu.UnzipFromBytes(data, "path/to/extract", gozu.AllowAll)
//
// # FileMap Operations
//
// For pure in-memory operations without filesystem access, use [FileMap]:
//
//	// Create a FileMap with file contents
//	files := gozu.FileMap{
//	    "readme.txt":     []byte("Hello, World!"),
//	    "src/main.go":    []byte("package main"),
//	    "data/config.json": []byte(`{"key": "value"}`),
//	}
//
//	// Zip the FileMap to bytes
//	zipData, err := gozu.ZipMapToBytes(files)
//
//	// Unzip bytes back to a FileMap
//	extracted, err := gozu.UnzipBytesToMap(zipData)
//
// FileMap provides helper methods for manipulation:
//
//	files.Set("new.txt", []byte("new content"))
//	content, ok := files.Get("readme.txt")
//	files.Delete("old.txt")
//	exists := files.Has("readme.txt")
//	paths := files.Paths()
//
// # Appending to Archives
//
// Add content to existing zip archives:
//
//	// Append new zip content under a specific path
//	updated, err := gozu.ZipAppend(existingZip, "subdir", newContent)
//
// # Error Handling
//
// All functions return wrapped errors with context using fmt.Errorf and %w,
// making it easy to understand the cause of failures and use errors.Is/errors.As.
//
// # Thread Safety
//
// The functions in this package are not safe for concurrent use on the same
// files or data. Callers should synchronize access when needed.
package gozu
