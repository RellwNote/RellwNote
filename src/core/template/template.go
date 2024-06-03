package template

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"rellwnote/core/TOCGenerator"
	"rellwnote/core/config"
	"rellwnote/core/library"
	"rellwnote/core/log"
	"rellwnote/core/models"
	"strings"
)

var LastLoadedTemplate *template.Template

// LibraryData 定义了 html 模版内需要的数据格式。
type LibraryData struct {
	Directory       models.TOCItem
	LibraryName     string
	FaviconFileName string
}

// NewLibraryData 会根据当前参数创建新的 LibraryData
func NewLibraryData() LibraryData {
	res := LibraryData{
		Directory:   TOCGenerator.GetTOCFromFile(path.Join(config.LibraryPath, config.SummaryFileName)),
		LibraryName: config.LibraryName,
	}
	if name, has := library.GetIconFileName(); has {
		res.FaviconFileName = name
	} else {
		res.FaviconFileName = "favicon.svg"
	}
	return res
}

// LoadFromDir 会读取一个目录中的全部模版文件，等同于重新加载模版
func LoadFromDir(root string) *template.Template {
	LastLoadedTemplate = template.New("main").Funcs(CustomFuncMap)
	_ = filepath.Walk(root, func(filePath string, info os.FileInfo, err error) error {
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
		key = key[strings.Index(key, "/")+1:]
		_, err = LastLoadedTemplate.Parse(fmt.Sprintf(`{{define "%s"}}%s{{end}}`, key, read))
		if err != nil {
			log.Error.Printf("模版 %s 中存在错误：%s\n", filePath, err.Error())
		}
		return nil
	})

	return LastLoadedTemplate
}

func BuildContentPage() ([]byte, error) {
	var buf bytes.Buffer
	err := LoadFromDir(config.TemplateDir).ExecuteTemplate(&buf, "content.gohtml", NewLibraryData())
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func BuildIndexPage() ([]byte, error) {
	if library.FileExists("index.html") {
		res, err := library.ReadFile("index.html")
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	var buf bytes.Buffer
	err := LoadFromDir(config.TemplateDir).ExecuteTemplate(&buf, "index/index.gohtml", NewLibraryData())
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
