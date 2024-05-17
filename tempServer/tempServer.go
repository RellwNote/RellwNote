package tempServer

import (
	"bytes"
	"fmt"
	"github.com/RellwNote/RellwNote/config"
	"github.com/RellwNote/RellwNote/directoryGenerator"
	"github.com/RellwNote/RellwNote/log"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var rootTemplate *template.Template

func init() {
	loadAllTemplate("template")
}

func loadAllTemplate(root string) {
	rootTemplate = template.New("main")
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
		_, err = rootTemplate.Parse(fmt.Sprintf(`{{define "%s"}}%s{{end}}`, key, read))
		if err != nil {
			log.Error.Println(err.Error())
		}
		return nil
	})
}

func printTemplates() (res []byte) {
	for _, t := range rootTemplate.Templates() {
		res = append(res, []byte(t.Name()+"\n")...)
	}
	return res
}

func printContent() (res []byte, err error) {
	filePath := config.GetPublicConfig.Directory.FilePath
	content := directoryGenerator.GetSummaryFileToByte(filePath)
	directory := directoryGenerator.ParseSummaryByte(content)

	contentTemplateData := struct {
		Directory directoryGenerator.Directory
	}{
		Directory: directory,
	}

	var buf bytes.Buffer
	err = rootTemplate.ExecuteTemplate(&buf, "html/content.gohtml", contentTemplateData)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	loadAllTemplate("template")

	var response []byte
	var err error
	switch r.URL.Path {
	case "/templates":
		response = printTemplates()
		break
	case "/content":
		response, err = printContent()
	}

	// write error
	if err != nil {
		log.Error.Println(err.Error())
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	// write response
	_, err = w.Write(response)
	if err != nil {
		log.Error.Println(err.Error())
		return
	}
}

func Start() {
	err := http.ListenAndServe(config.GetPublicConfig.Template.Server.Port, http.HandlerFunc(httpHandler))
	if err != nil {
		log.Error.Println(err.Error())
	}
}
