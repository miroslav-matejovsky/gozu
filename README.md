# gzu

gzu - go zip utils

## Usage

```go
import "github.com/miroslav-matejovsky/gzu"

// Zip a directory
err := gzu.Zip("path/to/dir", "output.zip")

// Unzip a zip file
err = gzu.Unzip("input.zip", "path/to/extract")
```

## Advanced Usage

Use filters to selectively include files:

```go
// Zip only .txt files
filter := func(path string, info fs.FileInfo) bool {
    return filepath.Ext(path) == ".txt"
}
err := gzu.ZipToFile("path/to/dir", "output.zip", filter)

// Unzip from bytes
data, err := gzu.ZipToBytes("path/to/dir", gzu.AllowAll)
err = gzu.UnzipFromBytes(data, "path/to/extract", gzu.AllowAll)
```

## Filters

- `AllowAll`: Includes all files
- `CombineFilters`: Combine multiple filters
