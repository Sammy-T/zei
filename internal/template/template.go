package template

import (
	"regexp"
	"slices"
)

// ParseFields returns the field names
// found in the provided template text.
func ParseFields(text string) []string {
	tmplFieldRe := regexp.MustCompile(`{{\.?[\w\d\s]+}}`)
	tmplFieldNameRe := regexp.MustCompile(`[\w\d]+`)

	fields := tmplFieldRe.FindAllString(text, -1)
	var names []string

	for _, field := range fields {
		name := tmplFieldNameRe.FindString(field)

		if !slices.Contains(names, name) {
			names = append(names, name)
		}
	}

	return names
}
