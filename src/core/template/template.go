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
	"path/filepath"
	"strings"
)

var LastLoadedTemplate *template.Template

// 读取一个目录中的全部模版
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
	content := TOCGenerator.GetSummaryFileToByte(config.LibraryPath, config.SummaryFileName)
	directory := TOCGenerator.ParseSummaryByte(content)

	contentTemplateData := struct {
		Directory models.TOCItem
	}{
		Directory: directory,
	}

	var buf bytes.Buffer
	err := LoadFromDir(config.TemplateDir).ExecuteTemplate(&buf, "content.gohtml", contentTemplateData)
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
	err := LoadFromDir(config.TemplateDir).ExecuteTemplate(&buf, "index/index.gohtml", nil)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
