package internal

import (
	"regexp"
	"strings"
)

var timeLayoutTagValueRe = regexp.MustCompile("^timeLayout=(.+)$")
var unixTimeUnitTagValueRe = regexp.MustCompile("^unixTimeUnit=(.+)$")

func ExtractTimeTag(tagValues []string) (string, string) {
	timeLayout := ""
	unixTimeUnit := ""

	for _, t := range tagValues {
		trimTag := strings.TrimSpace(t)

		submatch := timeLayoutTagValueRe.FindStringSubmatch(trimTag)
		if len(submatch) >= 2 {
			timeLayout = submatch[1]
			continue
		}

		submatch = unixTimeUnitTagValueRe.FindStringSubmatch(trimTag)
		if len(submatch) >= 2 {
			unixTimeUnit = submatch[1]
		}
	}

	return timeLayout, unixTimeUnit
}
