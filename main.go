package main

import (
	"fmt"
	"net/http"

	"github.com/mikerybka/util"
)

func main() {
	port := util.EnvVar("PORT", "2069")
	addr := ":" + port
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
