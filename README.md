# gozu

[![Go Reference](https://pkg.go.dev/badge/github.com/miroslav-matejovsky/gozu.svg)](https://pkg.go.dev/github.com/miroslav-matejovsky/gozu)

gozu - Go Zip Utilities

A comprehensive Go library for creating and extracting zip archives with support for filtering, in-memory operations, and archive manipulation.

## Features

- **Basic Operations**: Zip directories and unzip archives with a single function call
- **Filtered Operations**: Selectively include/exclude files during compression and extraction
- **In-Memory Operations**: Work with zip archives entirely in memory without filesystem access
- **FileMap Type**: Pure in-memory file manipulation using `map[string][]byte`
- **Archive Appending**: Merge content into existing zip archives
- **Cross-Platform**: Uses forward slashes for consistent path handling
- **Security**: Built-in protection against zip slip attacks

## Installation

```bash
go get github.com/miroslav-matejovsky/gozu
```

## Usage

### Basic Operations

```go
import "github.com/miroslav-matejovsky/gozu"

// Zip a directory
err := gozu.Zip("path/to/dir", "output.zip")

// Unzip a zip file
err = gozu.Unzip("input.zip", "path/to/extract")
```

### Filtered Operations

Use filters to selectively include files:

```go
// Zip only .txt files
filter := func(path string, info fs.FileInfo) bool {
    return filepath.Ext(path) == ".txt"
}
err := gozu.ZipToFile("path/to/dir", "output.zip", filter)

// Unzip only files smaller than 1MB
sizeFilter := func(path string, info fs.FileInfo) bool {
    return info.Size() < 1024*1024
}
err := gozu.UnzipFromFile("archive.zip", "output", sizeFilter)
```

### In-Memory Operations

Work with zip archives as byte slices:

```go
// Compress directory to bytes
data, err := gozu.ZipToBytes("path/to/dir", gozu.AllowAll)

// Extract bytes to directory
err = gozu.UnzipFromBytes(data, "path/to/extract", gozu.AllowAll)
```

### FileMap Operations

For pure in-memory operations without filesystem access:

```go
// Create a FileMap with file contents
files := gozu.FileMap{
    "readme.txt":      []byte("Hello, World!"),
    "src/main.go":     []byte("package main"),
    "data/config.json": []byte(`{"key": "value"}`),
}

// Zip the FileMap to bytes
zipData, err := gozu.ZipMapToBytes(files)

// Unzip bytes back to a FileMap
extracted, err := gozu.UnzipBytesToMap(zipData)

// FileMap helper methods
files.Set("new.txt", []byte("new content"))
content, ok := files.Get("readme.txt")
files.Delete("old.txt")
exists := files.Has("readme.txt")
paths := files.Paths()
```

### Archive Appending

Add content to existing zip archives:

```go
// Append new zip content under a specific path
updated, err := gozu.ZipAppend(existingZip, "subdir", newContent)
```

## Filters

### Predefined Filters

- `AllowAll`: Includes all files (default behavior)

### Combining Filters

```go
combined := gozu.CombineFilters(filter1, filter2, filter3)
```

### Constructing Filters

`NewFilterFunc` accepts either a full `FilterFunc`, a path-only predicate, or an info-only predicate:

```go
// Path-only filter
pathOnly := gozu.NewFilterFunc(func(path string) bool {
    return strings.HasPrefix(path, "docs/")
})

// FileInfo-only filter
infoOnly := gozu.NewFilterFunc(func(info fs.FileInfo) bool {
    return info.Size() > 0
})
```

## API Reference

### Types

| Type | Description |
|------|-------------|
| `FileMap` | `map[string][]byte` for in-memory file representation |
| `FilterFunc` | `func(path string, info fs.FileInfo) bool` for filtering |

### Functions

| Function | Description |
|----------|-------------|
| `Zip(src, dst)` | Zip a directory to a file |
| `Unzip(src, dst)` | Unzip a file to a directory |
| `ZipToFile(src, dst, filter)` | Zip with filtering |
| `ZipToBytes(src, filter)` | Zip to in-memory bytes |
| `UnzipFromFile(src, dst, filter)` | Unzip with filtering |
| `UnzipFromBytes(data, dst, filter)` | Unzip from in-memory bytes |
| `ZipMapToBytes(files)` | Zip a FileMap to bytes |
| `UnzipBytesToMap(data)` | Unzip bytes to a FileMap |
| `ZipAppend(src, path, content)` | Append content to a zip archive |

## License

See [LICENSE](LICENSE) for details.
