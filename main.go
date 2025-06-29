package main

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mikerybka/util"
)

//go:embed main.css
var css string

//go:embed main.js
var js string

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "2069"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		data, ok := getData(filepath.Join(util.HomeDir(), "schema-cafe/data", r.URL.Path))
		if !ok {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
  <link rel="stylesheet" href="/main.css">
</head>
<body>
  <div id="root"></div>
  <script id="data" type="application/json">%s</script>
  <script src="/main.js"></script>
</body>
</html>`, data)
	})

	http.HandleFunc("/main.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		fmt.Fprint(w, js)
	})

	http.HandleFunc("/main.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		fmt.Fprint(w, css)
	})

	fmt.Printf("Server listening on port %s...\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func getData(path string) ([]byte, bool) {
	fi, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, false
		}
		panic(err)
	}
	if fi.IsDir() {
		entries, err := os.ReadDir(path)
		if err != nil {
			panic(err)
		}
		data := []DirEntry{}
		for _, e := range entries {
			entry := DirEntry{
				Name: e.Name(),
			}
			if e.IsDir() {
				entry.Type = "dir"
			} else {
				entry.Type = "schema"
			}
			data = append(data, entry)
		}
		b, err := json.Marshal(Response{
			Type:  "dir",
			Value: data,
		})
		if err != nil {
			panic(err)
		}
		return b, true
	} else {
		s := &Schema{}
		f, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		err = json.NewDecoder(f).Decode(s)
		if err != nil {
			panic(err)
		}
		b, err := json.Marshal(Response{
			Type:  "schema",
			Value: s,
		})
		if err != nil {
			panic(err)
		}
		return b, true
	}
}

type Response struct {
	Type  string `json:"type"`
	Value any    `json:"value"`
}

type Schema struct {
	Fields []Field `json:"fields"`
}

type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type DirEntry struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
