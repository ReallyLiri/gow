package pkg

import (
	"fmt"
	"github.com/ReallyLiri/gow/pkg/gofiles"
	"github.com/go-faster/errors"
	"go.uber.org/multierr"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func RunForEachMod(cmdArgs []string) error {
	if len(cmdArgs) == 0 {
		return fmt.Errorf("command cannot be empty")
	}

	workFilePath, workFile, err := gofiles.GetWorkFile()
	if err != nil {
		return err
	}
	workRoot := filepath.Dir(workFilePath)

	modules := workFile.Use
	log.Printf("found %d modules", len(modules))

	for _, mod := range modules {
		log.Printf("running command '%s' for module '%s'", cmdArgs, mod.Path)
		modRoot := filepath.Join(workRoot, mod.Path)
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = modRoot

		err = cmd.Run()
		log.Printf("command '%s' for module '%s' finished with code %d", cmdArgs, mod.Path, cmd.ProcessState.ExitCode())
		if err != nil {
			err = multierr.Append(err, errors.Wrapf(err, "failed to run command '%s' for module '%s'", cmdArgs, mod.Path))
		}
	}

	return err
}
