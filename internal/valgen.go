package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
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

type valueFileValue struct {
	Name    string `yaml:"name"`
	Path    string `yaml:"path"`    //this is going to depend on the type - but let's assume it's grandparent.parent.child
	Default string `yaml:"default"` // The default if the value can't be found
}

type valueFile struct {
	Name   string           `yaml:"name"`
	Values []valueFileValue `yaml:"values"`
}

type surveyQuestionsFile struct {
	Questions  []surveyQuestion `yaml:"questions"`
	ValueFiles []valueFile      `yaml:"valueFiles"` // This is a list of files we can potentially source values from
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

func loadValuesFromValuesFiles(files []string) (map[string]interface{}, error) {
	vals := make(map[string]interface{})
	for _, file := range files {
		f, err := os.Stat(file)
		if err != nil {
			//this might be okay, just log it
			log.Default().Printf("Unable to find file `%v`", file)
		}
		extension := filepath.Ext(f.Name())
		switch extension {
		case ".env":
			myEnv, err := godotenv.Read(file)
			if err != nil {
				log.Default().Printf("Unable to read env file `%v`", file)
				continue
			}
			vals[f.Name()] = myEnv
		case ".json":
			fileData, err := os.ReadFile(file)
			if err != nil {
				log.Default().Printf("Unable to read json file `%v`", file)
				continue
			}
			var v interface{}
			err = json.Unmarshal(fileData, &v)
			if err != nil {
				log.Default().Printf("Unable to unmarshal json file `%v`", file)
				continue
			}
			vals[f.Name()] = v
		case ".yml":
			fallthrough
		case ".yaml":
			fileData, err := os.ReadFile(file)
			if err != nil {
				log.Default().Printf("Unable to read yaml file `%v`", file)
				continue
			}
			var v interface{}
			err = yaml.Unmarshal(fileData, &v)
			if err != nil {
				log.Default().Printf("Unable to unmarshal json file `%v`", file)
				continue
			}
			vals[f.Name()] = v
		default:
			return nil, errors.New(fmt.Sprintf("Unsupported file comprehension for `%v`", file))
		}
	}
	return vals, nil
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
				survey.AskOne(textQuestion, &resp, survey.WithValidator(survey.Required))
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
				survey.AskOne(selectQuestion, &resp, survey.WithValidator(survey.Required))
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
				survey.AskOne(selectQuestion, &resp, survey.WithValidator(survey.Required))
			}

			subinteractive := false
			if resp == "yes" {
				subinteractive = true
			}

			subVals, err := RunFromSurveyQuestions(question.Questions, subinteractive)
			if err != nil {
				return nil, err
			}

			unwoundVals := subVals.(map[string]interface{})
			for k, v := range subVals.(map[string]interface{}) {
				unwoundVals[k] = v
			}
			unwoundVals["answer"] = subinteractive

			vals[question.Name] = unwoundVals

		default:
			return nil, errors.New(fmt.Sprintf("Unknown question type `%v` for question `%v`", question.Type, question.Name))
		}
	}
	return vals, nil
}
