package main
import (
	"fmt"
	"net/http"
	"log"
	"sync/atomic"
	//"encoding/json"
	//"strings"
	"os"
	"github.com/Syn-gularity/httpserv/internal/database"
	"github.com/joho/godotenv"
	"database/sql"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

/*func (cfg *apiConfig) metricMaker(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	ret := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	w.Write([]byte(ret))
}*/

func (cfg *apiConfig) adminMetric(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	ret := fmt.Sprintf("<html>\n  <body>\n    <h1>Welcome, Chirpy Admin</h1>\n    <p>Chirpy has been visited %d times!</p>\n  </body>\n</html>", cfg.fileserverHits.Load())
	w.Write([]byte(ret))
}

func (cfg *apiConfig) metricReset(w http.ResponseWriter, req *http.Request) {
	if os.Getenv("PLATFORM") != "dev"{
		respondWithError(w, 403, "FORBIDDEN", nil)
		return
	}
	cfg.fileserverHits.Store(0)
	err := cfg.db.DeleteAllUsers(req.Context())
	if err != nil{
		respondWithError(w, 400, "Shit happened during DB Access: %v", err)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("reset"))
}

func health(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

/*func chirpValidator(w http.ResponseWriter, r *http.Request){
    type parameters struct {
        // these tags indicate how the keys in the JSON should be mapped to the struct fields
        // the struct fields must be exported (start with a capital letter) if you want them parsed
        Msg string `json:"body"`
    }
	type returnMsg struct {
		Valid bool `json:"valid"` 
		CleanedBody string `json:"cleaned_body"`
	}

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if len(params.Msg) > 140{
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	} 
	cleanedMsg := badWordDisplacer(params.Msg)
	respondWithJSON(w, http.StatusOK, returnMsg{
		Valid: true,
		CleanedBody: cleanedMsg,
	})
}

func badWordDisplacer(text string) string{
	ret := ""
	const displacer = "****"
	var badWords = [3]string{"kerfuffle", "sharbert", "fornax"}
	lowered := strings.ToLower(text)
	splitText := strings.Split(lowered," ")
	splitTextOriginal := strings.Split(text," ")
	for idx, word := range splitText{
		bad := false
		for _, badWord := range badWords{
			if word == badWord{
				bad = true
			}
		}
		if bad {
			apnd := fmt.Sprintf(" %s", displacer)
			ret += apnd
		} else {
			apnd := fmt.Sprintf(" %s", splitTextOriginal[idx])
			ret += apnd
		}
	}
	ret = strings.TrimSpace(ret)
	return ret
}*/



func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	db_conn, err := sql.Open("postgres", dbURL)
	if err != nil{
		log.Fatalf("fatal error opening database: %s", err)
	}
	dbQueries := database.New(db_conn)
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db: dbQueries,
	}
	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	
	mux.HandleFunc("GET /admin/metrics", apiCfg.adminMetric)
	mux.HandleFunc("POST /admin/reset", apiCfg.metricReset)

	mux.HandleFunc("GET /api/healthz", health)
	//mux.HandleFunc("POST /api/validate_chirp", chirpValidator)
	mux.HandleFunc("POST /api/chirps", apiCfg.handleChirps)
	mux.HandleFunc("POST /api/users", apiCfg.handleCreateUser)

	var srv http.Server
	srv.Handler = mux
	srv.Addr = ":8080"

	fmt.Println("Starting Server")
	log.Fatal(srv.ListenAndServe())
}