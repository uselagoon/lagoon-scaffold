package internal

import (
	_ "embed"
	"gopkg.in/yaml.v2"
)

const manifestUrl = "https://raw.githubusercontent.com/bomoko/lagoon-init/main/assets/scaffolds.yml"

//go:embed assets/scaffolds.yml
var defaultScaffolds []byte

func resolveScaffolds() (map[string]ScaffoldRepo, error) {
	loader := &ScaffoldLoader{
		Scaffolds: make([]ScaffoldRepo, 0),
	}
	err := yaml.Unmarshal(defaultScaffolds, loader)

	if err != nil {
		return nil, err
	}

	remapped := make(map[string]ScaffoldRepo)
	for _, scaffold := range loader.Scaffolds {
		remapped[scaffold.Name] = scaffold
	}

	return remapped, nil
}

func GetScaffolds() map[string]ScaffoldRepo {
	return scaffolds
}

var scaffolds = map[string]ScaffoldRepo{
	"laravel-init": {
		GitRepo:          "https://github.com/bomoko/lagoon-laravel-dir.git",
		Branch:           "main",
		ShortDescription: "Will add a minimal set of files to an existing Laravel 10 installation",
		Description:      "Will add a minimal set of files to an existing Laravel 10 installation",
	},
	"drupal-9": {
		GitRepo:          "https://github.com/lagoon-examples/drupal9-full.git",
		Branch:           "scaffold",
		ShortDescription: "Pulls and sets up a new Lagoon ready Drupal 9",
		Description:      "Pulls and sets up a new Lagoon ready Drupal 9",
	},
}

type ScaffoldRepo struct {
	Name             string `yaml:"name,omitempty"`
	GitRepo          string `yaml:"git_repo,omitempty"`
	Branch           string `yaml:"branch,omitempty"`
	Description      string `yaml:"description,omitempty"`
	ShortDescription string `yaml:"shortDescription,omitempty"`
}

type ScaffoldLoader struct {
	Scaffolds []ScaffoldRepo `yaml:"scaffolds,omitempty"`
}
