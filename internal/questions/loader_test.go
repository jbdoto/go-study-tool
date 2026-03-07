package questions_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jbdoto/go-study-tool/internal/questions"
)

func TestLoadChapters(t *testing.T) {
	dir := t.TempDir()
	yaml := `chapter: test-chapter
title: "Test Chapter"
source_url: "https://example.com"
questions:
  - id: "tc-1"
    type: multiple_choice
    difficulty: easy
    text: "What is 1+1?"
    correct: "2"
    wrong:
      - "1"
      - "3"
      - "4"
      - "0"
      - "5"
    explanation: "Basic arithmetic."
  - id: "tc-2"
    type: multiple_choice
    difficulty: medium
    text: "Which keyword declares a variable?"
    correct: "var"
    wrong:
      - "func"
      - "type"
      - "import"
      - "const"
      - "package"
    explanation: "var is used to declare variables."
`
	err := os.WriteFile(filepath.Join(dir, "test-chapter.yaml"), []byte(yaml), 0644)
	if err != nil {
		t.Fatal(err)
	}

	chapters, err := questions.LoadChapters(dir)
	if err != nil {
		t.Fatal(err)
	}

	if len(chapters) != 1 {
		t.Fatalf("expected 1 chapter, got %d", len(chapters))
	}

	ch := chapters[0]
	if ch.Slug != "test-chapter" {
		t.Errorf("expected slug 'test-chapter', got %q", ch.Slug)
	}
	if ch.Title != "Test Chapter" {
		t.Errorf("expected title 'Test Chapter', got %q", ch.Title)
	}
	if len(ch.Questions) != 2 {
		t.Fatalf("expected 2 questions, got %d", len(ch.Questions))
	}

	q := ch.Questions[0]
	if q.ID != "tc-1" {
		t.Errorf("expected id 'tc-1', got %q", q.ID)
	}
	if len(q.Wrong) != 5 {
		t.Errorf("expected 5 wrong choices, got %d", len(q.Wrong))
	}
	if q.Correct != "2" {
		t.Errorf("expected correct='2', got %q", q.Correct)
	}
}

func TestLoadChaptersEmptyDir(t *testing.T) {
	dir := t.TempDir()
	chapters, err := questions.LoadChapters(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(chapters) != 0 {
		t.Errorf("expected 0 chapters, got %d", len(chapters))
	}
}