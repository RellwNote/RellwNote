package theme

import (
	"errors"
	"os"
	"rellwnote/core/config"
	"rellwnote/core/files"
)

func Load(name string) (string, error) {
	paths := []string{
		files.LibraryPath(config.ThemeDir, name+".css"),
		files.ProgramPath(config.ThemeDir, name+".css"),
	}
	for _, path := range paths {
		if files.IsFile(path) {
			read, err := os.ReadFile(path)
			if err != nil {
				continue
			}
			return string(read), nil
		}
	}
	return "", errors.New("theme " + name + " not found")
}
