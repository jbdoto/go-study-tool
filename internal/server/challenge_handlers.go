package server

import (
	"encoding/json"
	"go/parser"
	"go/token"
	"math/rand/v2"
	"net/http"
	"strconv"
	"strings"

	"github.com/jbdoto/go-study-tool/internal/challenges"
)

func (s *Server) findChallenge(set *challenges.ChallengeSet, id string) *challenges.CodeChallenge {
	for i := range set.Challenges {
		if set.Challenges[i].ID == id {
			return &set.Challenges[i]
		}
	}
	return nil
}

func (s *Server) handleChallengeIndex(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("set")
	cs, ok := s.challengeSets[slug]
	if !ok {
		http.NotFound(w, r)
		return
	}

	ids := make([]string, len(cs.Challenges))
	for i, c := range cs.Challenges {
		ids[i] = c.ID
	}
	rand.Shuffle(len(ids), func(i, j int) { ids[i], ids[j] = ids[j], ids[i] })
	http.Redirect(w, r, "/challenges/"+slug+"/solve?order="+strings.Join(ids, ",")+"&i=0", http.StatusFound)
}

func (s *Server) handleChallengeSolve(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("set")
	cs, ok := s.challengeSets[slug]
	if !ok {
		http.NotFound(w, r)
		return
	}

	order := r.URL.Query().Get("order")
	idx, _ := strconv.Atoi(r.URL.Query().Get("i"))
	ids := strings.Split(order, ",")
	total := len(ids)

	if idx >= total {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	challenge := s.findChallenge(&cs, ids[idx])
	if challenge == nil {
		http.NotFound(w, r)
		return
	}

	scaffoldJSON, _ := json.Marshal(challenge.Scaffold)

	data := struct {
		SetTitle     string
		SetSlug      string
		Challenge    challenges.CodeChallenge
		ScaffoldJSON string
		Index        int
		Total        int
		Order        string
		IsLast       bool
	}{
		SetTitle:     cs.Title,
		SetSlug:      cs.Slug,
		Challenge:    *challenge,
		ScaffoldJSON: string(scaffoldJSON),
		Index:        idx,
		Total:        total,
		Order:        order,
		IsLast:       idx+1 >= total,
	}
	s.templates["challenge.html"].ExecuteTemplate(w, "layout", data)
}

func (s *Server) handleChallengeHint(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("set")
	cs, ok := s.challengeSets[slug]
	if !ok {
		http.NotFound(w, r)
		return
	}

	challengeID := r.FormValue("challenge_id")
	shown, _ := strconv.Atoi(r.FormValue("shown"))

	challenge := s.findChallenge(&cs, challengeID)
	if challenge == nil {
		http.NotFound(w, r)
		return
	}

	if shown >= len(challenge.Hints) {
		w.Write([]byte(`<div class="hint-box"><strong>No more hints.</strong></div>`))
		return
	}

	data := struct {
		Number int
		Text   string
	}{shown + 1, challenge.Hints[shown]}
	s.templates["challenge_hint.html"].ExecuteTemplate(w, "challenge_hint.html", data)
}

func (s *Server) handleChallengeSolution(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("set")
	cs, ok := s.challengeSets[slug]
	if !ok {
		http.NotFound(w, r)
		return
	}

	challengeID := r.FormValue("challenge_id")
	challenge := s.findChallenge(&cs, challengeID)
	if challenge == nil {
		http.NotFound(w, r)
		return
	}

	data := struct {
		Solution string
	}{challenge.Solution}
	s.templates["challenge_solution.html"].ExecuteTemplate(w, "challenge_solution.html", data)
}

func (s *Server) handleChallengeCheck(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")

	fset := token.NewFileSet()
	_, err := parser.ParseFile(fset, "code.go", code, parser.AllErrors)

	data := struct {
		OK    bool
		Error string
	}{err == nil, ""}
	if err != nil {
		data.Error = err.Error()
	}
	s.templates["challenge_check.html"].ExecuteTemplate(w, "challenge_check.html", data)
}
