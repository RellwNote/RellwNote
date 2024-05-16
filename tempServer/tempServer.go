package tempServer

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

var Templates = make(map[string]*template.Template)

func init() {
	loadAllTemplate("template")
}

// 加载一个目录下的全部模版，模版在 root 文件夹中的相对路径作为 key 存在 Templates 变量中。
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
		Templates[key] = template.Must(template.New(path).Parse(string(read)))
		return nil
	})
}

func Start() {
	http.HandleFunc("/templates", func(w http.ResponseWriter, r *http.Request) {
		for s := range Templates {
			w.Write([]byte(s + "\n"))
		}
	})
	_ = http.ListenAndServe(":8080", nil)
}
