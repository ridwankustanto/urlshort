package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	urlshort "github.com/ridwankustanto/urlshort"
)

func main() {
	yamlLocation := flag.String("yamlLocation", "pathToUrls.yaml", "location of yaml file that contain path and url")
	jsonLocation := flag.String("jsonLocation", "pathToUrls.json", "location of json file that contain path and url")
	flag.Parse()

	// Parse yaml
	yamlByte, err := parseFile(*yamlLocation)
	if err != nil {
		log.Fatalln(err)
		return
	}

	// Parse json
	jsonByte, err := parseFile(*jsonLocation)
	if err != nil {
		log.Fatalln(err)
		return
	}

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	// 	yaml := `- path: /urlshort
	//   url: https://github.com/gophercises/urlshort
	// - path: /urlshort-final
	//   url: https://github.com/gophercises/urlshort/tree/solution`
	yamlHandler, err := urlshort.YAMLHandler(yamlByte, mapHandler)
	if err != nil {
		panic(err)
	}

	jsonHandler, err := urlshort.JSONHandler(jsonByte, yamlHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func parseFile(location string) ([]byte, error) {

	f, err := os.Open(location)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var buf bytes.Buffer
	io.Copy(&buf, f)

	return buf.Bytes(), nil
}
