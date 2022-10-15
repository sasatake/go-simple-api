package main

import (
	"fmt"
	"log"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	name := r.URL.Query().Get("name")

	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "set name parameter.")
		return
	}

	message := fmt.Sprintf("Hello %s\n", name)
	log.Printf("parameter name value is %s", name)
	fmt.Fprint(w, message)
}
