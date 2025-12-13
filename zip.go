package gzu

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// Zip zips the directory at srcPath to a zip file at dstPath, including all files.
func Zip(srcPath, dstPath string) error {
	return ZipToFile(srcPath, dstPath, AllowAll)
}

// ZipToFile zips the directory at srcPath to a zip file at dstPath, filtering files with the provided filter.
func ZipToFile(srcPath, dstPath string, filter FilterFunc) error {
	file, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("create zip file: %w", err)
	}
	defer func() { _ = file.Close() }()

	zipWriter := zip.NewWriter(file)
	defer func() { _ = zipWriter.Close() }()

	return zipToWriter(srcPath, filter, zipWriter)
}

// ZipToBytes zips the directory at srcPath to a byte slice, filtering files with the provided filter.
func ZipToBytes(srcPath string, filter FilterFunc) ([]byte, error) {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)
	err := zipToWriter(srcPath, filter, zipWriter)
	if err != nil {
		return nil, err
	}
	err = zipWriter.Close()
	if err != nil {
		return nil, fmt.Errorf("close zip writer: %w", err)
	}
	return buf.Bytes(), nil
}

// zipToWriter zips the directory at srcPath to the provided zip.Writer, filtering files with the provided filter.
func zipToWriter(srcPath string, filter FilterFunc, zw *zip.Writer) error {
	err := filepath.WalkDir(srcPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(srcPath, path)
		if err != nil {
			return err
		}
		if relPath == "." {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		if !filter(relPath, info) {
			return nil
		}
		if d.IsDir() {
			_, err = zw.Create(relPath + "/")
			return err
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		w, err := zw.Create(relPath)
		if err != nil {
			_ = f.Close()
			return err
		}
		_, err = io.Copy(w, f)
		_ = f.Close()
		return err
	})
	if err != nil {
		return fmt.Errorf("walk directory: %w", err)
	}
	return nil
}
