package internal

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"gopkg.in/yaml.v2"
)

type surveyQuestion struct {
	Name      string           `yaml:"name"`
	Type      string           `yaml:"type"`
	Required  bool             `yaml:"required"`
	Help      string           `yaml:"help"`
	Prompt    string           `yaml:"prompt"`
	Default   string           `yaml:"default"`
	Options   []string         `yaml:"options"`
	Questions []surveyQuestion `yaml:"questions,omitempty"`
}

type surveyQuestionsFile struct {
	Questions []surveyQuestion `yaml:"questions"`
}

func UnmarshallSurveyQuestions(incoming []byte) ([]surveyQuestion, error) {
	//first we unmarshal the incoming interface into a map
	var incomingMap surveyQuestionsFile
	err := yaml.Unmarshal(incoming, &incomingMap)
	if err != nil {
		return nil, err
	}
	return incomingMap.Questions, nil
}

func RunFromSurveyQuestions(questions []surveyQuestion, interactive bool) (interface{}, error) {
	vals := make(map[string]interface{})
	for _, question := range questions {
		switch question.Type {
		case "text":
			textQuestion := &survey.Input{
				Message: question.Prompt,
				Default: question.Default,
				Help:    question.Help,
			}
			resp := ""
			if interactive {
				if err := survey.AskOne(textQuestion, &resp, survey.WithValidator(survey.Required)); err != nil {
					return nil, err
				}
			}
			vals[question.Name] = question.Default
			if resp != "" {
				vals[question.Name] = resp
			}
		case "select":
			selectQuestion := &survey.Select{
				Message: question.Prompt, Options: question.Options, Default: question.Default, Help: question.Help,
			}
			resp := ""
			if interactive {
				if err := survey.AskOne(selectQuestion, &resp, survey.WithValidator(survey.Required)); err != nil {
					return nil, err
				}
			}
			vals[question.Name] = question.Default
			if resp != "" {
				vals[question.Name] = resp
			}
		case "conditional": //This isn't strictly a survey question type, but it's a useful way to group questions
			selectQuestion := &survey.Select{
				Message: question.Prompt, Options: []string{"yes", "no"}, Default: "no", Help: question.Help,
			}
			resp := ""
			if interactive {
				if err := survey.AskOne(selectQuestion, &resp, survey.WithValidator(survey.Required)); err != nil {
					return nil, err
				}
			}

			subinteractive := resp == "yes"

			subVals, err := RunFromSurveyQuestions(question.Questions, subinteractive)
			if err != nil {
				return nil, err
			}

			unwoundVals := subVals.(map[string]interface{})
			unwoundVals["answer"] = subinteractive

			vals[question.Name] = unwoundVals

		default:
			return nil, fmt.Errorf("unknown question type `%v` for question `%v`", question.Type, question.Name)
		}
	}
	return vals, nil
}
