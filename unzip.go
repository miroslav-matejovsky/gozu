package gozu

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Unzip extracts all files from the zip archive at srcFile to the directory at dstPath.
// Creates the destination directory if it doesn't exist.
// Returns an error if the archive cannot be read or files cannot be extracted.
func Unzip(srcFile, dstPath string) error {
	return UnzipFromFile(srcFile, dstPath, AllowAll)
}

// UnzipFromFile extracts files from the zip archive at srcFile to the directory at dstPath.
// The filter function is called for each file to determine if it should be extracted.
// Use AllowAll to extract all files, or provide a custom FilterFunc for selective extraction.
// Returns an error if the archive cannot be read or files cannot be extracted.
func UnzipFromFile(srcFile, dstPath string, filter FilterFunc) error {
	if filter == nil {
		return fmt.Errorf("filter function cannot be nil")
	}

	r, err := zip.OpenReader(srcFile)
	if err != nil {
		return fmt.Errorf("open zip archive %q: %w", srcFile, err)
	}
	defer func() { _ = r.Close() }()

	return unzipFromFiles(dstPath, filter, r.File)
}

// UnzipFromBytes extracts files from an in-memory zip archive to the directory at dstPath.
// The filter function is called for each file to determine if it should be extracted.
// Use AllowAll to extract all files, or provide a custom FilterFunc for selective extraction.
// Returns an error if the archive is invalid or files cannot be extracted.
func UnzipFromBytes(data []byte, dstPath string, filter FilterFunc) error {
	if filter == nil {
		return fmt.Errorf("filter function cannot be nil")
	}
	if len(data) == 0 {
		return fmt.Errorf("zip data is empty")
	}

	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return fmt.Errorf("read zip archive from bytes: %w", err)
	}

	return unzipFromFiles(dstPath, filter, r.File)
}

// unzipFromFiles extracts the provided zip files to the directory at dstPath.
// Includes protection against zip slip attacks by validating extracted paths.
func unzipFromFiles(dstPath string, filter FilterFunc, files []*zip.File) error {
	// Clean and resolve destination path for zip slip protection
	dstPath, err := filepath.Abs(dstPath)
	if err != nil {
		return fmt.Errorf("resolve destination path: %w", err)
	}

	for _, f := range files {
		if !filter(f.Name, f.FileInfo()) {
			continue
		}

		// Construct and validate the target path (zip slip protection)
		targetPath := filepath.Join(dstPath, f.Name)
		if !strings.HasPrefix(filepath.Clean(targetPath), dstPath) {
			return fmt.Errorf("illegal file path %q: path traversal attempt detected", f.Name)
		}

		if f.FileInfo().IsDir() {
			err := os.MkdirAll(targetPath, f.Mode())
			if err != nil {
				return fmt.Errorf("create directory %q: %w", targetPath, err)
			}
			continue
		}

		err := os.MkdirAll(filepath.Dir(targetPath), 0755)
		if err != nil {
			return fmt.Errorf("create parent directory for %q: %w", targetPath, err)
		}

		err = extractFile(f, targetPath)
		if err != nil {
			return err
		}
	}
	return nil
}

// extractFile extracts a single file from the archive to the target path.
func extractFile(f *zip.File, targetPath string) error {
	out, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("create file %q: %w", targetPath, err)
	}
	defer func() { _ = out.Close() }()

	rc, err := f.Open()
	if err != nil {
		return fmt.Errorf("open archived file %q: %w", f.Name, err)
	}
	defer func() { _ = rc.Close() }()

	_, err = io.Copy(out, rc)
	if err != nil {
		return fmt.Errorf("extract file %q: %w", f.Name, err)
	}
	return nil
}
