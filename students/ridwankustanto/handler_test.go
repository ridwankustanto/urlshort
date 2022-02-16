package urlshort

import (
	"fmt"
	"testing"
)

func TestParseYAML(t *testing.T) {
	var test = []struct {
		yamlByte []byte
		want     []map[string]string
	}{
		{
			yamlByte: []byte(`- path: /google
  url: https://google.com`),
			want: []map[string]string{
				{"path": "/google", "url": "https://google.com"},
			},
		},
		{
			yamlByte: []byte(`- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution`),
			want: []map[string]string{
				{"path": "/urlshort", "url": "https://github.com/gophercises/urlshort"},
				{"path": "/urlshort-final", "url": "https://github.com/gophercises/urlshort/tree/solution"},
			},
		},
	}

	for _, tt := range test {
		testName := string(tt.yamlByte)
		t.Run(testName, func(t *testing.T) {
			ans, err := parseYAML(tt.yamlByte)
			if err != nil {
				t.Fatal(err)
			}

			for i, val := range ans {
				if val["path"] != tt.want[i]["path"] && val["url"] != tt.want[i]["url"] {
					t.Errorf("got %v, want %v", val, tt.want[i])
				}
			}
		})
	}
}

func BencmarkParseYAML(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseYAML([]byte(`- path: /urlshort
		url: https://github.com/gophercises/urlshort`))
	}
}

func TestBuildMap(t *testing.T) {
	var test = []struct {
		mapData []map[string]string
		want    map[string]string
	}{
		{
			mapData: []map[string]string{
				{"path": "/urlshort-godoc", "url": "https://godoc.org/github.com/gophercises/urlshort"},
			},
			want: map[string]string{
				"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
			},
		},
		{
			mapData: []map[string]string{
				{"path": "/urlshort-godoc", "url": "https://godoc.org/github.com/gophercises/urlshort"},
				{"path": "/github", "url": "https://github.com"},
			},
			want: map[string]string{
				"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
				"/github":         "https://github.com",
			},
		},
	}

	for _, tt := range test {
		testName := fmt.Sprintf("map with %v length", len(tt.mapData))

		t.Run(testName, func(t *testing.T) {
			ans := buildMap(tt.mapData)

			if ans["path"] != tt.want["path"] && ans["url"] != tt.want["url"] {
				t.Errorf("got %v, want %v", tt.want, ans)
			}
		})

	}

}

func BencmarkBuildMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buildMap([]map[string]string{
			{"path": "/urlshort-godoc", "url": "https://godoc.org/github.com/gophercises/urlshort"},
		})
	}
}
