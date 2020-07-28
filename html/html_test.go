package html

import (
	"regexp"
	"strings"
	"unicode"
)

var tagSpace0 = regexp.MustCompile(`\s+\\*>`)

func removeSpace(s string) string {
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", "")
	s = tagSpace0.ReplaceAllString(s, ">")
	s = reduceSpace(s)
	return strings.TrimSpace(s)
}

// reduceSpace takes all space not with " " and reduces the space count to 1.
func reduceSpace(s string) string {
	rs := make([]rune, 0, len(s))
	var inQuote bool
	var lastWasSpace bool
	for i := 0; i < len(s); i++ {
		r := rune(s[i])
		if r == '"' {
			if inQuote {
				inQuote = false
			} else {
				inQuote = true
			}
			lastWasSpace = false
			continue
		}
		if inQuote {
			rs = append(rs, r)
			continue
		}
		if unicode.IsSpace(r) {
			lastWasSpace = true
			continue
		}
		if lastWasSpace {
			lastWasSpace = false
			rs = append(rs, rune(' '))
		}
		rs = append(rs, r)
	}
	return string(rs)
}
