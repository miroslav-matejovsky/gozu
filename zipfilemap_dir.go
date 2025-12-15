package gozu

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// ExportToDir writes the FileMap contents to destDir, creating directories as needed.
// Existing files with the same paths are overwritten.
func (fm FileMap) ExportToDir(destDir string) error {
	if fm == nil {
		return fmt.Errorf("file map cannot be nil")
	}
	if destDir == "" {
		return fmt.Errorf("destination directory cannot be empty")
	}
	if len(fm) == 0 {
		return fmt.Errorf("file map is empty")
	}
	if err := fm.Validate(); err != nil {
		return fmt.Errorf("invalid file map: %w", err)
	}
	absDest, err := filepath.Abs(destDir)
	if err != nil {
		return fmt.Errorf("resolve destination directory: %w", err)
	}
	if err := os.MkdirAll(absDest, 0o755); err != nil {
		return fmt.Errorf("create destination directory %q: %w", absDest, err)
	}
	for relPath, data := range fm {
		targetPath := filepath.Join(absDest, filepath.FromSlash(relPath))
		dirPath := filepath.Dir(targetPath)
		if err := os.MkdirAll(dirPath, 0o755); err != nil {
			return fmt.Errorf("create parent directory for %q: %w", targetPath, err)
		}
		if err := os.WriteFile(targetPath, data, 0o644); err != nil {
			return fmt.Errorf("write file %q: %w", targetPath, err)
		}
	}
	return nil
}

// FileMapFromDir builds a FileMap from all files found under srcDir.
// Paths in the resulting map use forward slashes and are relative to srcDir.
func FileMapFromDir(srcDir string) (FileMap, error) {
	if srcDir == "" {
		return nil, fmt.Errorf("source directory cannot be empty")
	}
	absSrc, err := filepath.Abs(srcDir)
	if err != nil {
		return nil, fmt.Errorf("resolve source directory: %w", err)
	}
	info, err := os.Stat(absSrc)
	if err != nil {
		return nil, fmt.Errorf("stat source directory %q: %w", absSrc, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("source path %q is not a directory", absSrc)
	}
	files := make(FileMap)
	err = filepath.WalkDir(absSrc, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("access path %q: %w", path, walkErr)
		}
		if d.IsDir() {
			return nil
		}
		relPath, err := filepath.Rel(absSrc, path)
		if err != nil {
			return fmt.Errorf("compute relative path for %q: %w", path, err)
		}
		relPath = filepath.ToSlash(relPath)
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read file %q: %w", path, err)
		}
		files[relPath] = data
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
