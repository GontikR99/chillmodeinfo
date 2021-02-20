package eqspec

import (
	"sort"
	"strings"
)

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

type ItemTrie []itemTrieNode

func NewItemTrie() ItemTrie {
	return []itemTrieNode{newItemTrieNode()}
}

func (trie ItemTrie) With(itemName string) ItemTrie {
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
func (tr CompressedItemTrie) Scan(lineText string) []string {
	curState := make(map[int]int)
	curState[0]=0
	var found []string
	for curOffset, c := range lineText {
		nextState := make(map[int]int)
		for stateIdx, startOffset := range curState {
			if tr.isAccept(stateIdx) {
				name := lineText[startOffset:curOffset]
				found = append(found, name)
			}
			if nextStateIdx := tr.transition(stateIdx, c); nextStateIdx!=0 {
				nextState[nextStateIdx]=startOffset
			}
		}
		nextState[0]=curOffset+1
		curState = nextState
	}
	for stateIdx, startOffset := range curState {
		if tr.isAccept(stateIdx) {
			name := lineText[startOffset:]
			found = append(found, name)
		}
	}
	sort.Sort(byLengthReversed(found))
	return found
}

func (trie ItemTrie) Compress() CompressedItemTrie {
	var transitions CompressedItemTrieTransitions
	var accepts []int
	for sourceIdx, sourceNode := range trie {
		if sourceNode.IsItem {
			accepts = append(accepts, sourceIdx)
		}
		for _, child := range sourceNode.Children {
			transitions = append(transitions, newCompressedItemTrieNode(sourceIdx, child.Char, child.NextIdx))
		}
	}
	sort.Sort(transitions)
	sort.Sort(byValue(accepts))
	return CompressedItemTrie{
		Transitions: transitions,
		Accepts:     accepts,
	}
}

type LexOrderIgnoreCase []string
func (l LexOrderIgnoreCase) Len() int {return len(l)}
func (l LexOrderIgnoreCase) Less(i, j int) bool {return strings.ToUpper(l[i])<strings.ToUpper(l[j])}
func (l LexOrderIgnoreCase) Swap(i, j int) {l[i],l[j] = l[j],l[i]}

type byLengthReversed []string

func (b byLengthReversed) Len() int {return len(b)}
func (b byLengthReversed) Less(i, j int) bool {return len(b[i])>len(b[j])}
func (b byLengthReversed) Swap(i, j int) {b[i],b[j] = b[j],b[i]}

type CompressedItemTrieTransition uint64
func newCompressedItemTrieNode(sourceState int, char rune, destState int) CompressedItemTrieTransition {
	return CompressedItemTrieTransition((uint64(sourceState)<<40)|(uint64(char)<<32)|uint64(destState))
}
func (tn CompressedItemTrieTransition) sourceState() int {return int(tn>>40)}
func (tn CompressedItemTrieTransition) char() rune       {return rune((tn>>32)&0xff)}
func (tn CompressedItemTrieTransition) destState() int   {return int(tn&0xffffffff)}

type CompressedItemTrieTransitions []CompressedItemTrieTransition
func (tr CompressedItemTrieTransitions) Len() int           {return len(tr)}
func (tr CompressedItemTrieTransitions) Less(i, j int) bool {return tr[i]>>32 < tr[j]>>32}
func (tr CompressedItemTrieTransitions) Swap(i, j int)      { tr[i], tr[j] = tr[j], tr[i]}

type byValue []int
func (b byValue) Len() int {return len(b)}
func (b byValue) Less(i, j int) bool {return b[i]<b[j]}
func (b byValue) Swap(i, j int) {b[i],b[j] = b[j],b[i]}

type CompressedItemTrie struct {
	Transitions CompressedItemTrieTransitions
	Accepts     []int
}

func (tr CompressedItemTrie) transition(sourceState int, char rune) (deststate int) {
	target:=newCompressedItemTrieNode(sourceState, char, 0)
	loc := sort.Search(len(tr.Transitions), func(i int) bool {return target <= tr.Transitions[i]})
	if loc==len(tr.Transitions) {return 0}
	if tr.Transitions[loc].sourceState()!=sourceState || tr.Transitions[loc].char()!=char {return 0}
	return tr.Transitions[loc].destState()
}

func (tr CompressedItemTrie) isAccept(state int) bool {
	loc := sort.Search(len(tr.Accepts), func(i int)bool {return state<=tr.Accepts[i]})
	return loc!=len(tr.Accepts) && state== tr.Accepts[loc]
}