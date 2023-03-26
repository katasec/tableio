package reflectx

import (
	"regexp"
	"strings"
)

func TrimSuffix(s, suffix string) string {
	hasSuffix := strings.HasSuffix(s, suffix)
	if hasSuffix {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

func ToSnakeCase(str string) string {
	matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
