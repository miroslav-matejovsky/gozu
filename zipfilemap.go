package gozu

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
)

// FileMap represents a mapping of file paths to their byte contents.
// Keys are forward-slash separated paths (e.g., "dir/file.txt").
// This type is useful for in-memory zip operations without filesystem access.
type FileMap map[string][]byte

// Validate ensures the FileMap only contains safe, relative paths and non-nil data.
// Paths must not be empty, reference the current directory, be absolute, or include up-level segments.
func (fm FileMap) Validate() error {
	if fm == nil {
		return fmt.Errorf("file map cannot be nil")
	}
	for filePath := range fm {
		slashPath := filepath.ToSlash(filePath)
		cleanPath := path.Clean(slashPath)
		if cleanPath == "." || cleanPath == "" {
			return fmt.Errorf("invalid file path: %q", filePath)
		}
		if isAbsoluteArchivePath(cleanPath, filePath) {
			return fmt.Errorf("absolute paths are not allowed: %q", filePath)
		}
		if cleanPath != slashPath {
			return fmt.Errorf("invalid file path (contains up-level references or redundant separators): %q", filePath)
		}
	}
	return nil
}

// isAbsoluteArchivePath detects absolute paths in both Unix and Windows styles,
// regardless of the current runtime OS.
func isAbsoluteArchivePath(cleanSlashPath, original string) bool {
	if path.IsAbs(cleanSlashPath) {
		return true
	}
	if hasWindowsDrivePrefix(cleanSlashPath) {
		return true
	}
	if strings.HasPrefix(cleanSlashPath, "//") {
		return true
	}
	if strings.HasPrefix(original, `\\`) {
		return true
	}
	return false
}

func hasWindowsDrivePrefix(p string) bool {
	if len(p) < 2 {
		return false
	}
	letter := p[0]
	if (letter >= 'A' && letter <= 'Z') || (letter >= 'a' && letter <= 'z') {
		if p[1] == ':' {
			return true
		}
	}
	return false
}

// Get returns the contents of the file at the given path and a boolean indicating
// whether the file exists in the map.
func (fm FileMap) Get(filePath string) ([]byte, bool) {
	data, ok := fm[filePath]
	return data, ok
}

// Set adds or updates a file in the map with the given path and contents.
func (fm FileMap) Set(filePath string, data []byte) {
	fm[filePath] = data
}

// Delete removes a file from the map.
func (fm FileMap) Delete(filePath string) {
	delete(fm, filePath)
}

// Has returns true if the file exists in the map.
func (fm FileMap) Has(filePath string) bool {
	_, ok := fm[filePath]
	return ok
}

// Paths returns all file paths in the map.
func (fm FileMap) Paths() []string {
	paths := make([]string, 0, len(fm))
	for p := range fm {
		paths = append(paths, p)
	}
	return paths
}
