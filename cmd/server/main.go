package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/jbdoto/go-study-tool/internal/challenges"
	"github.com/jbdoto/go-study-tool/internal/questions"
	"github.com/jbdoto/go-study-tool/internal/server"
)

//go:embed all:templates
var templateFS embed.FS

func main() {
	chapters, err := questions.LoadChapters("questions")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Loaded %d chapter(s)\n", len(chapters))

	challengeSets, err := challenges.LoadChallengeSets("challenges")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Loaded %d challenge set(s)\n", len(challengeSets))

	srv, err := server.New(chapters, challengeSets, templateFS)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", srv.Routes()))
}
