# Go Study Tool

Web-based drill app for Go interview prep, built with Go stdlib `net/http`, HTMX, and PicoCSS.

## Project Structure

- `cmd/server/` — entry point, embeds templates via `//go:embed all:templates`
- `cmd/server/templates/` — HTML templates (layout, index, quiz, result, done)
- `internal/server/` — HTTP handlers and routing
- `internal/questions/` — YAML loader and types
- `questions/` — YAML question files, one per chapter

## Question YAML Format

Questions use `correct` (string, the right answer) and `wrong` (list of 5+ incorrect options). The handler picks 3 random wrong answers + the correct one and shuffles them each time.

```yaml
- id: "chapter-mc-1"
  type: multiple_choice
  difficulty: easy|medium|hard
  text: "Question text"
  correct: "The correct answer"
  wrong:
    - "Wrong answer 1"
    - "Wrong answer 2"
    - "Wrong answer 3"
    - "Wrong answer 4"
    - "Wrong answer 5"
  explanation: "Why the correct answer is right."
  source_anchor: "section-heading-slug"
```

## GitBook Anchor Slugs

When adding `source_anchor` values for questions, use the actual anchor from the GitBook page. GitBook preserves dots/ellipses in heading slugs rather than stripping them:

- "Hello, world... again" → `hello-world...-again` (NOT `hello-world-again`)
- "one...last...refactor?" → `one...last...refactor` (NOT `onelastrefactor`)

Always verify anchors by checking the actual page headings in browser rather than guessing the slug format.

## Running

```sh
cd ~/Desktop/projects/jbdoto/go-study-tool
go run ./cmd/server
# Open http://localhost:8080
```

## Adding New Chapters

1. Fetch chapter content with WebFetch to identify key concepts and section anchors
2. Create `questions/{chapter-slug}.yaml` with 8-12 questions
3. Include 5+ wrong answers per question for variety
4. Set `source_anchor` to the nearest heading slug from the gitbook page
5. Aim for difficulty spread: ~3 easy, ~5 medium, ~3 hard