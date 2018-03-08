package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"

	"gotask/build"
)

var compileFlag = cli.BoolFlag{Name: "compile, c", Usage: "compile the task binary to pkg.task but do not run it"}

func compileTasks(isDebug bool) (err error) {
	sourceDir, err := os.Getwd()
	if err != nil {
		return
	}

	fileName := fmt.Sprintf("%s.task", filepath.Base(sourceDir))
	outfile := filepath.Join(sourceDir, fileName)

	err = build.Compile(sourceDir, outfile, isDebug)
	return
}
