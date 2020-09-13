package internal

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_applyFilters(t *testing.T) {
	tests := []struct {
		name     string
		arg      []os.FileInfo
		config   *TreeConfig
		expected []os.FileInfo
	}{
		{
			name: "empty config",
			arg: []os.FileInfo{
				fakeFile{
					name:  ".idea",
					isDir: true,
				},
				fakeFile{
					name:  ".gitignore",
					isDir: false,
				},
				fakeFile{
					name:  "cmd",
					isDir: true,
				},
				fakeFile{
					name:  "go.mod",
					isDir: false,
				},
			},
			config: &TreeConfig{},
			expected: []os.FileInfo{
				fakeFile{
					name:  "cmd",
					isDir: true,
				},
				fakeFile{
					name:  "go.mod",
					isDir: false,
				},
			},
		},
		{
			name: "directories only",
			arg: []os.FileInfo{
				fakeFile{
					name:  ".idea",
					isDir: true,
				},
				fakeFile{
					name:  ".gitignore",
					isDir: false,
				},
				fakeFile{
					name:  "cmd",
					isDir: true,
				},
				fakeFile{
					name:  "go.mod",
					isDir: false,
				},
			},
			config: &TreeConfig{DirectoriesOnly: true},
			expected: []os.FileInfo{
				fakeFile{
					name:  "cmd",
					isDir: true,
				},
			},
		},
		{
			name: "all files",
			arg: []os.FileInfo{
				fakeFile{
					name:  ".idea",
					isDir: true,
				},
				fakeFile{
					name:  ".gitignore",
					isDir: false,
				},
				fakeFile{
					name:  "cmd",
					isDir: true,
				},
				fakeFile{
					name:  "go.mod",
					isDir: false,
				},
			},
			config: &TreeConfig{AllFiles: true},
			expected: []os.FileInfo{
				fakeFile{
					name:  ".idea",
					isDir: true,
				},
				fakeFile{
					name:  ".gitignore",
					isDir: false,
				},
				fakeFile{
					name:  "cmd",
					isDir: true,
				},
				fakeFile{
					name:  "go.mod",
					isDir: false,
				},
			},
		},
		{
			name: "all files and only directories",
			arg: []os.FileInfo{
				fakeFile{
					name:  ".idea",
					isDir: true,
				},
				fakeFile{
					name:  ".gitignore",
					isDir: false,
				},
				fakeFile{
					name:  "cmd",
					isDir: true,
				},
				fakeFile{
					name:  "go.mod",
					isDir: false,
				},
			},
			config: &TreeConfig{AllFiles: true, DirectoriesOnly: true},
			expected: []os.FileInfo{
				fakeFile{
					name:  ".idea",
					isDir: true,
				},
				fakeFile{
					name:  "cmd",
					isDir: true,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := applyFilters(test.arg, test.config)

			assert.Equal(t, test.expected, f)
		})
	}
}
