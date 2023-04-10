package cmd

import (
	"bomoko/lagoon-init/internal"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
)

var flowFile string

var flowCmd = &cobra.Command{
	Use:   "flow",
	Short: "Utilities for visualizing flow details",
	Long:  `Utilities for visualizing flow details`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if flowFile == "" {
			return errors.New("Please provide a file to visualize")
		} else {
			flowData, err := ioutil.ReadFile(flowFile)
			if err != nil {
				return fmt.Errorf("Error reading file: ", err)
			}
			data, _ := internal.UnmarshallSurveyQuestions(flowData)
			output, _ := internal.FlowToGraph(0, data)
			fmt.Printf("\n%s:\n\n", flowFile)
			fmt.Println(output)
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(flowCmd)
	flowCmd.Flags().StringVar(&flowFile, "file", "", "The flow file we'd like to visualize")
}
