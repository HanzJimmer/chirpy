package main

import (
	"net/http"
)

func main() {
	const port = "8080"
	serveMux := http.NewServeMux()

	serveMux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	serveMux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	localServer := http.Server{Handler: serveMux, Addr: ":" + port}
	localServer.ListenAndServe()
}
