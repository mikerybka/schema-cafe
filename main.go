package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"os"
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
		fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head>
  <title>Simple Page</title>
  <link rel="stylesheet" href="/main.css">
</head>
<body>
  <div id="root"></div>
  <script src="/main.js"></script>
</body>
</html>`)
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
