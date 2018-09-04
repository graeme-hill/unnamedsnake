package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.HandleFunc("/start", Start)
	http.HandleFunc("/move", Move)
	http.HandleFunc("/end", End)

	port := os.Getenv("PORT")
	if port == "" {
		port = "1234"
	}

	// Add filename into logging messages
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Printf("Running server on port %s...\n", port)
	http.ListenAndServe(":"+port, LoggingHandler(http.DefaultServeMux))
}
