package files

import (
	"os"
	"path/filepath"
	"rellwnote/core/config"
)

func IsFile(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !stat.IsDir()
}

func IsDir(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

func IsEmptyDir(path string) bool {
	if !IsDir(path) {
		return false
	}
	dirs, err := os.ReadDir(path)
	if err != nil {
		return false
	}
	return len(dirs) == 0
}

func Copy(sourcePath, targetPath string) error {
	file, err := os.ReadFile(sourcePath)
	if err != nil {
		return err
	}
	err = os.WriteFile(targetPath, file, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func CopyDirContentTo(sourceDir, targetDir string) error {
	absPath, err := filepath.Abs(sourceDir)
	if err != nil {
		return err
	}
	err = filepath.Walk(absPath, func(walkPath string, info os.FileInfo, err error) error {
		if absPath == walkPath {
			return nil
		}
		targetPath := filepath.Join(targetDir, walkPath[len(absPath)+1:])
		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}
		return Copy(walkPath, targetPath)
	})
	return err
}

func LibraryPath(path ...string) string {
	return filepath.Join(append([]string{config.LibraryPath}, path...)...)
}

func ProgramPath(path ...string) string {
	return filepath.Join(append([]string{config.ProgramDir}, path...)...)
}

func OutputPath(path ...string) string {
	return filepath.Join(append([]string{config.BuildOutput}, path...)...)
}
