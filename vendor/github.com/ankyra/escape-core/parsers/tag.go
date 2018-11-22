package parsers

import (
	"strings"
	"unicode"
)

var ForbiddenTags = []string{
	"latest",
	"@",
	"v@",
}

func IsValidTag(t string) bool {
	if t == "" {
		return false
	}
	for _, forbidden := range ForbiddenTags {
		if forbidden == t {
			return false
		}
	}
	if strings.HasPrefix(t, "v") {
		if maybeParseVersionQuery(t[1:]) != nil {
			return false
		}
	} else if unicode.IsDigit(rune(t[0])) {
		if maybeParseVersionQuery(t) != nil {
			return false
		}
	}
	return true
}
