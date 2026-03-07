package questions

type Chapter struct {
	Slug      string     `yaml:"chapter"`
	Title     string     `yaml:"title"`
	SourceURL string     `yaml:"source_url"`
	Questions []Question `yaml:"questions"`
}

type Question struct {
	ID          string   `yaml:"id"`
	Type        string   `yaml:"type"`
	Difficulty  string   `yaml:"difficulty"`
	Text        string   `yaml:"text"`
	Correct     string   `yaml:"correct"`
	Wrong       []string `yaml:"wrong"`
	Explanation  string   `yaml:"explanation"`
	SourceAnchor string   `yaml:"source_anchor,omitempty"`
}