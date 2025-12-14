package gozu

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"path"
	"strings"
)

// ZipAppend merges zippedContent into the existing zip archive provided by srcZipBytes.
// The appended files are placed under the specified path inside the archive
// (use an empty path to append at the root). If srcZipBytes is empty, the
// function creates a new archive containing only the appended content.
// Returns the updated archive as a byte slice.
func ZipAppend(srcZipBytes []byte, pathToAppend string, zippedContent []byte) ([]byte, error) {
	if len(zippedContent) == 0 {
		return nil, fmt.Errorf("zipped content is empty")
	}

	normalizedPath, err := sanitizeArchivePath(pathToAppend)
	if err != nil {
		return nil, fmt.Errorf("sanitize target path: %w", err)
	}

	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	if len(srcZipBytes) > 0 {
		transform := func(f *zip.File) (string, error) {
			return f.Name, nil
		}
		if err := appendArchiveData(zipWriter, srcZipBytes, "source zip", transform); err != nil {
			return nil, fmt.Errorf("append source archive entries: %w", err)
		}
	}

	transformNewContent := func(f *zip.File) (string, error) {
		sanitized, err := sanitizeZipEntryName(f.Name)
		if err != nil {
			return "", fmt.Errorf("sanitize entry %q: %w", f.Name, err)
		}
		if sanitized == "" {
			return "", nil
		}
		if normalizedPath == "" {
			return sanitized, nil
		}
		return path.Join(normalizedPath, sanitized), nil
	}

	if err := appendArchiveData(zipWriter, zippedContent, "zipped content", transformNewContent); err != nil {
		return nil, fmt.Errorf("append new archive entries: %w", err)
	}

	if err := zipWriter.Close(); err != nil {
		return nil, fmt.Errorf("finalize zip archive: %w", err)
	}

	return buf.Bytes(), nil
}

type zipNameTransform func(f *zip.File) (string, error)

// appendArchiveData reads zipBytes, applies transform to each entry name, and writes it into zw.
func appendArchiveData(zw *zip.Writer, zipBytes []byte, description string, transform zipNameTransform) error {
	reader, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		return fmt.Errorf("read %s: %w", description, err)
	}

	for _, file := range reader.File {
		targetName, err := transform(file)
		if err != nil {
			return err
		}
		if targetName == "" {
			continue
		}
		if err := writeZipEntry(file, zw, targetName); err != nil {
			return err
		}
	}
	return nil
}

// writeZipEntry writes a single entry to zw using the provided target name.
func writeZipEntry(src *zip.File, zw *zip.Writer, targetName string) error {
	header := src.FileHeader
	header.Name = targetName
	if src.FileInfo().IsDir() && !strings.HasSuffix(header.Name, "/") {
		header.Name += "/"
	}

	writer, err := zw.CreateHeader(&header)
	if err != nil {
		return fmt.Errorf("create entry %q in archive: %w", targetName, err)
	}
	if src.FileInfo().IsDir() {
		return nil
	}

	rc, err := src.Open()
	if err != nil {
		return fmt.Errorf("open source entry %q: %w", src.Name, err)
	}
	defer func() { _ = rc.Close() }()

	if _, err = io.Copy(writer, rc); err != nil {
		return fmt.Errorf("copy entry %q: %w", src.Name, err)
	}
	return nil
}

// sanitizeArchivePath sanitizes the user-provided base path.
func sanitizeArchivePath(p string) (string, error) {
	if strings.TrimSpace(p) == "" {
		return "", nil
	}
	return sanitizeZipEntryName(p)
}

// sanitizeZipEntryName cleans archive entry names and prevents path traversal.
func sanitizeZipEntryName(name string) (string, error) {
	replaced := strings.ReplaceAll(name, "\\", "/")
	cleaned := path.Clean(replaced)
	cleaned = strings.TrimPrefix(cleaned, "/")
	if cleaned == "." {
		return "", nil
	}
	if cleaned == "" {
		return "", nil
	}
	if cleaned == ".." || strings.HasPrefix(cleaned, "../") || strings.Contains(cleaned, "/../") {
		return "", fmt.Errorf("invalid archive entry path %q: path traversal detected", name)
	}
	return cleaned, nil
}
