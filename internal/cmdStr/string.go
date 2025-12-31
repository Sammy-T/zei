package cmdstr

import (
	"regexp"
	"strings"
)

// Split separates the provided text at spaces
// while preserving quoted text as a single unit.
func Split(text string) []string {
	var parts []string

	var builder strings.Builder
	var quoting string

	var cmdLen = len(text)

	quotesRe := regexp.MustCompile("\"|'|`")
	spaceRe := regexp.MustCompile(`\s`)

	for i, r := range text {
		char := string(r)

		if quotesRe.MatchString(char) {
			switch quoting {
			case "": // Not already quoting, start.
				quoting = char
			case char: // Matches current quoting, end.
				quoting = ""
			}
		} else if quoting == "" && spaceRe.MatchString(char) {
			parts = append(parts, builder.String())
			builder.Reset()
			continue
		}

		builder.WriteRune(r)

		if i == cmdLen-1 {
			parts = append(parts, builder.String())
		}
	}

	return parts
}
