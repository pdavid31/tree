package internal

import (
	"os"
	"strings"
)

func applyFilters(f []os.FileInfo, config *TreeConfig) []os.FileInfo {
	c := make([]os.FileInfo, len(f))
	copy(c, f)

	if !config.AllFiles {
		c = removeHidden(c)
	}

	if config.DirectoriesOnly {
		c = removeFiles(c)
	}

	return c
}

// filter out all hidden files
func removeHidden(f []os.FileInfo) []os.FileInfo {
	n := 0
	for _, x := range f {
		if !strings.HasPrefix(x.Name(), ".") {
			f[n] = x
			n++
		}
	}
	return f[:n]
}

// filter out all files (non directories)
func removeFiles(f []os.FileInfo) []os.FileInfo {
	n := 0
	for _, x := range f {
		if x.IsDir() {
			f[n] = x
			n++
		}
	}
	return f[:n]
}
