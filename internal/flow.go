package internal

import (
	"fmt"
	"github.com/fatih/color"
)

func printColor(depth int, s string) string {
	switch depth % 4 {
	case 0:
		return color.RedString(s)
	case 1:
		return color.GreenString(s)
	case 2:
		return color.YellowString(s)
	case 3:
		return color.RedString(s)
	}
	return s
}

func repeatColorWithDepth(s string, depth int) string {
	ret := ""
	for i := 0; i < depth; i++ {
		ret = ret + printColor(i, s)
	}
	return ret
}

func FlowToGraph(depth int, questions []surveyQuestion) (string, error) {
	graph := ""
	if depth > 0 {
		graph += fmt.Sprintf("%s%s\n", repeatColorWithDepth("|  ", depth), printColor(depth, "\\"))
	}
	for _, question := range questions {
		questionFormatted := printColor(depth, fmt.Sprintf("| %s:%s", question.Name, question.Prompt))
		graph += fmt.Sprintf("%s%s\n", repeatColorWithDepth("|  ", depth), questionFormatted)
		if question.Type == "conditional" {
			conditionalGraph, err := FlowToGraph(depth+1, question.Questions)
			if err != nil {
				return "", err
			}
			graph += conditionalGraph
		}
	}
	if depth > 0 {
		graph += fmt.Sprintf("%s%s\n", repeatColorWithDepth("|  ", depth), printColor(depth, "/"))
	}
	return graph, nil
}
