package main

import (
	"fmt"
	"net/http"

	"github.com/mikerybka/util"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hi")
	})
	port := util.EnvVar("PORT", "2069")
	addr := ":" + port
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
