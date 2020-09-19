package internal

import (
	"os"
	"strings"

	"github.com/gobwas/glob"
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

	if config.Pattern != "" {
		c = applyPattern(c, config.Pattern)
	}

	return c
}

// filter out all hidden files
func removeHidden(f []os.FileInfo) []os.FileInfo {
	return inplaceFilter(f, func(x os.FileInfo) bool {
		return !strings.HasPrefix(x.Name(), ".")
	})
}

// filter out all files (non directories)
func removeFiles(f []os.FileInfo) []os.FileInfo {
	return inplaceFilter(f, func(x os.FileInfo) bool {
		return x.IsDir()
	})
}

func applyPattern(f []os.FileInfo, pattern string) []os.FileInfo {
	g, err := glob.Compile(pattern)
	if err != nil {
		return nil
	}

	return inplaceFilter(f, func(x os.FileInfo) bool {
		return x.IsDir() || g.Match(x.Name())
	})
}

func inplaceFilter(f []os.FileInfo, fn func(os.FileInfo) bool) []os.FileInfo {
	n := 0
	for _, x := range f {
		if fn(x) {
			f[n] = x
			n++
		}
	}
	return f[:n]
}
