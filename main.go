package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func NewApiConfig() *apiConfig {
	return &apiConfig{}
}
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	incSrvHits := func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(incSrvHits)
}

func (cfg *apiConfig) Reset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}

func (cfg *apiConfig) FileServerHits(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
}

func main() {
	ns := http.NewServeMux()
	fs := http.FileServer(http.Dir("."))
	c := NewApiConfig()
	stripped := http.StripPrefix("/app/", fs)

	ns.Handle("/app/", c.middlewareMetricsInc(stripped))
	ns.Handle("/app/assets/logo.png", stripped)

	ns.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		w.Write([]byte("OK"))
	})

	ns.HandleFunc("/metrics", c.FileServerHits)
	ns.HandleFunc("/reset", c.Reset)
	server := http.Server{
		Handler: ns,
		Addr:    ":8080",
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
