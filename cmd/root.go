package cmd

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/AlecAivazis/survey/v2"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	cp "github.com/otiai10/copy"
	"github.com/spf13/cobra"
	"github.com/uselagoon/lagoon-scaffold/internal"
	"gopkg.in/yaml.v2"
)

var targetDirectory string
var localManifest string
var scaffold string
var noInteraction bool
var inputFile string
var privateKeyFile string

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

	return survey.AskOne(&prompt, scaffold)
}

var RootCmd = &cobra.Command{
	Use:   "lagoon-scaffold",
	Short: "Lagoon scaffold will pull a new site and fill in the details",
	Long:  `Lagoon scaffold will pull a new site and fill in the details`,
	RunE: func(cmd *cobra.Command, args []string) error {

		scaffolds, err := internal.GetScaffolds(localManifest)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if scaffold == "" && noInteraction {
			return fmt.Errorf("please select a scaffold")
		}

		if scaffold == "" {
			if err := selectScaffold(&scaffold); err != nil {
				return err
			}
		}

		repo, ok := scaffolds[scaffold]
		// If the key exists
		if !ok {
			return fmt.Errorf("scaffold `%v` does not exist", scaffold)
		}

		//let's checkout the scaffold into a tmp dir
		tDir, err := os.MkdirTemp(targetDirectory, ".lagoon-scaffold-")
		if err != nil {
			return err
		}
		defer cleanRemoveDir(tDir)

		// Here we deal with sshkeys, if one is passed to us

		cloneOptions := &git.CloneOptions{
			URL: repo.GitRepo,
			//Depth:         1,
			ReferenceName: plumbing.NewBranchReferenceName(repo.Branch),
			SingleBranch:  true,
			Progress:      os.Stdout,
		}

		if privateKeyFile != "" {
			_, err := os.Stat(privateKeyFile)
			if err != nil {
				log.Fatalf("read file %s failed %s\n", privateKeyFile, err.Error())
			}

			publicKeys, err := ssh.NewPublicKeysFromFile("git", privateKeyFile, "")
			if err != nil {
				log.Fatalf("generate publickeys failed: %s\n", err.Error())
			}

			cloneOptions.Auth = publicKeys
		}

		_, err = git.PlainClone(tDir, false, cloneOptions)

		if err != nil {
			return err
		}

		err = cleanRemoveDir(filepath.Join(tDir, ".git"))
		if err != nil {
			return err
		}

		rawYaml, err := os.ReadFile(filepath.Join(tDir, ".lagoon/flow.yml"))
		if err != nil {
			return err
		}

		questions, err := internal.UnmarshallSurveyQuestions(rawYaml)

		if err != nil {
			return err
		}

		var values interface{}

		if inputFile == "" {
			values, err = internal.RunFromSurveyQuestions(questions, !noInteraction)
			if err != nil {
				log.Fatalf("Error running survey: %v", err)
			}
		} else { // we're going to attempt to load these values from the file
			yamlFile, err := os.ReadFile(inputFile)
			if err != nil {
				log.Fatalf("Error reading YAML file: %v", err)
			}
			err = yaml.Unmarshal(yamlFile, &values)
			if err != nil {
				log.Fatalf("Error parsing YAML file: %v", err)
			}
		}

		if err = processTemplates(values, tDir); err != nil {
			return err
		}

		// Let's now dump the output of the flow file into a values file
		valuesYml, err := yaml.Marshal(values)
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(tDir, ".lagoon/values.yml"), valuesYml, 0644); err != nil {
			return err
		}

		showPostMessage(tDir)

		return cp.Copy(tDir, targetDirectory)
	},
}

func processTemplates(values interface{}, tempDir string) error {
	return filepath.WalkDir(tempDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(p) == ".lgtmpl" {
			templ, err := internal.GetTemplate(filepath.Base(p)).ParseFiles(p)
			if err != nil {
				return err
			}
			var buf bytes.Buffer
			err = templ.Execute(&buf, values)
			if err != nil {
				return err
			}

			outputName := p[:len(p)-len(filepath.Ext(p))]
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

func showPostMessage(tempDir string) {
	text, err := os.ReadFile(filepath.Join(tempDir, ".lagoon/post-message.txt"))
	if err != nil {
		return //no post-message
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
	Example: "lagoon-scaffold list",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("We currently support the following:")
		for _, name := range getScaffoldsKeys() {
			fmt.Println(name)
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
	RootCmd.PersistentFlags().StringVar(&scaffold, "scaffold", "", "Which scaffold to pull into directory")
	RootCmd.Flags().BoolVar(&noInteraction, "no-interaction", false, "Don't interactively fill in any values for the scaffold - use defaults")
	RootCmd.Flags().StringVar(&targetDirectory, "targetdir", "./", "Directory to check out project into - defaults to current directory")
	RootCmd.Flags().StringVar(&localManifest, "manifest", "", "Custom local manifest file for scaffold list - defaults to an empty string")
	RootCmd.Flags().StringVar(&inputFile, "values", "", "A Yaml file that provides defaults/answers for a scaffold - can be used in automation")
	RootCmd.Flags().StringVar(&privateKeyFile, "privatekey", "", "If private repository is used, this points to the private key used to access it")
}

func Execute() {

	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
