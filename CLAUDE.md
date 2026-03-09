# Go Study Tool

Web-based drill app for Go interview prep, built with Go stdlib `net/http`, HTMX, and PicoCSS. No database — stateless URL-param progression for quizzes, localStorage for stats.

## Architecture

- **Module**: `github.com/jbdoto/go-study-tool` (Go 1.25, depends on `gopkg.in/yaml.v3`)
- **Server**: stdlib `net/http` with Go 1.22+ route patterns, HTMX for partials, PicoCSS for styling
- **Templates**: embedded via `//go:embed all:templates`, parsed at startup
- **Data**: YAML files loaded at startup, no runtime persistence

## Project Structure

```
cmd/server/
  main.go                  — entry point, embeds templates
  templates/
    layout.html            — quiz layout (PicoCSS + HTMX)
    layout_challenge.html  — challenge layout (adds Monaco Editor CDN)
    index.html             — home page (chapters + challenge sets)
    quiz.html              — MC question page with progress bar
    result.html            — HTMX partial for answer feedback + stats tracking
    done.html              — quiz completion page
    stats.html             — stats dashboard (reads localStorage)
    challenge.html         — code challenge page with Monaco editor
    challenge_hint.html    — HTMX partial for progressive hints
    challenge_solution.html — HTMX partial for reference solution
    challenge_check.html   — HTMX partial for syntax check result
internal/
  questions/               — MC quiz types and YAML loader
  challenges/              — Code challenge types and YAML loader
  server/                  — HTTP handlers, routing, server setup
questions/                 — MC question YAML files (one per chapter)
challenges/                — Code challenge YAML files (one per topic)
```

## Two Modes

### MC Quizzes (`/chapter/{slug}`, `/quiz/all`)
- Shuffle question IDs into URL params (`?order=id1,id2&i=0`), walk through linearly
- 4 choices shown per question (1 correct + 3 random wrong from 5+)
- HTMX partial swap for answer feedback
- Stats tracked in localStorage (overall + per-chapter %)

### Code Challenges (`/challenges/{slug}`)
- Monaco Editor loads scaffold code from a JSON `<script>` element
- Buttons: Check Syntax (server-side `go/parser`), Show Hint (progressive), Show Solution
- No compilation/execution — syntax check only via AST parsing

## Stats
- Tracked client-side in localStorage key `quizStats`
- `/stats` page renders from localStorage via JS DOM construction
- Shows overall % and per-chapter % with color coding and drill links

## Question YAML Format

```yaml
chapter: chapter-slug
title: "Chapter Title"
source_url: "https://..."
questions:
  - id: "slug-mc-1"      # must be globally unique
    type: multiple_choice
    difficulty: easy|medium|hard
    text: "Question text"
    correct: "The correct answer"
    wrong:                 # exactly 5 wrong answers
      - "Wrong answer 1"
      - "Wrong answer 2"
      - "Wrong answer 3"
      - "Wrong answer 4"
      - "Wrong answer 5"
    explanation: "Why the correct answer is right."
    source_anchor: "section-heading-slug"  # optional
```

## Challenge YAML Format

```yaml
challenge_set: set-slug
title: "Set Title"
challenges:
  - id: "unique-id"
    type: fill_blank|free_form
    difficulty: easy|medium|hard
    title: "Challenge Title"
    description: |
      What to implement.
    scaffold: |
      // Starter code with TODOs
    solution: |
      // Complete reference solution
    hints:
      - "Hint 1"
      - "Hint 2"
    key_concepts:
      - "Concept 1"
```

## Key Conventions

- All question/challenge IDs must be globally unique across all YAML files
- Wrong answers should be similar in length to the correct answer (avoid length-based guessing)
- `source_anchor` uses GitBook slug format: "Hello, world... again" → `hello-world...-again`
- Templates are embedded at compile time — must rebuild after template changes

## Running

```sh
go run ./cmd/server
# Open http://localhost:8080
```

## Testing

```sh
go test ./...
```

## Adding Content

### New MC chapter
1. Create `questions/{slug}.yaml` following the format above
2. Use unique ID prefix (e.g., `myslug-mc-1`)
3. 5 wrong answers per question, balanced lengths
4. Difficulty spread: ~30% easy, ~50% medium, ~20% hard

### New challenge set
1. Create `challenges/{slug}.yaml` following the format above
2. Include scaffold code, complete solution, 2+ hints, key concepts