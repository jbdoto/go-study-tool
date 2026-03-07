package challenges

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func LoadChallengeSets(dir string) ([]ChallengeSet, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return nil, fmt.Errorf("glob challenge files: %w", err)
	}

	var sets []ChallengeSet
	for _, f := range files {
		cs, err := loadChallengeSet(f)
		if err != nil {
			return nil, fmt.Errorf("load %s: %w", f, err)
		}
		sets = append(sets, cs)
	}
	return sets, nil
}

func loadChallengeSet(path string) (ChallengeSet, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ChallengeSet{}, err
	}
	var cs ChallengeSet
	if err := yaml.Unmarshal(data, &cs); err != nil {
		return ChallengeSet{}, err
	}
	return cs, nil
}
