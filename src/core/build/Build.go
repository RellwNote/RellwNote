package build

import (
	"errors"
	"os"
	"rellwnote/core/config"
	"rellwnote/core/files"
	"rellwnote/core/library"
	"rellwnote/core/log"
	"rellwnote/core/template"
)

func Build() error {
	err := recreateOutputDirectory()
	if err != nil {
		return err
	}
	log.Info.Printf("[OK] clean build directory")
	err = files.CopyDirContentTo(config.LibraryPath, config.BuildOutput)
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
		err = os.WriteFile(files.OutputPath(fileName), page, os.ModePerm)
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
	return files.CopyDirContentTo(
		files.ProgramPath(config.TemplateDir, "favicon.svg"),
		files.OutputPath("favicon.svg"))
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

// 创建输出目录并在里面添加标志文件
func createOutputDirectory() error {
	err := os.MkdirAll(config.BuildOutput, os.ModePerm)
	if err != nil {
		return err
	}
	return os.WriteFile(files.OutputPath(config.OutputDirFlag), []byte("."), os.ModePerm)
}

// 尝试移除输出目录，并检测是否是正确的非空或由 RellwNote 管理的目录
func removeOutputDirectory() error {
	if !files.IsDir(config.BuildOutput) {
		return errors.New("output directory isn't a folder")
	}
	empty := files.IsEmptyDir(config.BuildOutput)
	if !empty {
		flag := files.IsFile(files.OutputPath(config.OutputDirFlag))
		if !flag {
			return errors.New("output directory is not empty")
		}
	}

	return os.RemoveAll(config.BuildOutput)
}
