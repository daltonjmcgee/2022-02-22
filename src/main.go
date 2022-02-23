package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var config map[string]string = returnConfig()

func loadFile(fileName string) (string, error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "File not found", err
	} else {
		return string(bytes), err
	}
}

func handleDynamic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path != "/static/" {
		http.ServeFile(w, r, "../public"+path)
	} else {
		http.NotFound(w, r)
	}
}

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

func handleUri(w http.ResponseWriter, r *http.Request) {
	var templates []string
	templateFiles, _ := ioutil.ReadDir("../templates")
	for _, file := range templateFiles {
		templates = append(templates, fmt.Sprintf("../templates/%s", file.Name()))
	}

	if r.Method != "GET" {
		handle405(w, r.Method)
		return
	}

	path := r.URL.Path

	// Checking for pattern used for dynamic pages and return 404 if found.
	// We don't want anyone grabbing that un-rendered page.
	matched, _ := regexp.Match(`\[\w+\]`, []byte(path))
	if matched {
		handle404(w)
		return
	}

	if path == "/" {
		jsonBytes, err := loadFile(config["databasePath"])

		if err != nil {
			// This should probably throw a different error
			fmt.Fprintf(w, config["databasePath"])
			return
		}

		jsonMap := map[string][]interface{}{}
		json.Unmarshal([]byte(jsonBytes), &jsonMap)

		files := append([]string{"../public/index.html"}, templates...)
		t, err := template.ParseFiles(files...)
		if err != nil {
			handle500(w)
		}
		t.Execute(w, jsonMap)
	} else {
		files := append([]string{fmt.Sprintf("../public%s.html", path)}, templates...)
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
			directoryFiles, err := ioutil.ReadDir(fmt.Sprintf("../public/%s", directory))

			if err != nil {
				handle404(w)
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

				jsonBytes, err := loadFile("../noSQL.json")

				if err != nil {
					// This should probably throw a different error
					handle404(w)
					return
				}

				jsonMap := map[string][]interface{}{}
				queryKey := regexp.MustCompile(`\[|\]`).Split(file.Name(), -1)[1]

				json.Unmarshal([]byte(jsonBytes), &jsonMap)

				for _, val := range jsonMap["data"] {
					for key, value := range val.(map[string]interface{}) {
						if key == queryKey && *queryableValue == value {
							fullDirectory := fmt.Sprintf("../public%s/%s", directory, file.Name())
							files := append([]string{fullDirectory}, templates...)
							t, err := template.ParseFiles(files...)
							if err != nil {
								fmt.Println(err)
								handle500(w)
							} else {
								t.Execute(w, val)
							}
							return
						}
					}
				}
			}
			handle404(w)
			return
		}
	}
}

func handleStatic(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		handle405(w, r.Method)
		return
	}
	path := r.URL.Path
	if path != "/static/" {
		http.ServeFile(w, r, "../public"+path)
	} else {
		handle404(w)
	}
}

func doNothing(w http.ResponseWriter, r *http.Request) {}

func main() {
	http.HandleFunc("/", handleUri)
	http.HandleFunc("/static/", handleStatic)
	http.HandleFunc("/favicon.ico", doNothing)
	log.Fatal(http.ListenAndServe(":"+config["port"], nil))
}
