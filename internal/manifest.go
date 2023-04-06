package internal

import (
	_ "embed"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
)

//const manifestUrl = "https://raw.githubusercontent.com/uselagoon/lagoon-scaffold/main/internal/assets/scaffolds.yml"
const manifestUrl = "https://gist.githubusercontent.com/bomoko/161bbeabc6d17d69e7d52f233cce749c/raw/09c800e9b3460d11bdf2716977ef9f7f7a3bf8f6/scaffolds.yml"

//go:embed assets/scaffolds.yml
var defaultScaffolds []byte

func getManifestFromUrl(manifestUrl string) ([]byte, error) {
	resp, err := http.Get(manifestUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func getDefaultScaffold() (map[string]ScaffoldRepo, error) {
	loader := &ScaffoldLoader{
		Scaffolds: make([]ScaffoldRepo, 0),
	}
	err := yaml.Unmarshal(defaultScaffolds, loader)

	if err != nil {
		return nil, err
	}

	return remapScaffoldLoader(loader), nil
}

func resolveScaffolds(manifestUrl string) map[string]ScaffoldRepo {

	defaultScaffold, err := getDefaultScaffold()
	if err != nil {
		panic(err)
	}

	if manifestUrl == "" {
		return defaultScaffold
	}

	scaffolds, err := getManifestFromUrl(manifestUrl)
	if err != nil {
		return defaultScaffold
	}

	loader := &ScaffoldLoader{
		Scaffolds: make([]ScaffoldRepo, 0),
	}
	err = yaml.Unmarshal(scaffolds, loader)

	if err != nil {
		return defaultScaffold
	}

	return remapScaffoldLoader(loader)
}

func remapScaffoldLoader(loader *ScaffoldLoader) map[string]ScaffoldRepo {
	remapped := make(map[string]ScaffoldRepo)
	for _, scaffold := range loader.Scaffolds {
		remapped[scaffold.Name] = scaffold
	}
	return remapped
}

func GetScaffolds() map[string]ScaffoldRepo {
	return resolveScaffolds(manifestUrl)
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
