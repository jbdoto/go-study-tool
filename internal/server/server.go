package server

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/jbdoto/go-study-tool/internal/challenges"
	"github.com/jbdoto/go-study-tool/internal/questions"
)

type Server struct {
	chapters      map[string]questions.Chapter
	challengeSets map[string]challenges.ChallengeSet
	templates     map[string]*template.Template
}

func New(chapters []questions.Chapter, challengeSets []challenges.ChallengeSet, templateFS fs.FS) (*Server, error) {
	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
	}

	pages := []string{"index.html", "quiz.html", "done.html"}
	templates := make(map[string]*template.Template, len(pages))
	for _, page := range pages {
		t, err := template.New("").Funcs(funcMap).ParseFS(templateFS, "templates/layout.html", "templates/"+page)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", page, err)
		}
		templates[page] = t
	}

	// result.html is a standalone partial (no layout)
	resultTmpl, err := template.New("").ParseFS(templateFS, "templates/result.html")
	if err != nil {
		return nil, fmt.Errorf("parse result.html: %w", err)
	}
	templates["result.html"] = resultTmpl

	// challenge.html uses its own layout
	challengeTmpl, err := template.New("").Funcs(funcMap).ParseFS(templateFS, "templates/layout_challenge.html", "templates/challenge.html")
	if err != nil {
		return nil, fmt.Errorf("parse challenge.html: %w", err)
	}
	templates["challenge.html"] = challengeTmpl

	// standalone partials for challenge HTMX responses
	for _, partial := range []string{"challenge_hint.html", "challenge_solution.html", "challenge_check.html"} {
		t, err := template.New("").ParseFS(templateFS, "templates/"+partial)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", partial, err)
		}
		templates[partial] = t
	}

	chMap := make(map[string]questions.Chapter, len(chapters))
	for _, ch := range chapters {
		chMap[ch.Slug] = ch
	}

	csMap := make(map[string]challenges.ChallengeSet, len(challengeSets))
	for _, cs := range challengeSets {
		csMap[cs.Slug] = cs
	}

	return &Server{chapters: chMap, challengeSets: csMap, templates: templates}, nil
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.handleIndex)
	mux.HandleFunc("GET /quiz/all", s.handleQuizAll)
	mux.HandleFunc("POST /quiz/all/submit", s.handleSubmitAll)
	mux.HandleFunc("GET /chapter/{name}", s.handleQuiz)
	mux.HandleFunc("POST /chapter/{name}/submit", s.handleSubmit)
	mux.HandleFunc("GET /challenges/{set}", s.handleChallengeIndex)
	mux.HandleFunc("GET /challenges/{set}/solve", s.handleChallengeSolve)
	mux.HandleFunc("POST /challenges/{set}/hint", s.handleChallengeHint)
	mux.HandleFunc("POST /challenges/{set}/solution", s.handleChallengeSolution)
	mux.HandleFunc("POST /challenges/{set}/check", s.handleChallengeCheck)
	return mux
}
