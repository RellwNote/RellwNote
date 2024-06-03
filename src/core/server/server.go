package server

import (
	"fmt"
	"math"
	"math/rand/v2"
	"net/http"
	"os"
	"path"
	"rellwnote/core/config"
	"rellwnote/core/library"
	"rellwnote/core/log"
	"rellwnote/core/template"
	"time"
)

func templatesPage() (res []byte, state int) {
	for _, t := range template.LoadFromDir(config.TemplateDir).Templates() {
		res = append(res, []byte(t.Name()+"\n")...)
	}
	return res, 200
}

func contentPage() (res []byte, state int) {
	build, err := template.BuildContentPage()
	if err != nil {
		return []byte(err.Error()), 500
	}
	return build, 200
}

func indexPage() ([]byte, int) {
	build, err := template.BuildIndexPage()
	if err != nil {
		return []byte(err.Error()), 500
	}
	return build, 200
}

func favicon() ([]byte, int) {
	if name, has := library.GetIconFileName(); has {
		icon, err := library.ReadFile(name)
		if err != nil {
			return nil, 500
		}
		return icon, 200
	}
	file, err := os.ReadFile(path.Join(config.TemplateDir, "favicon.svg"))
	if err != nil {
		return nil, 500
	}
	return file, 200
}

func staticFile(url string) (res []byte, state int) {
	if library.FileExists(url) {
		file, err := library.ReadFile(url)
		if err != nil {
			return nil, 500
		}
		return file, 200
	}
	return nil, 404
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
	} else if urlPath == "/content.html" {
		response, state = contentPage()
	} else if len(urlPath) >= 12 && urlPath[:9] == "/favicon." {
		response, state = favicon()
		if urlPath[8:12] == ".svg" {
			w.Header().Set("Content-Type", "image/svg+xml")
		} else {
			w.Header().Set("Content-Type", "image")
		}
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
