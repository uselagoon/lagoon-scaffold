package internal

import (
	"errors"
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
			}
			resp := ""
			if interactive {
				survey.AskOne(textQuestion, &resp, survey.WithValidator(survey.Required))
			}
			vals[question.Name] = question.Default
			if resp != "" {
				vals[question.Name] = resp
			}
		case "select":
			selectQuestion := &survey.Select{
				Message: question.Prompt, Options: question.Options, Default: question.Default,
			}
			resp := ""
			if interactive {
				survey.AskOne(selectQuestion, &resp, survey.WithValidator(survey.Required))
			}
			vals[question.Name] = selectQuestion.Default
			if resp != "" {
				vals[question.Name] = resp
			}
		case "conditional": //This isn't strictly a survey question type, but it's a useful way to group questions
			selectQuestion := &survey.Select{
				Message: question.Prompt, Options: []string{"yes", "no"}, Default: "no",
			}
			resp := ""
			if interactive {
				survey.AskOne(selectQuestion, &resp, survey.WithValidator(survey.Required))
			}

			subinteractive := false
			if resp == "yes" {
				subinteractive = true
			}

			vals[question.Name] = subinteractive

			subVals, err := RunFromSurveyQuestions(question.Questions, subinteractive)
			if err != nil {
				return nil, err
			}
			//var vals[question.Name] map[string]interface{}
			unwoundVals := make(map[string]interface{})
			for k, v := range subVals.(map[string]interface{}) {
				//vals[question.Name] =
				unwoundVals[k] = v
			}
			vals[question.Name] = unwoundVals

		default:
			return nil, errors.New(fmt.Sprintf("Unknown question type `%v` for question `%v`", question.Type, question.Name))
		}
	}
	return vals, nil
}
