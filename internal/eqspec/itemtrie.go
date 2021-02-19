// +build wasm,electron

package eqspec

import (
	"github.com/GontikR99/chillmodeinfo/pkg/console"
	"sort"
	"strings"
	"time"
)

func BuildTrie() {
	go func() {
		startTime := time.Now()
		iteration := 0
		for _, v := range everquestItems {
			if iteration%1000 == 0 {
				<-time.After(10*time.Millisecond)
			}
			iteration++
			builtTrie = builtTrie.With(v)
		}
		endTime := time.Now()
		console.Logf("Search trie built, after %v", endTime.Sub(startTime))
	}()
}

// Attempt to find any items
func ScanForItems(textLine string) []string {
	return builtTrie.Scan(textLine)
}

type trieChildEntry struct {
	Char rune
	NextIdx int
}

type trieEntryMap []trieChildEntry

func (tem trieEntryMap) find(c rune) (int, bool) {
	for _, v := range tem {
		if v.Char==c {
			return v.NextIdx, true
		}
	}
	return -1, false
}

func (tem trieEntryMap) with(c rune, idx int) trieEntryMap {
	for idx, v := range tem {
		if v.Char==c {
			tem[idx].NextIdx=idx
			return tem
		}
	}
	tem=append(tem, trieChildEntry{c, idx})
	return tem
}

type itemTrieNode struct {
	Children trieEntryMap
	IsItem   bool
}

func newItemTrieNode() itemTrieNode {
	return itemTrieNode{}
}

type itemTrie []itemTrieNode

func NewItemTrie() itemTrie {
	return []itemTrieNode{newItemTrieNode()}
}

func (trie itemTrie) With(itemName string) itemTrie {
	curIdx := 0
	for _, c := range itemName {
		if newIdx, ok := trie[curIdx].Children.find(c); ok {
			curIdx = newIdx
		} else {
			childIdx := len(trie)
			trie = append(trie, newItemTrieNode())
			trie[curIdx].Children=trie[curIdx].Children.with(c,childIdx)
			curIdx = childIdx
		}
	}
	trie[curIdx].IsItem=true
	return trie
}

// Search for mentions of recognized EverQuest items within a string, using a prebuilt item trie
func (trie itemTrie) Scan(lineText string) []string {
	curState := make(map[int]int)
	curState[0]=0
	var found []string
	for curOffset, c := range lineText {
		nextState := make(map[int]int)
		for stateIdx, startOffset := range curState {
			if trie[stateIdx].IsItem {
				name := lineText[startOffset:curOffset]
				found = append(found, name)
			}
			if nextStateIdx, ok := trie[stateIdx].Children.find(c); ok {
				nextState[nextStateIdx]=startOffset
			}
		}
		nextState[0]=curOffset+1
		curState = nextState
	}
	for stateIdx, startOffset := range curState {
		if trie[stateIdx].IsItem {
			name := lineText[startOffset:]
			found = append(found, name)
		}
	}
	sort.Sort(byLengthReversed(found))
	return found
}


type LexOrderIgnoreCase []string
func (l LexOrderIgnoreCase) Len() int {return len(l)}
func (l LexOrderIgnoreCase) Less(i, j int) bool {return strings.ToUpper(l[i])<strings.ToUpper(l[j])}
func (l LexOrderIgnoreCase) Swap(i, j int) {l[i],l[j] = l[j],l[i]}

type byLengthReversed []string

func (b byLengthReversed) Len() int {return len(b)}
func (b byLengthReversed) Less(i, j int) bool {return len(b[i])>len(b[j])}
func (b byLengthReversed) Swap(i, j int) {b[i],b[j] = b[j],b[i]}

var builtTrie = NewItemTrie()
