# gozu

gozu - go zip utils

## Usage

```go
import "github.com/miroslav-matejovsky/gozu"

// Zip a directory
err := gozu.Zip("path/to/dir", "output.zip")

// Unzip a zip file
err = gozu.Unzip("input.zip", "path/to/extract")
```

## Advanced Usage

Use filters to selectively include files:

```go
// Zip only .txt files
filter := func(path string, info fs.FileInfo) bool {
    return filepath.Ext(path) == ".txt"
}
err := gozu.ZipToFile("path/to/dir", "output.zip", filter)

// Unzip from bytes
data, err := gozu.ZipToBytes("path/to/dir", gozu.AllowAll)
err = gozu.UnzipFromBytes(data, "path/to/extract", gozu.AllowAll)
```

## Filters

- `AllowAll`: Includes all files
- `CombineFilters`: Combine multiple filters
