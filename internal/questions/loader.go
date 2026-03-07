package questions

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func LoadChapters(dir string) ([]Chapter, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return nil, fmt.Errorf("glob question files: %w", err)
	}

	var chapters []Chapter
	for _, f := range files {
		ch, err := loadChapter(f)
		if err != nil {
			return nil, fmt.Errorf("load %s: %w", f, err)
		}
		chapters = append(chapters, ch)
	}
	return chapters, nil
}

func loadChapter(path string) (Chapter, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Chapter{}, err
	}
	var ch Chapter
	if err := yaml.Unmarshal(data, &ch); err != nil {
		return Chapter{}, err
	}
	return ch, nil
}