package tempServer

import (
	"bytes"
	"github.com/RellwNote/RellwNote/TOCGenerator"
	"github.com/RellwNote/RellwNote/config"
	"github.com/RellwNote/RellwNote/log"
	"github.com/RellwNote/RellwNote/models"
	"github.com/RellwNote/RellwNote/template"
	"net/http"
	"os"
	"path/filepath"
)

func templatesPage() (res []byte) {
	for _, t := range template.LoadFromDir(config.TemplateDir).Templates() {
		res = append(res, []byte(t.Name()+"\n")...)
	}
	return res
}

func contentPage() (res []byte, err error) {

	content := TOCGenerator.GetSummaryFileToByte("test/", "Summary.md")
	directory := TOCGenerator.ParseSummaryByte(content)

	contentTemplateData := struct {
		Directory models.TOCItem
	}{
		Directory: directory,
	}

	var buf bytes.Buffer
	err = template.LoadFromDir(config.TemplateDir).ExecuteTemplate(&buf, "content.gohtml", contentTemplateData)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func serveStaticFiles(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join(config.LibraryPath, r.URL.Path)

	file, err := os.Stat(filePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error.Println(err)
		return
	}
	if file.IsDir() {
		log.Error.Println("访问的文件是一个文件夹，非md文件", filePath)
		_, err := w.Write([]byte("404 访问的文件是一个文件夹，非md文件"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	fileByte, err := os.ReadFile(filePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error.Println(err)
		return
	}
	_, err = w.Write(fileByte)
	if err != nil {
		log.Error.Println(err)
	}
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	var response []byte
	var err error

	urlPath := r.URL.Path
	if urlPath == "/templates" {
		response = templatesPage()
		_, err := w.Write(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error.Println(err)
			return
		}
		return
	}

	if urlPath == "/content" {
		response, err = contentPage()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error.Println(err)
			return
		}
		_, err := w.Write(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error.Println(err)
			return
		}
		return
	}

	if len(urlPath) > 3 && urlPath[len(urlPath)-3:len(urlPath)] == ".md" {
		serveStaticFiles(w, r)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	return
}

func Start() {
	err := http.ListenAndServe(config.PreviewServer, http.HandlerFunc(httpHandler))
	if err != nil {
		log.Error.Println(err.Error())
	}
}
