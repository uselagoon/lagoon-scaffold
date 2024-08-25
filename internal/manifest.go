package internal

import (
	_ "embed"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
)

const manifestUrl = "https://raw.githubusercontent.com/uselagoon/lagoon-scaffold/main/internal/assets/scaffolds.yml"

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

func GetScaffolds(localmanifest string) (map[string]ScaffoldRepo, error) {

	if localmanifest != "" { // we try load up a manifest given locally
		// try open and read the file

		loader := &ScaffoldLoader{
			Scaffolds: make([]ScaffoldRepo, 0),
		}
		dat, err := os.ReadFile(localmanifest)
		err = yaml.Unmarshal(dat, loader)
		if err != nil {
			return map[string]ScaffoldRepo{}, err
		}

		return remapScaffoldLoader(loader), nil

	}

	return resolveScaffolds(manifestUrl), nil
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
