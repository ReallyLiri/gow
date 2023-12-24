package gofiles

import (
	"fmt"
	"github.com/go-faster/errors"
	"golang.org/x/mod/modfile"
	"log"
	"os"
	"path/filepath"
)

func GetModFile(modPath string) (*modfile.File, error) {
	gomodPath := filepath.Join(modPath, "go.mod")
	if !fileExists(gomodPath) {
		return nil, fmt.Errorf("cannot find main module at '%s'", gomodPath)
	}
	modData, err := os.ReadFile(gomodPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read go.mod file: %s", gomodPath)
	}
	gomod, err := modfile.Parse(gomodPath, modData, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse go.mod file: %s", gomodPath)
	}
	return gomod, nil
}

func WriteModFile(path string, modFile *modfile.File) error {
	log.Printf("writing go.mod file: %s", path)
	return writeFile(path, modFile, func() *modfile.FileSyntax {
		return modFile.Syntax
	})
}
