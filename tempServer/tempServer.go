package tempServer

import (
	"bytes"
	"fmt"
	"github.com/RellwNote/RellwNote/config"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var rootTemplate = template.New("main")

func init() {
	loadAllTemplate("template")
}

func loadAllTemplate(root string) {
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		read, _ := os.ReadFile(path)
		key := path[len(root)+1:]
		key = strings.ReplaceAll(key, "\\", "/")
		rootTemplate.Parse(fmt.Sprintf(`{{define "%s"}}%s{{end}}`, key, read))
		return nil
	})
}

func Start() {
	templateFile := config.GetPublicConfig.Template.FilePath
	templateServerPort := config.GetPublicConfig.Template.Server.Port

	http.HandleFunc(templateFile, func(w http.ResponseWriter, r *http.Request) {

		for _, v := range rootTemplate.Templates() {
			w.Write([]byte(v.Name() + "\n"))
		}
	})
	http.HandleFunc("/content", func(writer http.ResponseWriter, request *http.Request) {
		var buf bytes.Buffer
		err := rootTemplate.ExecuteTemplate(&buf, "html/content.gohtml", nil)
		if err != nil {
			println(err.Error())
			return
		}
		writer.Write(buf.Bytes())
	})

	_ = http.ListenAndServe(templateServerPort, nil)
}
