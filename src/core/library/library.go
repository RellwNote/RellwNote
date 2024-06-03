package library

import (
	"os"
	"path"
	"rellwnote/core/config"
)

func FileExists(filePath string) bool {
	stat, err := os.Stat(path.Join(config.LibraryPath, filePath))
	if err != nil {
		return false
	}
	return stat.IsDir() == false
}

func ReadFile(filePath string) ([]byte, error) {
	return os.ReadFile(path.Join(config.LibraryPath, filePath))
}

func GetIconFileName() (name string, has bool) {
	iconNames := []string{
		"favicon.svg",
		"favicon.ico",
		"favicon.png",
		"favicon.jpg",
		"favicon.jpeg",
	}
	for _, name := range iconNames {
		if FileExists(name) {
			return name, true
		}
	}
	return "", false
}
