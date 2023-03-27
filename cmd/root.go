package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/go-git/go-git/v5" // with go modules disabled
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

var targetDirectory string
var scaffold string

var scaffolds = map[string]string{
	"laravel": "https://github.com/bomoko/lagoon-laravel-dir.git",
}

var rootCmd = &cobra.Command{
	Use:   "scaffold",
	Short: "Lagoon scaffold will pull a new site and fill in the details",
	Long:  `Lagoon scaffold will pull a new site and fill in the details`,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Lagoon scaffold will pull a new site and fill in the details",
	Long:  `Lagoon scaffold will pull a new site and fill in the details`,
	Run: func(cmd *cobra.Command, args []string) {
		if scaffold == "" {
			fmt.Println("Please select a scaffold")
			os.Exit(1)
		}

		repo, ok := scaffolds[scaffold]
		// If the key exists
		if !ok {
			fmt.Printf("Scaffold `%v` does not exist\n", scaffold)
		}

		//We'll use this when we want to use templates
		//let's checkout the scaffold into a tmp dir
		tDir, err := ioutil.TempDir(targetDirectory, "prefix")
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		//defer cleanRemoveDir(tDir)

		fmt.Println(tDir)

		_, err = git.PlainClone(tDir, false, &git.CloneOptions{
			URL:      repo,
			Progress: os.Stdout,
		})

		if err != nil {
			fmt.Println(err)
			return
		}

		err = cleanRemoveDir(tDir + "/.git")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err = processTemplates(nil, tDir); err != nil {
			fmt.Println(err)
			return
		}
		// For now we're just testing the dir traversal
		//err = cp.Copy(tDir, targetDirectory)
		//if err != nil {
		//	fmt.Println(err)
		//	os.Exit(1)
		//}
	},
}

func processTemplates(values interface{}, tempDir string) error {

	//we should find a values file in the root
	valfilename := tempDir + "/values.yml"
	if _, err := os.Stat(valfilename); errors.Is(err, os.ErrNotExist) {
		return errors.New(valfilename + " does not exist")
	}

	valuesDefaults, err := os.ReadFile(valfilename)
	if err != nil {
		return err
	}

	//let's open and edit the values file - this can move into proper survey questions in the future

	prompt := &survey.Editor{
		Renderer:      survey.Renderer{},
		Message:       "Shell code snippet",
		Default:       string(valuesDefaults),
		Help:          "",
		Editor:        "",
		HideDefault:   true,
		AppendDefault: true,
		FileName:      "*.yml",
	}
	var content string
	survey.AskOne(prompt, &content)

	var parsedContent interface{}
	//yam.
	err = yaml.Unmarshal([]byte(content), &parsedContent)
	if err != nil {
		return err
	}

	err = filepath.WalkDir(tempDir, func(p string, d fs.DirEntry, err error) error {
		if !d.IsDir() && filepath.Ext(p) == ".tmpl" {
			templ, err := template.ParseFiles(p)
			if err != nil {
				return err
			}
			var buf bytes.Buffer
			err = templ.Execute(&buf, parsedContent)
			if err != nil {
				return err
			}
			extension := path.Ext(p)
			outputName := p[:len(p)-len(extension)]
			//TODO: better permission handling?
			err = ioutil.WriteFile(outputName, buf.Bytes(), 0644)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
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
		fmt.Println("We currently support the following:")
		for pagage, _ := range scaffolds {
			fmt.Println(pagage)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(initCmd)
	initCmd.PersistentFlags().StringVar(&scaffold, "scaffold", "", "Which scaffold to pull into directory")
	initCmd.Flags().StringVar(&targetDirectory, "targetdir", "./", "Directory to check out project into - defaults to current directory")
}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
