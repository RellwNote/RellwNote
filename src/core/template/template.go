package template

import (
	"bytes"
	"fmt"
	"github.com/RellwNote/RellwNote/TOCGenerator"
	"github.com/RellwNote/RellwNote/config"
	"github.com/RellwNote/RellwNote/log"
	"github.com/RellwNote/RellwNote/models"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var LastLoadedTemplate *template.Template

// LibraryData 定义了 html 模版内需要的数据格式。
type LibraryData struct {
	Directory   models.TOCItem
	LibraryName string
}

// NewLibraryData 会根据当前参数创建新的 LibraryData
func NewLibraryData() LibraryData {
	return LibraryData{
		Directory:   TOCGenerator.GetTOCFromFile(path.Join(config.LibraryPath, config.SummaryFileName)),
		LibraryName: config.LibraryName,
	}
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
	if config.LibraryFileExists("index.html") {
		res, err := config.ReadLibraryFile("index.html")
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
