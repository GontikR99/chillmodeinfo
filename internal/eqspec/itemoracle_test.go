package eqspec

import (
	"strings"
	"testing"
)

func TestSuggestCompletions(t *testing.T) {
	sug := SuggestCompletions("water ")
	if len(sug)==0 {
		t.Fatal("No suggestions generated")
	}
	if len(sug)>100 {
		t.Fatal("Too many suggestions")
	}
	for _, s := range sug {
		if !strings.HasPrefix(strings.ToUpper(s), "WATER ") {
			t.Fatal("Bad suggestion")
		}
	}
}
