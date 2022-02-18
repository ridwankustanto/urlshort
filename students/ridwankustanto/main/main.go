package main

import (
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	urlshort "github.com/ridwankustanto/urlshort"
)

func main() {
	yamlLocation := flag.String("yamlLocation", "pathToUrls.yaml", "location of yaml file that contain path and url")
	jsonLocation := flag.String("jsonLocation", "pathToUrls.json", "location of json file that contain path and url")
	sqliteLocation := flag.String("sqliteLocation", "pathToUrls.db", "location of sqlite db file that contain path and url")
	flag.Parse()

	// connect to sqlite3
	dbConnection, err := connectSqlite(*sqliteLocation)
	if err != nil {
		log.Fatalln(err)
		return
	}

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

	sqliteHandler, err := urlshort.SqliteHandler(dbConnection, jsonHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", sqliteHandler)
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

func connectSqlite(location string) (*sql.DB, error) {
	sqliteDatabase, err := sql.Open("sqlite3", location)
	if err != nil {
		return nil, err
	}
	// defer sqliteDatabase.Close()

	return sqliteDatabase, nil
}
