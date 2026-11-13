package main

import (
	"net/http"
	"strconv"
	"sync/atomic"
)

func main() {
	const port = "8080"
	apiCfg := apiConfig{fileserverHits: atomic.Int32{}}
	serveMux := http.NewServeMux()

	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	serveMux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	serveMux.HandleFunc("GET /metrics", apiCfg.getHitsHandler)
	serveMux.HandleFunc("POST /reset", apiCfg.resetHitsHandler)

	localServer := http.Server{Handler: serveMux, Addr: ":" + port}
	localServer.ListenAndServe()
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) getHitsHandler(w http.ResponseWriter, r *http.Request) {
	numHits := int(cfg.fileserverHits.Load())
	w.Write([]byte("Hits: " + strconv.Itoa(numHits) + "\n"))
}

func (cfg *apiConfig) resetHitsHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Swap(0)
}
