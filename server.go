package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"
)

func serveUI() {
	mux := http.NewServeMux()

	// setup server with in memory user sessions
	server := Server{
		jobs: make(map[string]*Job),
	}

	mux.HandleFunc("GET /metrics", server.getMetrics)
	mux.HandleFunc("POST /metrics", server.createJob)
	mux.HandleFunc("DELETE /metrics", server.cancelJob)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "./static/ui.html")
	})

	handler := cors.AllowAll().Handler(mux)

	log.Println("Serving static page at '/'")
	log.Panicln(http.ListenAndServe(":8000", handler))
}
