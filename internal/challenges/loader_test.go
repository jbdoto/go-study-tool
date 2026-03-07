package challenges_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jbdoto/go-study-tool/internal/challenges"
)

func TestLoadChallengeSets(t *testing.T) {
	dir := t.TempDir()
	yaml := `challenge_set: test-set
title: "Test Challenge Set"
challenges:
  - id: "ts-1"
    type: code
    difficulty: easy
    title: "Hello World"
    description: "Write a function that returns hello world."
    scaffold: |
      package main

      func Hello() string {
          // your code here
      }
    solution: |
      package main

      func Hello() string {
          return "hello world"
      }
    hints:
      - "Use the return keyword"
      - "Return a string literal"
    key_concepts:
      - functions
      - strings
  - id: "ts-2"
    type: code
    difficulty: medium
    title: "Add Two Numbers"
    description: "Write a function that adds two integers."
    scaffold: |
      package main

      func Add(a, b int) int {
          // your code here
      }
    solution: |
      package main

      func Add(a, b int) int {
          return a + b
      }
    hints:
      - "Use the + operator"
    key_concepts:
      - functions
      - integers
`
	err := os.WriteFile(filepath.Join(dir, "test-set.yaml"), []byte(yaml), 0644)
	if err != nil {
		t.Fatal(err)
	}

	sets, err := challenges.LoadChallengeSets(dir)
	if err != nil {
		t.Fatal(err)
	}

	if len(sets) != 1 {
		t.Fatalf("expected 1 challenge set, got %d", len(sets))
	}

	cs := sets[0]
	if cs.Slug != "test-set" {
		t.Errorf("expected slug 'test-set', got %q", cs.Slug)
	}
	if cs.Title != "Test Challenge Set" {
		t.Errorf("expected title 'Test Challenge Set', got %q", cs.Title)
	}
	if len(cs.Challenges) != 2 {
		t.Fatalf("expected 2 challenges, got %d", len(cs.Challenges))
	}

	c := cs.Challenges[0]
	if c.ID != "ts-1" {
		t.Errorf("expected id 'ts-1', got %q", c.ID)
	}
	if len(c.Hints) != 2 {
		t.Errorf("expected 2 hints, got %d", len(c.Hints))
	}
	if len(c.KeyConcepts) != 2 {
		t.Errorf("expected 2 key concepts, got %d", len(c.KeyConcepts))
	}
}

func TestLoadChallengeSetsEmptyDir(t *testing.T) {
	dir := t.TempDir()
	sets, err := challenges.LoadChallengeSets(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(sets) != 0 {
		t.Errorf("expected 0 challenge sets, got %d", len(sets))
	}
}
