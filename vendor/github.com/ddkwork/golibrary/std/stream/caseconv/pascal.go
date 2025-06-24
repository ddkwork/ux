package caseconv

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func ToPascal(str string) string {
	chunks := chunk(str)
	for i, c := range chunks {
		chunks[i] = cases.Title(language.English).String(c)
	}
	return strings.Join(chunks, "")
}
