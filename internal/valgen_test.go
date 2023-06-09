package internal

import (
	"reflect"
	"testing"
)

func Test_unmarshallSurveyQuestions(t *testing.T) {
	type args struct {
		incoming []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []surveyQuestion
		wantErr bool
	}{
		{
			name: "Test UnmarshallSurveyQuestions",
			args: args{
				incoming: []byte(`
questions:
- name: project_name
  help: The name of the project
  type: text
  required: true
  prompt: What is your project's name
  default: MyProject
  type: text`),
			},
			want: []surveyQuestion{
				{
					Name:     "project_name",
					Help:     "The name of the project",
					Type:     "text",
					Required: true,
					Prompt:   "What is your project's name",
					Default:  "MyProject",
				},
			},
		},
		{
			name: "Test UnmarshallSurveyQuestions conditional",
			args: args{
				incoming: []byte(`
questions:
- name: select_list
  help: select one of these options
  options:
    - option1
    - option2
    - option3
  type: select
  required: true
  prompt: Select one of these options
  default: option1
- name: a_conditional
  help: This is a conditional question
  type: conditional
  questions:
  - name: conditional_question
    type: text
    required: true
    prompt: This is a sub question
    default: default value
`),
			},
			want: []surveyQuestion{
				{
					Name:     "select_list",
					Help:     "select one of these options",
					Type:     "select",
					Required: true,
					Prompt:   "Select one of these options",
					Default:  "option1",
					Options:  []string{"option1", "option2", "option3"},
				},
				{
					Name: "a_conditional",
					Help: "This is a conditional question",
					Type: "conditional",
					Questions: []surveyQuestion{
						{
							Name:     "conditional_question",
							Type:     "text",
							Required: true,
							Prompt:   "This is a sub question",
							Default:  "default value",
						},
					},
				},
			},
		},
		{
			name: "Test UnmarshallSurveyQuestions select",
			args: args{
				incoming: []byte(`
questions:
- name: select_list
  help: select one of these options
  options:
    - option1
    - option2
    - option3
  type: select
  required: true
  prompt: Select one of these options
  default: option1
`),
			},
			want: []surveyQuestion{
				{
					Name:     "select_list",
					Help:     "select one of these options",
					Type:     "select",
					Required: true,
					Prompt:   "Select one of these options",
					Default:  "option1",
					Options:  []string{"option1", "option2", "option3"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshallSurveyQuestions(tt.args.incoming)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshallSurveyQuestions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnmarshallSurveyQuestions() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_runFromSurveyQuestions(t *testing.T) {
	type args struct {
		questions   []surveyQuestion
		interactive bool
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{{
		name: "Test Single Question",
		args: args{
			questions: []surveyQuestion{
				{
					Name:     "project_name",
					Help:     "The name of the project",
					Type:     "text",
					Required: true,
					Prompt:   "What is your project's name",
					Default:  "MyProject",
				},
			},
			interactive: false,
		},
		want: map[string]interface{}{
			"project_name": "MyProject",
		},
	},
		{
			name: "Test Conditional with sub questions",
			args: args{
				questions: []surveyQuestion{
					{
						Name:     "conditional",
						Type:     "conditional",
						Required: true,
						Prompt:   "Yes or no",
						Questions: []surveyQuestion{
							{
								Name:     "conditional_question",
								Type:     "text",
								Required: true,
								Prompt:   "This is a sub question",
								Default:  "value",
							},
						},
					},
				},
				interactive: false,
			},
			want: map[string]interface{}{
				"conditional":                      false,
				"conditional.conditional_question": "value",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RunFromSurveyQuestions(tt.args.questions, tt.args.interactive)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunFromSurveyQuestions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RunFromSurveyQuestions() got = %v, want %v", got, tt.want)
			}
		})
	}
}
