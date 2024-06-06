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

	// ProgramDir 指 rellwnote.exe 这个文件所在的目录
	ProgramDir = ""
)

const (
	TemplateDir              = "templates"
	ExtensionDir             = "extensions"
	BuiltinExtensionFileName = "builtin"
	SummaryFileName          = "SUMMARY.md"
	OutputDirFlag            = ".build_by_rellwnote"
)

func init() {
	var err error
	ProgramDir, err = os.Executable()
	if err != nil {
		println(err.Error())
	}
	ProgramDir = filepath.Dir(ProgramDir)
}
