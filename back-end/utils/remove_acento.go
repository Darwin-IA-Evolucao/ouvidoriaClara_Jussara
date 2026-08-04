package utils

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

func RemoveAcento(str string) string {
	t := norm.NFD.String(str)

	var builder strings.Builder
	builder.Grow(len(t))

	for _, r := range t {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		builder.WriteRune(r)
	}
	return builder.String()
}
