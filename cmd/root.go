package cmd

import (
	"bomoko/lagoon-init/internal"
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"text/template"

	"github.com/AlecAivazis/survey/v2"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	cp "github.com/otiai10/copy"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var targetDirectory string
var localManifest string
var scaffold string
var noInteraction bool

func getScaffoldsKeys() []string {
	scaffolds, _ := internal.GetScaffolds(localManifest)
	var ret []string
	for k := range scaffolds {
		ret = append(ret, k)
	}
	sort.Strings(ret)
	return ret
}

func selectScaffold(scaffold *string) error {
	scaffolds, err := internal.GetScaffolds(localManifest)
	if err != nil {
		return err
	}
	prompt := survey.Select{
		Message: "Select a scaffold to run",
		Options: getScaffoldsKeys(),
		Description: func(value string, index int) string {
			return scaffolds[value].ShortDescription
		},
	}

	survey.AskOne(&prompt, scaffold)
	return nil
}

var RootCmd = &cobra.Command{
	Use:   "scaffold",
	Short: "Lagoon scaffold will pull a new site and fill in the details",
	Long:  `Lagoon scaffold will pull a new site and fill in the details`,
	RunE: func(cmd *cobra.Command, args []string) error {

		scaffolds, err := internal.GetScaffolds(localManifest)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if scaffold == "" && noInteraction {
			return errors.New("Please select a scaffold")
		}

		if scaffold == "" {
			selectScaffold(&scaffold)
		}

		repo, ok := scaffolds[scaffold]
		// If the key exists
		if !ok {
			return errors.New(fmt.Sprintf("Scaffold `%v` does not exist", scaffold))
		}

		//We'll use this when we want to use templates
		//let's checkout the scaffold into a tmp dir
		tDir, err := ioutil.TempDir(targetDirectory, "prefix")
		if err != nil {
			return err
		}
		defer cleanRemoveDir(tDir)

		fmt.Println(tDir)

		_, err = git.PlainClone(tDir, false, &git.CloneOptions{
			URL: repo.GitRepo,
			//Depth:         1,
			ReferenceName: plumbing.NewBranchReferenceName(repo.Branch),
			SingleBranch:  true,
			Progress:      os.Stdout,
		})

		if err != nil {
			return err
		}

		err = cleanRemoveDir(tDir + "/.git")
		if err != nil {
			return err
		}

		rawYaml, err := ioutil.ReadFile(tDir + "/.lagoon/flow.yml")
		if err != nil {
			return err
		}

		questions, err := internal.UnmarshallSurveyQuestions(rawYaml)

		if err != nil {
			return err
		}

		values, err := internal.RunFromSurveyQuestions(questions, !noInteraction)

		if err = processTemplates(values, tDir); err != nil {
			return err
		}

		// Let's now dump the output of the flow file into a values file
		valuesYml, err := yaml.Marshal(values)
		if err != nil {
			return err
		}
		if err := os.WriteFile(tDir+"/.lagoon/values.yml", valuesYml, 0644); err != nil {
			return err
		}

		showPostMessage(tDir)

		// For now we're just testing the dir traversal
		err = cp.Copy(tDir, targetDirectory)
		if err != nil {
			return err
		}

		return nil
	},
}

func processTemplates(values interface{}, tempDir string) error {
	return filepath.WalkDir(tempDir, func(p string, d fs.DirEntry, err error) error {
		if !d.IsDir() && filepath.Ext(p) == ".lgtmpl" {
			templ, err := template.ParseFiles(p)
			if err != nil {
				return err
			}
			var buf bytes.Buffer
			err = templ.Execute(&buf, values)
			if err != nil {
				return err
			}

			outputName := p[:len(p)-len(path.Ext(p))]
			err = os.WriteFile(outputName, buf.Bytes(), 0644)
			if err != nil {
				return err
			}
			//remove the file from the temp dir
			err = os.Remove(p)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func readValuesFile(tempDir string, noInteraction bool) (interface{}, error) {
	//we should find a values file in the root
	valfilename := tempDir + "/.lagoon/values.yml"
	if _, err := os.Stat(valfilename); errors.Is(err, os.ErrNotExist) {
		return nil, errors.New(valfilename + " does not exist")
	}

	valuesDefaults, err := os.ReadFile(valfilename)
	if err != nil {
		return nil, err
	}

	var content string

	if !noInteraction {
		//let's open and edit the values file - this can move into proper survey questions in the future
		prompt := &survey.Editor{
			Renderer:      survey.Renderer{},
			Message:       "We will now open your values file for editing",
			Default:       string(valuesDefaults),
			Help:          "",
			Editor:        "",
			HideDefault:   true,
			AppendDefault: true,
			FileName:      "*.yml",
		}
		survey.AskOne(prompt, &content)
	} else { //we simply use the defaults...
		content = string(valuesDefaults)
	}

	var parsedContent interface{}
	err = yaml.Unmarshal([]byte(content), &parsedContent)
	if err != nil {
		return nil, err
	}
	return parsedContent, err
}

func showPostMessage(tempDir string) {
	valfilename := tempDir + "/.lagoon/post-message.txt"
	if _, err := os.Stat(valfilename); errors.Is(err, os.ErrNotExist) {
		return //no post-message
	}

	text, err := ioutil.ReadFile(valfilename)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(string(text))
}

func cleanRemoveDir(dir string) error {
	return os.RemoveAll(dir)
}

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List currently supported templates",
	Long:    "Lists all currently supported Lagoon scaffolds",
	Example: "lagoon-init-prot list",
	Run: func(cmd *cobra.Command, args []string) {
		scaffolds, _ := internal.GetScaffolds(localManifest)
		fmt.Println("We currently support the following:")
		for pagage := range scaffolds {
			fmt.Println(pagage)
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
	RootCmd.PersistentFlags().StringVar(&scaffold, "scaffold", "", "Which scaffold to pull into directory")
	RootCmd.Flags().BoolVar(&noInteraction, "no-interaction", false, "Don't interactively fill in any values for the scaffold - use defaults")
	RootCmd.Flags().StringVar(&targetDirectory, "targetdir", "./", "Directory to check out project into - defaults to current directory")
	RootCmd.Flags().StringVar(&localManifest, "manifest", "", "Custom local manifest file for scaffold list - defaults to an empty string")
}

func Execute() {

	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
