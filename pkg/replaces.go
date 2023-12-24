package pkg

import (
	"github.com/ReallyLiri/gow/pkg/gofiles"
	"github.com/go-faster/errors"
	"golang.org/x/mod/modfile"
	"log"
	"path/filepath"
	"strings"
)

func WorkspaceSyncReplaces() error {
	workFilePath, workFile, err := gofiles.GetWorkFile()
	if err != nil {
		return err
	}
	workRoot := filepath.Dir(workFilePath)

	replaces := workFile.Replace
	modules := workFile.Use
	log.Printf("found %d replace directives and %d modules", len(replaces), len(modules))
	if len(replaces) == 0 || len(modules) == 0 {
		return nil
	}

	replacesByPath := groupReplaces(replaces)
	modPathToModFile := make(map[string]*modfile.File, len(modules))
	modPathToRelPath := make(map[string]string, len(modules))

	for _, mod := range modules {
		modFilePath := filepath.Join(workRoot, mod.Path)
		modFile, err := gofiles.GetModFile(modFilePath)
		if err != nil {
			return err
		}
		modPathToModFile[mod.Path] = modFile
		err = syncModReplaces(mod.Path, modFile, replacesByPath)
		if err != nil {
			return err
		}

		modPath := modFile.Module.Mod.Path
		modPathToRelPath[modPath] = mod.Path
	}

	for _, mod := range modules {
		modFile := modPathToModFile[mod.Path]
		err = syncModDeps(mod.Path, modFile, modPathToRelPath)
		if err != nil {
			return err
		}

		err = gofiles.WriteModFile(filepath.Join(workRoot, mod.Path, "go.mod"), modFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func groupReplaces(replaces []*modfile.Replace) map[string]*modfile.Replace {
	replaceByPkg := make(map[string]*modfile.Replace, len(replaces))
	for _, replace := range replaces {
		replaceByPkg[replace.Old.Path] = replace
	}
	return replaceByPkg
}

func syncModReplaces(modRelPath string, modFile *modfile.File, workReplaces map[string]*modfile.Replace) error {
	for _, require := range modFile.Require {
		if workReplace, ok := workReplaces[require.Mod.Path]; ok {
			log.Printf("%s: adding replace for %s", modFile.Module.Mod.Path, require.Mod.Path)
			newPath := adjustNewPath(modRelPath, workReplace)
			err := modFile.AddReplace(require.Mod.Path, "", newPath, workReplace.New.Version)
			if err != nil {
				return errors.Wrapf(err, "%s: failed to add replace for %s", modFile.Module.Mod.Path, require.Mod.Path)
			}
		}
	}
	return nil
}

func adjustNewPath(modRelPath string, workReplace *modfile.Replace) string {
	newPath := workReplace.New.Path
	if strings.Contains(newPath, "../") {
		additionalUp := strings.Count(modRelPath, "/") + 1
		newPath = strings.Replace(newPath, "../", strings.Repeat("../", additionalUp+1), 1)
	}
	return newPath
}

func syncModDeps(modRelPath string, modFile *modfile.File, modPathToRelPath map[string]string) error {
	for _, require := range modFile.Require {
		if otherModRelPath, ok := modPathToRelPath[require.Mod.Path]; ok {
			relToEachOther, err := filepath.Rel(modRelPath, otherModRelPath)
			if err != nil {
				return errors.Wrapf(err, "%s: failed to get relative path from %s to %s", modFile.Module.Mod.Path, modRelPath, otherModRelPath)
			}
			relToEachOther = "./" + relToEachOther
			log.Printf("%s: adding replace for dep %s to '%s'", modFile.Module.Mod.Path, require.Mod.Path, relToEachOther)
			err = modFile.AddReplace(require.Mod.Path, "", relToEachOther, "")
			if err != nil {
				return errors.Wrapf(err, "%s: failed to add replace for %s", modFile.Module.Mod.Path, require.Mod.Path)
			}
		}
	}
	return nil
}
