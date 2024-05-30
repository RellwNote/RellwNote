package config

import (
	"os"
	"path"
)

var (
	LibraryPath      = "library"
	ServerPort       = 8333
	ServerHost       = "localhost"
	ServerDebugDelay = 0.0
)

const (
	TemplateDir     = "templates"
	SummaryFileName = "SUMMARY.md"
)

func LibraryFileExists(filePath string) bool {
	stat, err := os.Stat(path.Join(LibraryPath, filePath))
	if err != nil {
		return false
	}
	return stat.IsDir() == false
}

func ReadLibraryFile(filePath string) ([]byte, error) {
	return os.ReadFile(path.Join(LibraryPath, filePath))
}
