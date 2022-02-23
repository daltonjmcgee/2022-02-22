package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func handle404(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	file, err := loadFile("../public/404.html")
	if err != nil {
		fmt.Fprintf(w, "Some dipshit deleted the default 404 and didn't replace it. At any rate, your page wasn't found.")
	} else {
		fmt.Fprintf(w, file)
	}
}

func handle405(w http.ResponseWriter, method string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	data, _ := json.Marshal(map[string]string{
		"message": method + " is not allowed on this endpoint. Try something else.",
		"status":  "405",
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func handle500(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	file, err := loadFile("../public/500.html")
	if err != nil {
		fmt.Fprintf(w, "Some dipshit deleted the default 404 and didn't replace it. At any rate, there was a server error. Try again later.")
	} else {
		fmt.Fprintf(w, file)
	}
}
