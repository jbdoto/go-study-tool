package server

import (
	"math/rand/v2"
	"net/http"
	"strconv"
	"strings"

	"github.com/jbdoto/go-study-tool/internal/questions"
)

type presentedChoice struct {
	Text    string
	Index   int
}

// pickChoices selects 3 random wrong answers + the correct answer, shuffles them,
// and returns the choices and the index of the correct one.
func pickChoices(q *questions.Question) ([]string, int) {
	// Shuffle wrong answers and pick up to 3
	wrong := make([]string, len(q.Wrong))
	copy(wrong, q.Wrong)
	rand.Shuffle(len(wrong), func(i, j int) { wrong[i], wrong[j] = wrong[j], wrong[i] })
	n := 3
	if len(wrong) < n {
		n = len(wrong)
	}
	choices := make([]string, 0, n+1)
	choices = append(choices, q.Correct)
	choices = append(choices, wrong[:n]...)

	// Shuffle all choices
	rand.Shuffle(len(choices), func(i, j int) { choices[i], choices[j] = choices[j], choices[i] })

	correctIdx := 0
	for i, c := range choices {
		if c == q.Correct {
			correctIdx = i
			break
		}
	}
	return choices, correctIdx
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	chapters := make([]questions.Chapter, 0, len(s.chapters))
	totalQuestions := 0
	for _, ch := range s.chapters {
		chapters = append(chapters, ch)
		totalQuestions += len(ch.Questions)
	}

	type challengeSetInfo struct {
		Title      string
		Slug       string
		Challenges int
	}
	sets := make([]challengeSetInfo, 0, len(s.challengeSets))
	for _, cs := range s.challengeSets {
		sets = append(sets, challengeSetInfo{Title: cs.Title, Slug: cs.Slug, Challenges: len(cs.Challenges)})
	}

	data := struct {
		Chapters       []questions.Chapter
		TotalQuestions int
		ChallengeSets  []challengeSetInfo
	}{chapters, totalQuestions, sets}
	s.templates["index.html"].ExecuteTemplate(w, "layout", data)
}

func (s *Server) allQuestions() []questions.Question {
	var all []questions.Question
	for _, ch := range s.chapters {
		all = append(all, ch.Questions...)
	}
	return all
}

func (s *Server) findQuestion(id string) *questions.Question {
	for _, ch := range s.chapters {
		for i := range ch.Questions {
			if ch.Questions[i].ID == id {
				return &ch.Questions[i]
			}
		}
	}
	return nil
}

func (s *Server) chapterForQuestion(id string) *questions.Chapter {
	for _, ch := range s.chapters {
		for _, q := range ch.Questions {
			if q.ID == id {
				return &ch
			}
		}
	}
	return nil
}

func (s *Server) sourceURLForQuestion(id string) string {
	ch := s.chapterForQuestion(id)
	if ch == nil {
		return ""
	}
	url := ch.SourceURL
	for _, q := range ch.Questions {
		if q.ID == id && q.SourceAnchor != "" {
			url += "#" + q.SourceAnchor
			break
		}
	}
	return url
}

func (s *Server) handleQuizAll(w http.ResponseWriter, r *http.Request) {
	order := r.URL.Query().Get("order")
	idx, _ := strconv.Atoi(r.URL.Query().Get("i"))

	if order == "" {
		all := s.allQuestions()
		ids := make([]string, len(all))
		for i, q := range all {
			ids[i] = q.ID
		}
		rand.Shuffle(len(ids), func(i, j int) { ids[i], ids[j] = ids[j], ids[i] })
		http.Redirect(w, r, "/quiz/all?order="+strings.Join(ids, ",")+"&i=0", http.StatusFound)
		return
	}

	ids := strings.Split(order, ",")
	total := len(ids)

	if idx >= total {
		data := struct {
			Title    string
			Total    int
			QuizPath string
		}{"All Chapters", total, "/quiz/all"}
		s.templates["done.html"].ExecuteTemplate(w, "layout", data)
		return
	}

	q := s.findQuestion(ids[idx])
	if q == nil {
		http.NotFound(w, r)
		return
	}

	choices, correctIdx := pickChoices(q)
	s.renderQuiz(w, "All Chapters", q, choices, correctIdx, idx, total, order, "/quiz/all/submit")
}

func (s *Server) handleSubmitAll(w http.ResponseWriter, r *http.Request) {
	questionID := r.FormValue("question_id")
	order := r.FormValue("order")
	idx, _ := strconv.Atoi(r.FormValue("i"))

	q := s.findQuestion(questionID)
	if q == nil {
		http.Error(w, "question not found", http.StatusNotFound)
		return
	}

	ids := strings.Split(order, ",")
	nextIdx := idx + 1
	nextURL := "/quiz/all?order=" + order + "&i=" + strconv.Itoa(nextIdx)
	isLast := nextIdx >= len(ids)

	s.renderSubmit(w, r, q, nextURL, isLast, "/quiz/all")
}

func (s *Server) handleQuiz(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	ch, ok := s.chapters[name]
	if !ok {
		http.NotFound(w, r)
		return
	}

	order := r.URL.Query().Get("order")
	idx, _ := strconv.Atoi(r.URL.Query().Get("i"))

	if order == "" {
		ids := make([]string, len(ch.Questions))
		for i, q := range ch.Questions {
			ids[i] = q.ID
		}
		rand.Shuffle(len(ids), func(i, j int) { ids[i], ids[j] = ids[j], ids[i] })
		http.Redirect(w, r, "/chapter/"+name+"?order="+strings.Join(ids, ",")+"&i=0", http.StatusFound)
		return
	}

	ids := strings.Split(order, ",")
	total := len(ids)

	if idx >= total {
		data := struct {
			Title    string
			Total    int
			QuizPath string
		}{ch.Title, total, "/chapter/" + name}
		s.templates["done.html"].ExecuteTemplate(w, "layout", data)
		return
	}

	var q *questions.Question
	for i := range ch.Questions {
		if ch.Questions[i].ID == ids[idx] {
			q = &ch.Questions[i]
			break
		}
	}
	if q == nil {
		http.NotFound(w, r)
		return
	}

	choices, correctIdx := pickChoices(q)
	s.renderQuiz(w, ch.Title, q, choices, correctIdx, idx, total, order, "/chapter/"+name+"/submit")
}

func (s *Server) handleSubmit(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	ch, ok := s.chapters[name]
	if !ok {
		http.NotFound(w, r)
		return
	}

	questionID := r.FormValue("question_id")
	order := r.FormValue("order")
	idx, _ := strconv.Atoi(r.FormValue("i"))

	var q *questions.Question
	for i := range ch.Questions {
		if ch.Questions[i].ID == questionID {
			q = &ch.Questions[i]
			break
		}
	}
	if q == nil {
		http.Error(w, "question not found", http.StatusNotFound)
		return
	}

	ids := strings.Split(order, ",")
	nextIdx := idx + 1
	nextURL := "/chapter/" + name + "?order=" + order + "&i=" + strconv.Itoa(nextIdx)
	isLast := nextIdx >= len(ids)

	s.renderSubmit(w, r, q, nextURL, isLast, "/chapter/"+name)
}

func (s *Server) renderQuiz(w http.ResponseWriter, title string, q *questions.Question, choices []string, correctIdx, idx, total int, order, submitPath string) {
	data := struct {
		Title      string
		Question   questions.Question
		Choices    []string
		CorrectIdx int
		Index      int
		Total      int
		Order      string
		SubmitPath string
	}{
		Title:      title,
		Question:   *q,
		Choices:    choices,
		CorrectIdx: correctIdx,
		Index:      idx,
		Total:      total,
		Order:      order,
		SubmitPath: submitPath,
	}
	s.templates["quiz.html"].ExecuteTemplate(w, "layout", data)
}

func (s *Server) renderSubmit(w http.ResponseWriter, r *http.Request, q *questions.Question, nextURL string, isLast bool, quizPath string) {
	selected := r.FormValue("answer")
	correct := selected == q.Correct

	var chapterTitle, chapterSlug string
	if ch := s.chapterForQuestion(q.ID); ch != nil {
		chapterTitle = ch.Title
		chapterSlug = ch.Slug
	}

	data := struct {
		Correct       bool
		CorrectAnswer string
		Explanation   string
		SourceURL     string
		NextURL       string
		IsLast        bool
		QuizPath      string
		ChapterTitle  string
		ChapterSlug   string
	}{
		Correct:       correct,
		CorrectAnswer: q.Correct,
		Explanation:   q.Explanation,
		SourceURL:     s.sourceURLForQuestion(q.ID),
		NextURL:       nextURL,
		IsLast:        isLast,
		QuizPath:      quizPath,
		ChapterTitle:  chapterTitle,
		ChapterSlug:   chapterSlug,
	}
	s.templates["result.html"].ExecuteTemplate(w, "result.html", data)
}
