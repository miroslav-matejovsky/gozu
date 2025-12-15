package gozu

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"path"
)

// ZipMapToBytes creates an in-memory zip archive from the provided map of file paths to their byte contents.
// File paths in the map should use forward slashes (e.g., "dir/file.txt").
// Returns the zip archive as a byte slice or an error if the operation fails.
func ZipMapToBytes(fileMap FileMap) ([]byte, error) {
	if fileMap == nil {
		return nil, fmt.Errorf("files map cannot be nil")
	}

	if err := fileMap.Validate(); err != nil {
		return nil, fmt.Errorf("invalid file map: %w", err)
	}

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	for filePath, content := range fileMap {
		// Clean the path and ensure forward slashes
		cleanPath := path.Clean(filePath)
		if cleanPath == "." || cleanPath == "" {
			return nil, fmt.Errorf("invalid file path: %q", filePath)
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
