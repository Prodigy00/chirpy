package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

type chirpReq struct {
	Body string `json:"body"`
}

type validateErrResponse struct {
	Error string `json:"error"`
}

type chirpRes struct {
	CleanedBody string `json:"cleaned_body"`
}

func NewApiConfig() *apiConfig {
	return &apiConfig{}
}

func (cfg *apiConfig) ToJSON(w http.ResponseWriter, v any) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	incSrvHits := func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(incSrvHits)
}
func bannedWords() map[string]int {
	return map[string]int{
		"kerfuffle": 1,
		"sharbert":  1,
		"fornax":    1,
	}
}
func (cfg *apiConfig) ValidateChirp(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)

	var validated chirpReq

	if err := decoder.Decode(&validated); err != nil {
		somethingWentWrong := validateErrResponse{
			Error: "something went wrong.",
		}
		w.WriteHeader(http.StatusBadRequest)
		cfg.ToJSON(w, somethingWentWrong)
		return
	}

	if len(validated.Body) > 140 {
		chirpTooLongErr := validateErrResponse{
			Error: "chirp is too long",
		}
		w.WriteHeader(http.StatusBadRequest)
		cfg.ToJSON(w, chirpTooLongErr)
		return
	}

	//check for banned words
	profMap := bannedWords()

	formatted := make([]string, len(validated.Body))

	for _, w := range strings.Split(validated.Body, " ") {
		word := w
		if _, ok := profMap[strings.ToLower(w)]; ok {
			word = "****"
		}
		formatted = append(formatted, " "+word)
	}

	cb := ""
	for _, wr := range formatted {
		cb += wr
	}
	res := chirpRes{
		CleanedBody: strings.TrimSpace(cb),
	}
	w.WriteHeader(http.StatusOK)
	cfg.ToJSON(w, res)
	return
}

func (cfg *apiConfig) Reset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "<html>\n  <body>\n   <p>Reset Ok!</p>\n  </body>\n</html>")
}

func (cfg *apiConfig) FileServerHits(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<html>\n  <body>\n    <h1>Welcome, Chirpy Admin</h1>\n    <p>Chirpy has been visited %d times!</p>\n  </body>\n</html>", cfg.fileserverHits.Load())
}

func main() {
	ns := http.NewServeMux()
	fs := http.FileServer(http.Dir("."))
	c := NewApiConfig()
	stripped := http.StripPrefix("/app/", fs)

	ns.Handle("/app/", c.middlewareMetricsInc(stripped))
	ns.Handle("/app/assets/logo.png", stripped)

	ns.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		w.Write([]byte("OK"))
	})

	ns.HandleFunc("GET /admin/metrics", c.FileServerHits)
	ns.HandleFunc("POST /admin/reset", c.Reset)
	server := http.Server{
		Handler: ns,
		Addr:    ":8080",
	}

	ns.HandleFunc("POST /api/validate_chirp", c.ValidateChirp)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
