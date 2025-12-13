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
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	err = filepath.WalkDir(srcPath, func(path string, d fs.DirEntry, err error) error {
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
			_, err = zipWriter.Create(relPath + "/")
			return err
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		w, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}
		_, err = io.Copy(w, f)
		return err
	})
	if err != nil {
		return fmt.Errorf("walk directory: %w", err)
	}
	return nil
}

// ZipToBytes zips the directory at srcPath to a byte slice, filtering files with the provided filter.
func ZipToBytes(srcPath string, filter FilterFunc) ([]byte, error) {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)
	defer zipWriter.Close()

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
			_, err = zipWriter.Create(relPath + "/")
			return err
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		w, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}
		_, err = io.Copy(w, f)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("walk directory: %w", err)
	}
	zipWriter.Close()
	return buf.Bytes(), nil
}
