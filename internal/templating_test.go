package internal

import (
	"bytes"
	"testing"
)

func TestRegexMatch(t *testing.T) {
	tests := []struct {
		pattern string
		input   string
		expect  bool
	}{
		{"^hello", "hello world", true},
		{"^world", "hello world", false},
		{"world$", "hello world", true},
		{"[0-9]+", "abc123", true},
		{"^[a-z]{3}$", "abc", true},
		{"^[a-z]{3}$", "abcd", false},
		{"[A-Z]+", "lowercase", false},
		{"[A-Z]+", "UPPERCASE", true},
	}

	for _, test := range tests {
		result := TemplatingExtensions["regexMatch"].(func(string, string) bool)(test.pattern, test.input)
		if result != test.expect {
			t.Errorf("regexMatch(%q, %q) = %v; want %v", test.pattern, test.input, result, test.expect)
		}
	}
}

func TestGetTemplate(t *testing.T) {
	tests := []struct {
		name     string
		template string
		input    string
		output   string
	}{
		{"templateWithRegexMatch", "{{if regexMatch \"^hello\" .}}Matched!{{else}}No Match{{end}}", "hello world", "Matched!"},
		{"templateWithRegexMatch", "{{if regexMatch \"^hello\" .}}Matched!{{else}}No Match{{end}}", "goodbye world", "No Match"},
	}

	for _, test := range tests {
		ct := GetTemplate("")
		pt, err := ct.Parse(test.template)
		var buf bytes.Buffer
		err = pt.Execute(&buf, test.input)
		if err != nil {
			return
		}
		result := buf.String()

		//result := TemplatingExtensions["regexMatch"].(func(string, string) bool)(test.pattern, test.input)
		if result != test.output {
			t.Errorf("parsed template %v, with %v :- Got %v; want %v", test.template, test.input, result, test.output)
		}
	}
}
