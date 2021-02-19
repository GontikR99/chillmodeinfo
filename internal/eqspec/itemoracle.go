package eqspec

import (
	"sort"
	"strings"
)

func SuggestCompletions(beginningText string) []string {
	upperBeg:=strings.ToUpper(beginningText)
	start:=sort.Search(len(everquestItems), func(i int) bool {
		upperItemName := strings.ToUpper(everquestItems[i])
		return upperBeg<=upperItemName
	})
	var suggestions []string
	for idx:=start;idx<len(everquestItems)&&idx<start+10;idx++ {
		upperItemName := strings.ToUpper(everquestItems[idx])
		if strings.HasPrefix(upperItemName, upperBeg) {
			suggestions=append(suggestions, everquestItems[idx])
		}
	}
	return suggestions
}
