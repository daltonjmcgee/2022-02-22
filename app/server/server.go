package server

import (
	"encoding/json"
	"fmt"
	"goserve/config"
	"goserve/helpers"
	"goserve/httpErrorHandler"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var conf map[string]string = config.ReturnConfig()

func methodCheck(w http.ResponseWriter, r *http.Request, method string) {
	if r.Method != method {
		httpErrorHandler.Handle405(w, method)
		return
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) {

	methodCheck(w, r, "GET")

	var templates []string
	templateFiles, _ := ioutil.ReadDir(conf["templatesPath"])
	for _, file := range templateFiles {
		templates = append(templates, conf["templatesPath"]+fmt.Sprintf("/%s", file.Name()))
	}

	path := r.URL.Path

	// Checking for pattern used for dynamic pages and return 404 if found.
	// We don't want anyone grabbing that un-rendered page.
	matched, _ := regexp.Match(`\[\w+\]`, []byte(path))
	if matched {
		httpErrorHandler.Handle404(w)
		return
	}

	if path == "/" {
		jsonBytes, err := helpers.LoadFile(conf["databasePath"])

		if err != nil {
			// This should probably throw a different error
			fmt.Fprintf(w, conf["databasePath"])
			return
		}

		jsonMap := map[string][]interface{}{}
		json.Unmarshal([]byte(jsonBytes), &jsonMap)

		files := append([]string{conf["publicPath"] + "/index.html"}, templates...)
		t, err := template.ParseFiles(files...)
		if err != nil {
			httpErrorHandler.Handle500(w)
		}
		t.Execute(w, jsonMap)
	} else {
		files := append([]string{conf["publicPath"] + path + ".html"}, templates...)
		t, err := template.ParseFiles(files...)

		if err == nil {
			t.Execute(w, nil)
		} else {

			// Take the URI, strip off the last data to get the directory, then
			// get a list of all files in the directory to be looped over.
			// If the directory doesn't exist throw a 404.
			fileName := strings.Split(path, "/")
			queryableValue := &fileName[len(fileName)-1]
			directory := strings.Join(fileName[:len(fileName)-1], "/")
			directoryFiles, err := ioutil.ReadDir(conf["publicPath"] + fmt.Sprintf("/%s", directory))

			if err != nil {
				httpErrorHandler.Handle404(w)
				return
			}

			// Loop over all files in the directory and see if the template name matches
			// any of the keys in the JSON data provided. If so, serve the first one found.
			for _, file := range directoryFiles {

				// Skip subdirectories.
				if file.IsDir() {
					continue
				}

				// Skip file if it doesn't match the template format.
				isFile, _ := regexp.Match(`\[\w+\]`, []byte(file.Name()))

				if !isFile {
					continue
				}

				jsonBytes, err := helpers.LoadFile(conf["databasePath"])

				if err != nil {
					// This should probably throw a different error
					httpErrorHandler.Handle404(w)
					return
				}

				jsonMap := map[string][]interface{}{}
				queryKey := regexp.MustCompile(`\[|\]`).Split(file.Name(), -1)[1]

				json.Unmarshal([]byte(jsonBytes), &jsonMap)

				for _, val := range jsonMap["data"] {
					for key, value := range val.(map[string]interface{}) {
						if key == queryKey && *queryableValue == value {
							fullDirectory := conf["publicPath"] + fmt.Sprintf("%s/%s", directory, file.Name())
							files := append([]string{fullDirectory}, templates...)
							t, err := template.ParseFiles(files...)
							if err != nil {
								fmt.Println(err)
								httpErrorHandler.Handle500(w)
							} else {
								t.Execute(w, val)
							}
							return
						}
					}
				}
			}
			httpErrorHandler.Handle404(w)
			return
		}
	}
}

func handleStatic(w http.ResponseWriter, r *http.Request) {
	methodCheck(w, r, "GET")

	path := r.URL.Path
	if path != "/static/" {
		http.ServeFile(w, r, conf["publicPath"]+path)
	} else {
		httpErrorHandler.Handle404(w)
	}
}

func doNothing(w http.ResponseWriter, r *http.Request) {}

func CreateServer() {
	http.HandleFunc("/", handleGet)
	http.HandleFunc("/static/", handleStatic)
	http.HandleFunc("/favicon.ico", doNothing)
	log.Fatal(http.ListenAndServe(":"+conf["port"], nil))
}
