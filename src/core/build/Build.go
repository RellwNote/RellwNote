package build

import (
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"rellwnote/core/config"
	"rellwnote/core/library"
	"rellwnote/core/log"
	"rellwnote/core/template"
)

const (
	OutputDirFlag = ".build_by_rellwnote"
)

func Build() error {
	err := recreateOutputDirectory()
	if err != nil {
		return err
	}
	log.Info.Printf("[OK] clean build directory")
	err = copyDirContentTo(config.LibraryPath, config.BuildOutput)
	if err != nil {
		return err
	}
	log.Info.Printf("[OK] copy library")
	err = copyFaviconIfMission()
	if err != nil {
		return err
	}
	log.Info.Printf("[OK] favicon process")
	err = buildTemplate()
	if err != nil {
		return err
	}
	log.Info.Printf("[OK] build html")
	return nil
}

func buildTemplate() error {
	templateFuncs := map[string]func() ([]byte, error){
		"content.html": template.BuildContentPage,
		"index.html":   template.BuildIndexPage,
	}
	for fileName, builder := range templateFuncs {
		page, err := builder()
		if err != nil {
			return err
		}
		f, err := os.OpenFile(path.Join(config.BuildOutput, fileName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.Write(page)
		if err != nil {
			return err
		}
	}
	return nil
}

func copyFaviconIfMission() error {
	if _, hasIcon := library.GetIconFileName(); hasIcon {
		return nil
	}
	return copyFileTo(path.Join(config.TemplateDir, "favicon.svg"), path.Join(config.BuildOutput, "favicon.svg"))
}

func recreateOutputDirectory() error {
	err := removeOutputDirectory()
	if err != nil {
		return err
	}
	err = createOutputDirectory()
	if err != nil {
		return err
	}
	return nil
}

// 复制一个目录下的全部内容到另一个位置
func copyDirContentTo(sourceDir, targetDir string) error {
	absPath, err := filepath.Abs(sourceDir)
	if err != nil {
		return err
	}
	err = filepath.Walk(absPath, func(walkPath string, info os.FileInfo, err error) error {
		if absPath == walkPath {
			return nil
		}
		targetPath := path.Join(targetDir, walkPath[len(absPath)+1:])
		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}
		return copyFileTo(walkPath, targetPath)
	})
	return err
}

func copyFileTo(sourcePath, targetPath string) error {
	target, err := os.OpenFile(targetPath, os.O_RDWR|os.O_CREATE, 0644)
	origin, err := os.OpenFile(sourcePath, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer origin.Close()
	defer target.Close()
	_, err = io.Copy(target, origin)
	return err
}

// 创建输出目录并在里面添加标志文件
func createOutputDirectory() error {
	err := os.MkdirAll(config.BuildOutput, os.ModePerm)
	if err != nil {
		return err
	}
	create, err := os.Create(path.Join(config.BuildOutput, OutputDirFlag))
	if err != nil {
		return err
	}
	err = create.Close()
	if err != nil {
		return err
	}
	return nil
}

// 尝试移除输出目录，并检测是否是正确的非空或由 RellwNote 管理的目录
func removeOutputDirectory() error {
	stat, err := os.Stat(config.BuildOutput)
	if err != nil {
		return nil
	}
	if !stat.IsDir() {
		return errors.New("output directory isn't a folder")
	}

	empty, err := checkDirIsEmpty(config.BuildOutput)
	if err != nil {
		return err
	}
	if !empty {
		stat, err = os.Stat(path.Join(config.BuildOutput, OutputDirFlag))
		if err != nil {
			return errors.New("output directory is not empty")
		}
	}

	err = os.RemoveAll(config.BuildOutput)
	if err != nil {
		return err
	}
	return nil
}

func checkDirIsEmpty(dir string) (bool, error) {
	open, err := os.Open(dir)
	if err != nil {
		return false, err
	}
	defer open.Close()

	_, err = open.Readdirnames(1)
	if err != nil {
		return true, nil
	} else {
		return false, nil
	}
}

func isFile(filepath string) bool {
	stat, err := os.Stat(filepath)
	if err != nil {
		return false
	}
	return !stat.IsDir()
}

func isDir(filepath string) bool {
	stat, err := os.Stat(filepath)
	if err != nil {
		return false
	}
	return stat.IsDir()
}
