package gozu

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// Zip compresses the directory at srcPath into a zip archive at dstPath.
// All files and subdirectories are included recursively.
// Returns an error if the source directory cannot be read or the destination
// file cannot be created.
func Zip(srcPath, dstPath string) error {
	return ZipToFile(srcPath, dstPath, AllowAll)
}

// ZipToFile compresses the directory at srcPath into a zip archive at dstPath.
// The filter function is called for each file to determine if it should be included.
// Use AllowAll to include all files, or provide a custom FilterFunc for selective archiving.
// Returns an error if the source directory cannot be read or the destination file cannot be created.
func ZipToFile(srcPath, dstPath string, filter FilterFunc) error {
	if filter == nil {
		return fmt.Errorf("filter function cannot be nil")
	}

	file, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("create zip file %q: %w", dstPath, err)
	}
	defer func() { _ = file.Close() }()

	zipWriter := zip.NewWriter(file)
	defer func() { _ = zipWriter.Close() }()

	return zipToWriter(srcPath, filter, zipWriter)
}

// ZipToBytes compresses the directory at srcPath into an in-memory zip archive.
// The filter function is called for each file to determine if it should be included.
// Use AllowAll to include all files, or provide a custom FilterFunc for selective archiving.
// Returns the zip archive as a byte slice, or an error if compression fails.
func ZipToBytes(srcPath string, filter FilterFunc) ([]byte, error) {
	if filter == nil {
		return nil, fmt.Errorf("filter function cannot be nil")
	}

	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)
	err := zipToWriter(srcPath, filter, zipWriter)
	if err != nil {
		return nil, err
	}
	err = zipWriter.Close()
	if err != nil {
		return nil, fmt.Errorf("finalize zip archive: %w", err)
	}
	return buf.Bytes(), nil
}

// zipToWriter compresses the directory at srcPath into the provided zip.Writer.
// Uses forward slashes in archive paths for cross-platform compatibility.
func zipToWriter(srcPath string, filter FilterFunc, zw *zip.Writer) error {
	err := filepath.WalkDir(srcPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("access path %q: %w", path, err)
		}
		relPath, err := filepath.Rel(srcPath, path)
		if err != nil {
			return fmt.Errorf("compute relative path for %q: %w", path, err)
		}
		if relPath == "." {
			return nil
		}
		// Use forward slashes for zip compatibility across platforms
		relPath = filepath.ToSlash(relPath)

		info, err := d.Info()
		if err != nil {
			return fmt.Errorf("get file info for %q: %w", path, err)
		}
		if !filter(relPath, info) {
			return nil
		}
		if d.IsDir() {
			_, err = zw.Create(relPath + "/")
			if err != nil {
				return fmt.Errorf("create directory entry %q in archive: %w", relPath, err)
			}
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("open source file %q: %w", path, err)
		}
		w, err := zw.Create(relPath)
		if err != nil {
			_ = f.Close()
			return fmt.Errorf("create file entry %q in archive: %w", relPath, err)
		}
		_, err = io.Copy(w, f)
		_ = f.Close()
		if err != nil {
			return fmt.Errorf("write file %q to archive: %w", relPath, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("walk directory %q: %w", srcPath, err)
	}
	return nil
}
