package main

import (
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/mikerybka/util"
)

//go:embed favicon.ico
var favicon []byte

type API struct {
	StorageURL  string
	StorageUser string
	StoragePass string
}

func (api *API) NewRequest(method string, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, api.StorageURL+path, body)
	req.SetBasicAuth(api.StorageUser, api.StoragePass)
	return req, err
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := api.NewRequest(r.Method, strings.TrimPrefix(r.URL.Path, "/api"), r.Body)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
}

func main() {
	port := util.RequireEnvVar("PORT")
	jsURL := util.EnvVar("JS_URL", "https://brass.dev/main.js")
	cssURL := util.EnvVar("CSS_URL", "https://brass.dev/main.css")
	storageURL := util.RequireEnvVar("STORAGE_URL")
	storageUser := os.Getenv("STORAGE_USER")
	storagePass := os.Getenv("STORAGE_PASS")
	api := &API{
		StorageURL:  storageURL,
		StorageUser: storageUser,
		StoragePass: storagePass,
	}

	http.Handle("/api/", api)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Get data from storage
		req, err := api.NewRequest("GET", r.URL.Path, nil)
		if err != nil {
			panic(err)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}

		// If err, return that
		if res.StatusCode != 200 {
			http.Error(w, res.Status, res.StatusCode)
			return
		}

		// Return data with HTML wrapper if request is from a browser
		isHTML := util.Accept(r, "text/html")
		if isHTML {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, `<!DOCTYPE html>
<html>
	<head>
	  <link rel="stylesheet" href="%s">
	</head>
	<body id="body">
	  <div id="app">Please wait...</div>
	  <script id="data" type="application/json">`, cssURL)
		}
		_, err = io.Copy(w, res.Body)
		if err != nil {
			panic(err)
		}
		if isHTML {
			fmt.Fprintf(w, `</script>
	  <script src="%s"></script>
	</body>
</html>`, jsURL)
		}
	})

	http.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/x-icon")
		w.Write(favicon)
	})

	fmt.Printf("Server listening on port %s...\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
