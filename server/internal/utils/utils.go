package domain

import (
	"regexp"
	"strings"
)

var slugRegex = regexp.MustCompile("[^a-z0-9]+")

func GenerateSlug(s string) string {
	s = strings.ToLower(s)
	s = slugRegex.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}
