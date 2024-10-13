package main

import (
	"log"
	"net/http"
)

func main() {
	ns := http.NewServeMux()
	fs := http.FileServer(http.Dir("."))

	stripped := http.StripPrefix("/app/", fs)

	ns.Handle("/app/", stripped)
	ns.Handle("/app/assets/logo.png", stripped)

	ns.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		w.Write([]byte("OK"))
	})

	server := http.Server{
		Handler: ns,
		Addr:    ":8080",
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
