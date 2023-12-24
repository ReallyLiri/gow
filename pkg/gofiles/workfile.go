package gofiles

import (
	"fmt"
	"github.com/go-faster/errors"
	"golang.org/x/mod/modfile"
	"log"
	"os"
	"path/filepath"
)

func WriteWorkFile(path string, workFile *modfile.WorkFile) error {
	log.Printf("writing go.work file: %s", path)
	return writeFile(path, workFile, func() *modfile.FileSyntax {
		return workFile.Syntax
	})
}

func GetWorkFile() (string, *modfile.WorkFile, error) {
	workFilePath, err := findGoWork()
	if err != nil {
		return "", nil, err
	}
	if workFilePath == "" {
		return "", nil, errors.New("go.work file not found")
	}
	log.Printf("found go.work file: %s", workFilePath)
	workData, err := os.ReadFile(workFilePath)
	if err != nil {
		return "", nil, errors.Wrapf(err, "failed to read go.work file: %s", workFilePath)
	}

	gowork, err := modfile.ParseWork(workFilePath, workData, nil)
	if err != nil {
		return "", nil, errors.Wrapf(err, "failed to parse go.work file: %s", workFilePath)
	}
	return workFilePath, gowork, nil
}

func findGoWork() (string, error) {
	switch gowork := os.Getenv("GOWORK"); gowork {
	case "off":
		return "", nil
	case "", "auto":
		wd, err := os.Getwd()
		if err != nil {
			return "", errors.Wrap(err, "cannot determine working directory")
		}
		return findEnclosingFile(wd, "go.work"), nil
	default:
		if !filepath.IsAbs(gowork) {
			return "", fmt.Errorf("go: invalid GOWORK: not an absolute path: %s", gowork)
		}
		if !fileExists(gowork) {
			return "", fmt.Errorf("go: invalid GOWORK: does not exist: %s", gowork)
		}
		return gowork, nil
	}
}

func findEnclosingFile(dirPath string, fileName string) string {
	dirPath = filepath.Clean(dirPath)
	for {
		f := filepath.Join(dirPath, fileName)
		if fileExists(f) {
			return f
		}
		d := filepath.Dir(dirPath)
		if d == dirPath {
			break
		}
		if d == os.Getenv("GOROOT") {
			return ""
		}
		dirPath = d
	}
	return ""
}

func fileExists(fpath string) bool {
	if fi, err := os.Stat(fpath); err == nil && !fi.IsDir() {
		return true
	}
	return false
}
