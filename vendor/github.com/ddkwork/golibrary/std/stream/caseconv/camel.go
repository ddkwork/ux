package caseconv

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func ToCamel(str string) string {
	chunks := chunk(str)
	for i, c := range chunks {
		if i == 0 {
			chunks[i] = strings.ToLower(c)
			continue
		}
		chunks[i] = cases.Title(language.AmericanEnglish).String(c) // 分词并把首字母大写
	}
	return strings.Join(chunks, "")
}
