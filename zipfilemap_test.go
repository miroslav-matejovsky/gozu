package gozu

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateFileMap(t *testing.T) {
	var nilMap FileMap

	tests := []struct {
		name    string
		fm      FileMap
		wantErr string
	}{
		{name: "valid", fm: FileMap{"file.txt": []byte("content"), "dir/sub.go": []byte("package main")}},
		{name: "empty map", fm: FileMap{}},
		{name: "nil map", fm: nilMap, wantErr: "file map cannot be nil"},
		{name: "empty path", fm: FileMap{"": []byte("x")}, wantErr: "invalid file path"},
		{name: "dot path", fm: FileMap{".": []byte("x")}, wantErr: "invalid file path"},
		{name: "absolute path", fm: FileMap{"/abs/path.txt": []byte("x")}, wantErr: "absolute paths are not allowed"},
		{name: "windows absolute", fm: FileMap{"C:/abs/path.txt": []byte("x")}, wantErr: "absolute paths are not allowed"},
		{name: "up-level reference", fm: FileMap{"dir/../evil.txt": []byte("x")}, wantErr: "invalid file path (contains up-level references or redundant separators)"},
		{name: "redundant separators", fm: FileMap{"dir//file.txt": []byte("x")}, wantErr: "invalid file path (contains up-level references or redundant separators)"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := tc.fm.Validate()
			if tc.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.wantErr)
		})
	}
}
