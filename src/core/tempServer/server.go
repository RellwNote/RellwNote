package tempServer

import (
	"bytes"
	"fmt"
	"github.com/RellwNote/RellwNote/TOCGenerator"
	"github.com/RellwNote/RellwNote/config"
	"github.com/RellwNote/RellwNote/log"
	"github.com/RellwNote/RellwNote/models"
	"github.com/RellwNote/RellwNote/template"
	"math"
	"math/rand/v2"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func templatesPage() (res []byte, state int) {
	for _, t := range template.LoadFromDir(config.TemplateDir).Templates() {
		res = append(res, []byte(t.Name()+"\n")...)
	}
	return res, 200
}

func contentPage() (res []byte, state int) {
	content := TOCGenerator.GetSummaryFileToByte(config.LibraryPath, config.SummaryFileName)
	directory := TOCGenerator.ParseSummaryByte(content)

	contentTemplateData := struct {
		Directory models.TOCItem
	}{
		Directory: directory,
	}

	var buf bytes.Buffer
	err := template.LoadFromDir(config.TemplateDir).ExecuteTemplate(&buf, "content.gohtml", contentTemplateData)
	if err != nil {
		return []byte(err.Error()), 500
	}
	return buf.Bytes(), 200
}

func indexPage() ([]byte, int) {
	if config.LibraryFileExists("index.html") {
		res, err := config.ReadLibraryFile("index.html")
		if err != nil {
			return []byte("load library index.html error"), 500
		}
		return res, 200
	}

	var buf bytes.Buffer
	err := template.LoadFromDir(config.TemplateDir).ExecuteTemplate(&buf, "index/index.gohtml", nil)
	if err != nil {
		return []byte(err.Error()), 500
	}
	return buf.Bytes(), 200
}

func staticFile(path string) (res []byte, state int) {
	filePath := filepath.Join(config.LibraryPath, path)
	file, err := os.Stat(filePath)
	if err != nil {
		return nil, 404
	}
	if file.IsDir() {
		log.Error.Println("尝试访问静态文件", filePath, "，但这是一个目录")
		return nil, 404
	}

	fileByte, err := os.ReadFile(filePath)
	if err != nil {
		log.Error.Println("尝试访问静态文件", filePath, "出现错误：", err.Error())
		return nil, 500
	}

	return fileByte, 200
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	if config.ServerDebugDelay > 0 {
		intTime := int64(math.Ceil(config.ServerDebugDelay * float64(time.Second)))
		time.Sleep(time.Duration(rand.Int64N(intTime)))
	}
	var response []byte
	var state int
	urlPath := r.URL.Path
	if urlPath == "/templates" {
		response, state = templatesPage()
	} else if urlPath == "/content" {
		response, state = contentPage()
	} else if urlPath == "/" || urlPath == "/index.html" {
		response, state = indexPage()
	} else {
		response, state = staticFile(urlPath)
	}

	w.WriteHeader(state)
	_, err := w.Write(response)
	if err != nil {
		log.Error.Println(err.Error())
		return
	}
}

func Start() {
	server := fmt.Sprintf("%s:%d", config.ServerHost, config.ServerPort)
	fmt.Printf("Server setup on http://%s\n", server)
	err := http.ListenAndServe(server, http.HandlerFunc(httpHandler))
	if err != nil {
		log.Error.Println(err.Error())
	}
}
