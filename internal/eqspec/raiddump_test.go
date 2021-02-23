package eqspec

import (
	"strings"
	"testing"
)

const raidDump = "1\tUgdar\t70\tWarrior\tGroup Leader\t\t\tYes\t\n" +
	"1\tTyryn\t70\tWarrior\t\t\t\tYes\t\n" +
	"1\tVaikyria\t70\tBard\t\t\t\tNo\t\n" +
	"1\tDalamin\t70\tCleric\t\t\t\tYes\t\n" +
	"1\tKaleesi\t70\tShaman\t\t\t\tYes\t\n" +
	"1\tDoctora\t70\tCleric\t\t\t\tYes\t\n" +
	"2\tCrogg\t70\tBerserker\tGroup Leader\t\t\tYes\t\n" +
	"2\tJephine\t70\tShadow Knight\t\t\t\tYes\t\n" +
	"2\tJoram\t70\tShaman\t\t\t\tNo\t\n" +
	"2\tAryani\t70\tMonk\t\t\t\tYes\t\n" +
	"2\tTwistering\t70\tBard\t\t\t\tYes\t\n"

func TestParseRaidDump(t *testing.T) {
	members, err := ParseRaidDump(strings.NewReader(raidDump))
	if err != nil {
		t.Fatal(err)
	}

	for _, member := range members {
		t.Log(member)
	}
}
