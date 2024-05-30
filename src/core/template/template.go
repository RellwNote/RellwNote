package template

import (
	"fmt"
	"github.com/RellwNote/RellwNote/log"
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
