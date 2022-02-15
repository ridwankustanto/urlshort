package urlshort

import (
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	//	TODO: Implement this...
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		shortPath := r.URL.String()

		// find match path
		for path, url := range pathsToUrls {
			if path == shortPath {
				http.Redirect(rw, r, url, http.StatusMovedPermanently)
				return
			}
		}

		// fallback if nothing found
		fallback.ServeHTTP(rw, r)

	})

	return mux.ServeHTTP
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// TODO: Implement this...

	yamlData, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}

	yamlMap := buildMap(yamlData)

	return MapHandler(yamlMap, fallback), nil
}

func parseYAML(yamlString []byte) ([]map[string]string, error) {
	yamlData := []map[string]string{}

	if err := yaml.Unmarshal(yamlString, &yamlData); err != nil {
		return nil, fmt.Errorf("can't unmarshal yaml data with error %v", err)
	}

	return yamlData, nil
}

func buildMap(yamlData []map[string]string) map[string]string {
	yamlMap := make(map[string]string)
	for _, data := range yamlData {
		yamlMap[data["path"]] = data["url"]
	}
	return yamlMap
}
