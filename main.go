package main

import (
	"fmt"
	"net/http"
	"log"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricMaker(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	ret := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	w.Write([]byte(ret))
}

func (cfg *apiConfig) metricReset(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0) 
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("reset"))
}

func health(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func main() {
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}
	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	
	mux.HandleFunc("/metrics", apiCfg.metricMaker)
	mux.HandleFunc("/reset", apiCfg.metricReset)
	mux.HandleFunc("/healthz", health)

	var srv http.Server
	srv.Handler = mux
	srv.Addr = ":8080"

	fmt.Println("Starting Server")
	log.Fatal(srv.ListenAndServe())
}