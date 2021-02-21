package eqspec

import (
	"sort"
	"strings"
)

func IsItem(text string) bool {
	upperBeg:=strings.ToUpper(text)
	start:=sort.Search(len(everquestItems), func(i int) bool {
		upperItemName := strings.ToUpper(everquestItems[i])
		return upperBeg<=upperItemName
	})
	return start!=len(everquestItems) && strings.EqualFold(everquestItems[start], text)
}

func SuggestCompletions(beginningText string) []string {
	upperBeg:=strings.ToUpper(beginningText)
	start:=sort.Search(len(everquestItems), func(i int) bool {
		upperItemName := strings.ToUpper(everquestItems[i])
		return upperBeg<=upperItemName
	})
	end:=sort.Search(len(everquestItems), func(i int) bool {
		upperItemName := strings.ToUpper(everquestItems[i])
		return upperBeg+"\uffff"<=upperItemName
	})
	return everquestItems[start:end]
}
