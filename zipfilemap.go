package gozu

import (
	"fmt"
	"path"
)

// FileMap represents a mapping of file paths to their byte contents.
// Keys are forward-slash separated paths (e.g., "dir/file.txt").
// This type is useful for in-memory zip operations without filesystem access.
type FileMap map[string][]byte

// ValidateFileMap checks if the provided FileMap is valid.
// It ensures that no file paths are empty or represent the current directory.
// Also only relative paths without up-level references are allowed.
func ValidateFileMap(fm FileMap) error {
	if fm == nil {
		return fmt.Errorf("file map cannot be nil")
	}
	for filePath := range fm {
		cleanPath := path.Clean(filePath)
		if cleanPath == "." || cleanPath == "" {
			return fmt.Errorf("invalid file path: %q", filePath)
		}
		if path.IsAbs(cleanPath) {
			return fmt.Errorf("absolute paths are not allowed: %q", filePath)
		}
		if cleanPath != filePath {
			return fmt.Errorf("invalid file path (contains up-level references or redundant separators): %q", filePath)
		}
	}
	return nil
}
