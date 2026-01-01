package cmdstr

import (
	"regexp"
	"strings"
)

// Split separates the provided text at spaces
// while preserving quoted text as a single unit.
//
// Setting `includeQuotes` controls whether outer quotes
// are included in the output. Inner quotes are always included.
func Split(text string, includeQuotes bool) []string {
	var parts []string

	var builder strings.Builder
	var quoting string

	quotesRe := regexp.MustCompile("\"|'|`")
	spaceRe := regexp.MustCompile(`\s`)

	for _, r := range text {
		char := string(r)

		if quotesRe.MatchString(char) {
			switch quoting {
			case "": // Not already quoting, start.
				quoting = char

				if !includeQuotes {
					continue
				}
			case char: // Matches current quoting, end.
				quoting = ""

				if !includeQuotes {
					continue
				}
			}
			// Inner quotes are always included
		} else if quoting == "" && spaceRe.MatchString(char) {
			parts = append(parts, builder.String())
			builder.Reset()
			continue
		}

		builder.WriteRune(r)
	}

	if builder.Len() > 0 {
		parts = append(parts, builder.String())
	}

	return parts
}
