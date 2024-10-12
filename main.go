package main

import (
	"log"
	"net/http"
)

func main() {
	ns := http.NewServeMux()
	server := http.Server{
		Handler: ns,
		Addr:    ":8080",
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}