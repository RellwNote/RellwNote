package tempServer

import "net/http"

func Start() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

	})
	_ = http.ListenAndServe(":8080", nil)
}
