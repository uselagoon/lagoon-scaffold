package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/uselagoon/lagoon-scaffold/internal"
)

var flowFile string

var flowCmd = &cobra.Command{
	Use:   "flow",
	Short: "Utilities for visualizing flow details",
	Long:  `Utilities for visualizing flow details`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if flowFile == "" {
			return errors.New("please provide a file to visualize")
		}
		flowData, err := os.ReadFile(flowFile)
		if err != nil {
			return fmt.Errorf("error reading file: %w", err)
		}
		data, err := internal.UnmarshallSurveyQuestions(flowData)
		if err != nil {
			return fmt.Errorf("error parsing flow file: %w", err)
		}
		output, err := internal.FlowToGraph(0, data)
		if err != nil {
			return err
		}
		fmt.Printf("\n%s:\n\n", flowFile)
		fmt.Println(output)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(flowCmd)
	flowCmd.Flags().StringVar(&flowFile, "file", "", "The flow file we'd like to visualize")
}
