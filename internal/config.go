package internal

import (
	"github.com/gobwas/glob"
)

type TreeConfig struct {
	AllFiles           bool
	DirectoriesOnly    bool
	DisableIndentation bool
	FullPaths          bool
	Pattern            glob.Glob
}
