package admin

import (
	"fmt"
	"goserve/helpers"
	"goserve/httpErrorHandler"
	"net/http"
)

var adminPath string = "./admin/public/"

func handleGet(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	file, err := helpers.LoadFile(adminPath + path[1:] + ".html")

	if err != nil {
		httpErrorHandler.Handle404(w)
	} else {
		fmt.Fprintf(w, file)
	}
}

func handleStatic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fmt.Println(path)
	if path != "/admin/static/" {
		http.ServeFile(w, r, adminPath+path)
	} else {
		httpErrorHandler.Handle404(w)
	}
}

func AdminPanel() {
	http.HandleFunc("/admin", handleGet)
	http.HandleFunc("/admin/", handleGet)
	http.HandleFunc("/admin/static/", handleStatic)
}
