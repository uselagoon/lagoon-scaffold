package internal

import (
	"regexp"
	"text/template"
)

// templating.go contains, primarily, extra functions for templating.

var TemplatingExtensions = template.FuncMap{
	"regexMatch": func(pattern, str string) bool {
		matched, _ := regexp.MatchString(pattern, str)
		return matched
	},
}

func GetTemplate() *template.Template {
	return template.New("LagoonsyncTemplate").Funcs(TemplatingExtensions)
}
