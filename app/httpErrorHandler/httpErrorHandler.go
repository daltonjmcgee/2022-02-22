package httpErrorHandler

import (
	"encoding/json"
	"fmt"
	"goserve/helpers"
	"net/http"
)

func Handle404(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	file, err := helpers.LoadFile("./httpErrorHandler/404.html")
	if err != nil {
		fmt.Fprintf(w, "Some dipshit deleted the default 404 and didn't replace it. At any rate, your page wasn't found.")
	} else {
		fmt.Fprintf(w, file)
	}
}

func Handle405(w http.ResponseWriter, method string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	data, _ := json.Marshal(map[string]string{
		"message": method + " is not allowed on this endpoint. Try something else.",
		"status":  "405",
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func Handle500(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	file, err := helpers.LoadFile("./httpErrorHandler/500.html")
	if err != nil {
		fmt.Fprintf(w, "Some dipshit deleted the default 404 and didn't replace it. At any rate, there was a server error. Try again later.")
	} else {
		fmt.Fprintf(w, file)
	}
}
