package template

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"rellwnote/core/config"
	"rellwnote/core/extensions"
	"rellwnote/core/files"
	"rellwnote/core/library"
	"rellwnote/core/library/toc"
	"rellwnote/core/log"
	"strings"
)

var LastLoadedTemplate *template.Template

// LibraryData 定义了 html 模版内需要的数据格式。
type LibraryData struct {
	Directory       toc.Item
	LibraryName     string
	FaviconFileName string
	Extensions      []extensions.Extension
}

// NewLibraryData 会根据当前参数创建新的 LibraryData
func NewLibraryData() LibraryData {
	TOCFormFile, err := toc.GetTOCFromFile(files.LibraryPath(config.SummaryFileName))
	if err != nil {
		log.Error.Println(err)
	}
	res := LibraryData{
		Directory:   TOCFormFile,
		LibraryName: config.LibraryName,
	}
	if name, has := library.GetIconFileName(); has {
		res.FaviconFileName = name
	} else {
		res.FaviconFileName = "favicon.svg"
	}

	for _, v := range extensions.LoadAll() {
		res.Extensions = append(res.Extensions, v)
	}

	return res
}

// Load 会读取一个目录中的全部模版文件，等同于重新加载模版
func Load() *template.Template {
	LastLoadedTemplate = template.New("main").Funcs(CustomFuncMap)
	startPath := files.ProgramPath(config.TemplateDir)
	startPath, _ = filepath.Abs(startPath)
	_ = filepath.Walk(startPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.ToLower(filepath.Ext(filePath)) != ".gohtml" {
			return nil
		}
		read, _ := os.ReadFile(filePath)
		key := filepath.ToSlash(filePath)
		key = key[len(startPath)+1:]
		_, err = LastLoadedTemplate.Parse(fmt.Sprintf(`{{define "%s"}}%s{{end}}`, key, read))
		if err != nil {
			log.Error.Printf("template %s has error：%s\n", filePath, err.Error())
		}
		return nil
	})

	return LastLoadedTemplate
}

func BuildContentPage() ([]byte, error) {
	var buf bytes.Buffer
	err := Load().ExecuteTemplate(&buf, "content.gohtml", NewLibraryData())
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func BuildIndexPage() ([]byte, error) {
	if files.IsFile(files.LibraryPath("index.html")) {
		res, err := os.ReadFile(files.LibraryPath("index.html"))
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	var buf bytes.Buffer
	err := Load().ExecuteTemplate(&buf, "index/index.gohtml", NewLibraryData())
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
