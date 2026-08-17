// Command the-chat-server serves the built frontend as a single-page application.
package main

import (
	"flag"
	"io/fs"
	"log"
	"net/http"

	"the-chat/internal/server"
	"the-chat/web"
)

func main() {
	addr := flag.String("addr", "0.0.0.0:8080", "address to listen on")
	flag.Parse()

	dist, err := fs.Sub(web.Dist, "dist")
	if err != nil {
		log.Fatalf("locate embedded frontend: %v", err)
	}

	log.Printf("listening on %s", *addr)
	if err := http.ListenAndServe(*addr, server.New(dist)); err != nil {
		log.Fatal(err)
	}
}
