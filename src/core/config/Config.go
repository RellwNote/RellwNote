package config

import (
	"os"
	"path/filepath"
)

var (
	LibraryPath      = "library"
	LibraryName      = ""
	ServerPort       = 8333
	ServerHost       = "localhost"
	ServerDebugDelay = 0.0
	BuildOutput      = "build"

	BaseDir = ""
)

const (
	TemplateDir     = "templates"
	SummaryFileName = "SUMMARY.md"
)

func init() {
	var err error
	BaseDir, err = os.Executable()
	if err != nil {
		println(err.Error())
	}
	BaseDir = filepath.Dir(BaseDir)
}
