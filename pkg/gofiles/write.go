package gofiles

import (
	"golang.org/x/mod/modfile"
	"os"
)

type goFile interface {
	SortBlocks()
	Cleanup()
}

func writeFile(path string, f goFile, getSyntax func() *modfile.FileSyntax) error {
	f.SortBlocks()
	f.Cleanup()
	out := modfile.Format(getSyntax())
	return os.WriteFile(path, out, 0666)
}
