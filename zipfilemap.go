package gozu

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"path"
)

// FileMap represents a mapping of file paths to their byte contents.
// Keys are forward-slash separated paths (e.g., "dir/file.txt").
// This type is useful for in-memory zip operations without filesystem access.
type FileMap map[string][]byte

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

// ZipMapToBytes creates an in-memory zip archive from the provided map of file paths to their byte contents.
// File paths in the map should use forward slashes (e.g., "dir/file.txt").
// Returns the zip archive as a byte slice or an error if the operation fails.
func ZipMapToBytes(files FileMap) ([]byte, error) {
	if files == nil {
		return nil, fmt.Errorf("files map cannot be nil")
	}

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	for filePath, content := range files {
		// Clean the path and ensure forward slashes
		cleanPath := path.Clean(filePath)
		if cleanPath == "." || cleanPath == "" {
			continue
		}

		w, err := zw.Create(cleanPath)
		if err != nil {
			return nil, fmt.Errorf("create entry %q in archive: %w", cleanPath, err)
		}

		_, err = w.Write(content)
		if err != nil {
			return nil, fmt.Errorf("write content for %q: %w", cleanPath, err)
		}
	}

	err := zw.Close()
	if err != nil {
		return nil, fmt.Errorf("finalize zip archive: %w", err)
	}

	return buf.Bytes(), nil
}

// UnzipBytesToMap extracts files from an in-memory zip archive represented by the provided byte slice.
// Returns a FileMap containing file paths mapped to their byte contents.
// Directories are not included in the returned map.
// Returns an error if the archive is invalid or files cannot be read.
func UnzipBytesToMap(data []byte) (FileMap, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("zip data is empty")
	}

	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("read zip archive from bytes: %w", err)
	}

	files := make(FileMap, len(r.File))
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("open file %q in archive: %w", f.Name, err)
		}

		content, err := io.ReadAll(rc)
		if err != nil {
			_ = rc.Close()
			return nil, fmt.Errorf("read file %q from archive: %w", f.Name, err)
		}
		_ = rc.Close()

		files[f.Name] = content
	}

	return files, nil
}
