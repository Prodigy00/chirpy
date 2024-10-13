package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Prodigy00/chirpy/internal/db"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	queries        *db.Queries
}

type chirpReq struct {
	Body   string `json:"body"`
	UserID string `json:"user_id"`
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
		http.Error(w, fmt.Sprintf("error encoding json:%s\n", err.Error()), http.StatusInternalServerError)
		return
	}
	return
}

func (cfg *apiConfig) ToStruct(w http.ResponseWriter, r *http.Request, v any) {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		http.Error(w, fmt.Sprintf("error decoding json:%s\n", err.Error()), http.StatusBadRequest)
		return
	}
	return
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

type createUserReq struct {
	Email string `json:"email"`
}

func (cfg *apiConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var createUser createUserReq

	cfg.ToStruct(w, r, &createUser)

	u, err := cfg.queries.CreateUser(r.Context(), createUser.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		cfg.ToJSON(w, validateErrResponse{
			Error: fmt.Sprintf("error creating user:%v", err),
		})
	}

	w.WriteHeader(http.StatusCreated)
	user := db.User{
		ID:        u.ID,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}

	cfg.ToJSON(w, user)
	return
}

func (cfg *apiConfig) ValidateChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var validated db.CreateChirpParams

	cfg.ToStruct(w, r, &validated)

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

	cleanedParams := db.CreateChirpParams{
		Body:   strings.TrimSpace(cb),
		UserID: validated.UserID,
	}

	ch, err := cfg.queries.CreateChirp(r.Context(), cleanedParams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		cfg.ToJSON(w, fmt.Sprintf("error creating chirp:%v", err))
	}

	created := db.Chirp{
		ID:        ch.ID,
		Body:      ch.Body,
		CreatedAt: ch.CreatedAt,
		UpdatedAt: ch.UpdatedAt,
		UserID:    ch.UserID,
	}

	w.WriteHeader(http.StatusOK)
	cfg.ToJSON(w, created)
	return
}

func (cfg *apiConfig) Reset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	err := cfg.queries.DeleteUsers(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		cfg.ToJSON(w, fmt.Errorf("error deleting users: %v\n", err))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "<html>\n  <body>\n   <p>Reset Ok!</p>\n  </body>\n</html>")
	return
}

func (cfg *apiConfig) FileServerHits(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<html>\n  <body>\n    <h1>Welcome, Chirpy Admin</h1>\n    <p>Chirpy has been visited %d times!</p>\n  </body>\n</html>", cfg.fileserverHits.Load())
}

func main() {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")

	database, sqlErr := sql.Open("postgres", dbUrl)
	if sqlErr != nil {
		log.Fatalf("error accessing db:%v\n", sqlErr)
	}

	dbQueries := db.New(database)

	ns := http.NewServeMux()
	fs := http.FileServer(http.Dir("."))
	c := NewApiConfig()
	c.queries = dbQueries

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
	ns.HandleFunc("POST /api/users", c.CreateUser)
	ns.HandleFunc("POST /api/chirps", c.ValidateChirp)
	server := http.Server{
		Handler: ns,
		Addr:    ":8080",
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
