package gozu

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileMapFromDir_ReadsFiles(t *testing.T) {
	tempDir := t.TempDir()
	nested := filepath.Join(tempDir, "nested", "deep")
	err := os.MkdirAll(nested, 0o755)
	require.NoError(t, err)

	files := map[string][]byte{
		"root.txt":            []byte("root"),
		"nested/file.txt":     []byte("nested"),
		"nested/deep/bin.bin": {0x00, 0x01},
	}

	for rel, data := range files {
		absPath := filepath.Join(tempDir, filepath.FromSlash(rel))
		err = os.MkdirAll(filepath.Dir(absPath), 0o755)
		require.NoError(t, err)
		err = os.WriteFile(absPath, data, 0o644)
		require.NoError(t, err)
	}

	fm, err := FileMapFromDir(tempDir)
	require.NoError(t, err)
	require.Len(t, fm, len(files))
	for rel, data := range files {
		require.Equal(t, data, fm[rel])
	}
}

func TestFileMapFromDir_EmptyDir(t *testing.T) {
	fm, err := FileMapFromDir(t.TempDir())
	require.NoError(t, err)
	require.Empty(t, fm)
}

func TestFileMapFromDir_InvalidInput(t *testing.T) {
	_, err := FileMapFromDir("nonexistent")
	require.Error(t, err)
	require.Contains(t, err.Error(), "stat source directory")

	file := filepath.Join(t.TempDir(), "not_a_dir.txt")
	err = os.WriteFile(file, []byte("content"), 0o644)
	require.NoError(t, err)

	_, err = FileMapFromDir(file)
	require.Error(t, err)
	require.Contains(t, err.Error(), "is not a directory")
}

func TestFileMap_ExportToDir(t *testing.T) {
	fm := FileMap{
		"root.txt":    []byte("root"),
		"dir/sub.txt": []byte("sub"),
		"a/b/c.bin":   []byte{0x01, 0x02},
	}

	dest := filepath.Join(t.TempDir(), "output")
	err := fm.ExportToDir(dest)
	require.NoError(t, err)

	for rel, expect := range fm {
		abs := filepath.Join(dest, filepath.FromSlash(rel))
		data, err := os.ReadFile(abs)
		require.NoError(t, err)
		require.Equal(t, expect, data)
	}

	var exported []string
	err = filepath.Walk(dest, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(dest, path)
		if err != nil {
			return err
		}
		exported = append(exported, filepath.ToSlash(rel))
		return nil
	})
	require.NoError(t, err)
	sort.Strings(exported)
	var expected []string
	for rel := range fm {
		expected = append(expected, rel)
	}
	sort.Strings(expected)
	require.Equal(t, expected, exported)
}

func TestFileMap_ExportToDir_InvalidInput(t *testing.T) {
	var nilMap FileMap
	err := nilMap.ExportToDir(t.TempDir())
	require.Error(t, err)
	require.Contains(t, err.Error(), "file map cannot be nil")

	empty := FileMap{}
	err = empty.ExportToDir(t.TempDir())
	require.Error(t, err)
	require.Contains(t, err.Error(), "file map is empty")

	fm := FileMap{"../escape.txt": []byte("nope")}
	err = fm.ExportToDir(t.TempDir())
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid file map")

	fm = FileMap{"//share/file.txt": []byte("nope")}
	err = fm.ExportToDir(t.TempDir())
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid file map")

	fm = FileMap{"file.txt": []byte("ok")}
	err = fm.ExportToDir("")
	require.Error(t, err)
	require.Contains(t, err.Error(), "destination directory cannot be empty")
}
