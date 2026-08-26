package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/frozenf1sh/fish-interview/internal/content"
	"github.com/frozenf1sh/fish-interview/internal/web"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP listen address")
	contentDir := flag.String("content", "content", "content directory")
	flag.Parse()

	catalog, err := content.Load(*contentDir)
	if err != nil {
		log.Fatalf("load content: %v", err)
	}

	server, err := web.New(catalog)
	if err != nil {
		log.Fatalf("build web server: %v", err)
	}
	log.Printf("Fish Interview listening on http://localhost%s", *addr)
	if err := http.ListenAndServe(*addr, server); err != nil {
		log.Fatal(err)
	}
}
