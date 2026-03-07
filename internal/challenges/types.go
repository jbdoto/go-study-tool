package challenges

type ChallengeSet struct {
	Slug       string          `yaml:"challenge_set"`
	Title      string          `yaml:"title"`
	Challenges []CodeChallenge `yaml:"challenges"`
}

type CodeChallenge struct {
	ID          string   `yaml:"id"`
	Type        string   `yaml:"type"`
	Difficulty  string   `yaml:"difficulty"`
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Scaffold    string   `yaml:"scaffold"`
	Solution    string   `yaml:"solution"`
	Hints       []string `yaml:"hints,omitempty"`
	KeyConcepts []string `yaml:"key_concepts"`
}
