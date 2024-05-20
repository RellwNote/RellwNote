package tempServer

import (
	"bytes"
	"github.com/RellwNote/RellwNote/config"
	"github.com/RellwNote/RellwNote/directoryGenerator"
	"github.com/RellwNote/RellwNote/log"
	"github.com/RellwNote/RellwNote/template"
	"net/http"
)

func templatesPage() (res []byte) {
	for _, t := range template.LoadFromDir(config.GetPublicConfig.Template.FilePath).Templates() {
		res = append(res, []byte(t.Name()+"\n")...)
	}
	return res
}

func contentPage() (res []byte, err error) {
	filePath := config.GetPublicConfig.Directory.FilePath
	content := directoryGenerator.GetSummaryFileToByte(filePath)
	directory := directoryGenerator.ParseSummaryByte(content)

	contentTemplateData := struct {
		Directory directoryGenerator.Directory
	}{
		Directory: directory,
	}

	var buf bytes.Buffer
	err = template.LoadFromDir(config.GetPublicConfig.Template.FilePath).ExecuteTemplate(&buf, "html/content.gohtml", contentTemplateData)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	var response []byte
	var err error
	switch r.URL.Path {
	case "/templates":
		response = templatesPage()
		break
	case "/content":
		response, err = contentPage()
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
	err := http.ListenAndServe(config.GetPublicConfig.Temp.Server.Port, http.HandlerFunc(httpHandler))
	if err != nil {
		log.Error.Println(err.Error())
	}
}