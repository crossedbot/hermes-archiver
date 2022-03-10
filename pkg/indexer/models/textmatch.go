package models

import (
	"fmt"
	"strings"
)

type TextMatch int

const (
	TextMatchExact TextMatch = iota
	TextMatchPartial
)

var TextMatchStrings = []string{
	"exact",
	"partial",
}

func (tm TextMatch) String() (match string) {
	if int(tm) > -1 && len(TextMatchStrings) > int(tm) {
		match = TextMatchStrings[tm]
	}
	return
}

func ToTextMatch(m string) (TextMatch, error) {
	for i, tm := range TextMatchStrings {
		if strings.EqualFold(tm, m) {
			return TextMatch(i), nil
		}
	}
	return TextMatch(-1), fmt.Errorf("unkown text matching string")
}
