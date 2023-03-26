package cmd

import (
	"fmt"
	"github.com/go-git/go-git/v5" // with go modules disabled
	cp "github.com/otiai10/copy"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
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
		tDir, err := ioutil.TempDir("./", "prefix")
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		defer cleanRemoveDir(tDir)

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

		err = cp.Copy(tDir, targetDirectory)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
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
