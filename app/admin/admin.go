package admin

import (
	"fmt"
	"goserve/helpers"
	"goserve/httpErrorHandler"
	"net/http"
)

func handleAdmin(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	file, err := helpers.LoadFile("./admin/public/" + path[1:] + ".html")

	if err != nil {
		httpErrorHandler.Handle404(w)
	} else {
		fmt.Fprintf(w, file)
	}
}

func handleStatic(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		httpErrorHandler.Handle405(w, r.Method)
		return
	}

	path := r.URL.Path
	if path == "/admin/style.css" {
		http.ServeFile(w, r, "./admin/public/style.css")
	} else {
		httpErrorHandler.Handle404(w)
	}
}

func AdminPanel() {
	http.HandleFunc("/admin", handleAdmin)
	http.HandleFunc("/admin/", handleAdmin)
	http.HandleFunc("/admin/style.css", handleStatic)
}
