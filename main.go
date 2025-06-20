package main

import (
	"fmt"
	"net/http"

	"github.com/mikerybka/util"
)

func main() {
	port := util.EnvVar("PORT", "2069")
	addr := ":" + port
	err := http.ListenAndServe(addr, &SchemaCafe{"data"})
	if err != nil {
		fmt.Println(err)
		return
	}
}

type SchemaCafe struct {
	DataDir string
}

func (cafe *SchemaCafe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, cafe.DataDir)
}
