package gzu

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Unzip unzips the zip file at srcFile to the directory at dstPath, extracting all files.
func Unzip(srcFile, dstPath string) error {
	return UnzipFromFile(srcFile, dstPath, AllowAll)
}

// UnzipFromFile unzips the zip file at srcFile to the directory at dstPath, filtering files with the provided filter.
func UnzipFromFile(srcFile, dstPath string, filter FilterFunc) error {
	r, err := zip.OpenReader(srcFile)
	if err != nil {
		return fmt.Errorf("open zip reader: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		if !filter(f.Name, f.FileInfo()) {
			continue
		}
		path := filepath.Join(dstPath, f.Name)
		if f.FileInfo().IsDir() {
			err = os.MkdirAll(path, f.Mode())
			if err != nil {
				return fmt.Errorf("create directory %s: %w", path, err)
			}
			continue
		}
		err = os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			return fmt.Errorf("create directory %s: %w", filepath.Dir(path), err)
		}
		out, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("create file %s: %w", path, err)
		}
		defer out.Close()
		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("open zip file %s: %w", f.Name, err)
		}
		defer rc.Close()
		_, err = io.Copy(out, rc)
		if err != nil {
			return fmt.Errorf("copy file %s: %w", f.Name, err)
		}
	}
	return nil
}

// UnzipFromBytes unzips the zip data in the byte slice to the directory at dstPath, filtering files with the provided filter.
func UnzipFromBytes(data []byte, dstPath string, filter FilterFunc) error {
	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return fmt.Errorf("create zip reader: %w", err)
	}

	for _, f := range r.File {
		if !filter(f.Name, f.FileInfo()) {
			continue
		}
		path := filepath.Join(dstPath, f.Name)
		if f.FileInfo().IsDir() {
			err = os.MkdirAll(path, f.Mode())
			if err != nil {
				return fmt.Errorf("create directory %s: %w", path, err)
			}
			continue
		}
		err = os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			return fmt.Errorf("create directory %s: %w", filepath.Dir(path), err)
		}
		out, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("create file %s: %w", path, err)
		}
		defer out.Close()
		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("open zip file %s: %w", f.Name, err)
		}
		defer rc.Close()
		_, err = io.Copy(out, rc)
		if err != nil {
			return fmt.Errorf("copy file %s: %w", f.Name, err)
		}
	}
	return nil
}
