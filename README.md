# Go Study Tool

A web-based study tool for Go interview prep. Multiple-choice quizzes and hands-on code challenges with an in-browser Monaco editor.

Built with Go stdlib, HTMX, and PicoCSS. No frameworks, no database, no build step.

## Features

- **Multiple-choice quizzes** — 250+ questions across 30 chapters covering Go fundamentals, gRPC, observability, CLI design, testing patterns, and more
- **Code challenges** — 17 hands-on coding exercises with Monaco Editor, progressive hints, reference solutions, and server-side syntax checking via `go/parser`
- **Stats tracking** — per-chapter accuracy percentages stored in localStorage, with drill-down links to weak areas
- **Shuffle mode** — quiz all questions across every chapter in randomized order

## Quick Start

```sh
go run ./cmd/server
```

Open [http://localhost:8080](http://localhost:8080).

## How It Works

**Quizzes**: Question IDs are shuffled into URL parameters on first visit. Each "Next Question" increments the index. Four answer choices are shown (1 correct + 3 random from 5+ wrong options). Answer feedback is swapped in via HTMX without a full page reload.

**Code Challenges**: Monaco Editor loads scaffold code with TODO comments. Write your solution, click "Check Syntax" for instant `go/parser` feedback, reveal hints progressively, or compare against the reference solution.

**Stats**: Every answer is recorded to localStorage. The `/stats` page shows your overall accuracy and a per-chapter breakdown with color-coded scores and links to drill weak chapters.

## Content

### MC Question Banks

| Category | Chapters | Questions |
|----------|----------|-----------|
| Go Fundamentals | Hello World, Variables, Iteration, Arrays, Structs, Maps, Pointers, etc. | ~130 |
| Concurrency | Goroutines, Channels & Select, Context, Sync | ~40 |
| Interview Prep | gRPC, Observability, CLI Design, Testing Patterns, Retries | ~60 |
| Go Syntax | "Which statement is correct?" format across all Go features | 20 |

### Code Challenge Sets

| Topic | Challenges |
|-------|-----------|
| gRPC Client Connection | 2 |
| Retry Interceptor | 2 |
| Context & Timeout | 2 |
| CLI Commands | 2 |
| gRPC Testing | 2 |
| OpenTelemetry Tracing | 2 |
| Structured Logging | 2 |
| Graceful Shutdown | 1 |
| Error Handling | 2 |

## Adding Questions

Create a YAML file in `questions/` or `challenges/`:

```yaml
# questions/my-topic.yaml
chapter: my-topic
title: "My Topic"
source_url: "https://example.com"
questions:
  - id: "mytopic-1"
    type: multiple_choice
    difficulty: medium
    text: "Your question here?"
    correct: "The right answer"
    wrong:
      - "Plausible wrong answer 1"
      - "Plausible wrong answer 2"
      - "Plausible wrong answer 3"
      - "Plausible wrong answer 4"
      - "Plausible wrong answer 5"
    explanation: "Why the correct answer is right."
```

Restart the server to pick up new files.

## Tech Stack

- **Go** stdlib `net/http` with Go 1.22+ routing patterns
- **HTMX** for partial page updates
- **PicoCSS** for minimal styling
- **Monaco Editor** (CDN) for code challenges
- **gopkg.in/yaml.v3** for question loading
