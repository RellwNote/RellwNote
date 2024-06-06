package library

import (
	"rellwnote/core/files"
)

func GetIconFileName() (name string, has bool) {
	iconNames := []string{
		"favicon.svg",
		"favicon.ico",
		"favicon.png",
		"favicon.jpg",
		"favicon.jpeg",
	}
	for _, name := range iconNames {
		if files.IsFile(files.LibraryPath(name)) {
			return name, true
		}
	}
	return "", false
}
