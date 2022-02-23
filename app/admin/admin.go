package admin

import (
	"fmt"
	"goserve/helpers"
	"goserve/httpErrorHandler"
	"net/http"
)

func handleAdmin(w http.ResponseWriter, r *http.Request) {
	file, err := helpers.LoadFile("./admin/admin.html")

	if err != nil {
		httpErrorHandler.Handle404(w)
	} else {
		fmt.Fprintf(w, file)
	}
}

func AdminPanel() {
	http.HandleFunc("/admin", handleAdmin)
}
