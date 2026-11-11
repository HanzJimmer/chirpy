package main

import (
	"net/http"
)

func main() {
	const port = "8080"
	serveMux := http.NewServeMux()
	localServer := http.Server{Handler: serveMux, Addr: ":" + port}
	serveMux.Handle("/", http.FileServer(http.Dir(".")))
	localServer.ListenAndServe()
}
