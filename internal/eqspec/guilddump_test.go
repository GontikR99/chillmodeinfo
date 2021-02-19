package eqspec

import (
	"fmt"
	"strings"
	"testing"
)

const guildDump="Adax\t65\tBeastlord\tRetired\t\t06/06/20\t\t\t\toff\toff\t0\t\t\t\n" +
	"Alexandre\t70\tPaladin\tRaider\t\t02/13/21\t\t\t\ton\toff\t1290\t11/25/20\t\t\n" +
	"Allak\t65\tMagician\tRetired\t\t08/07/20\t\tNeedsbuffs\t\toff\toff\t0\t\tNeedsbuffs\t\n" +
	"Allegory\t65\tEnchanter\tRetired\t\t08/10/20\t\tWordsworth\t\toff\toff\t0\t\tWordsworth\t\n" +
	"Amras\t65\tRanger\tRetired\t\t06/05/20\t\t\t\toff\toff\t0\t\t\t\n" +
	"Anterra\t70\tMagician\tOfficer Box/Alt\tA\t02/08/21\t\tTheophania\t\toff\toff\t0\t\tTheophania\t\n" +
	"Araushnaee\t70\tNecromancer\tRaider\t\t02/16/21\tGuild Lobby\t\t\ton\toff\t44719\t02/09/21\t\t\n" +
	"Artyom\t70\tEnchanter\tBox/Alt\tA\t02/08/21\t\tSweed\t\toff\toff\t0\t\tSweed\t\n" +
	"Aryani\t70\tMonk\tRaider\t\t02/16/21\t\t\t\ton\toff\t134937\t01/12/21\t\t\n" +
	"Asteria\t67\tWarrior\tRetired\t\t02/15/21\t\tTheophania\t\toff\toff\t0\t\tTheophania\t\n" +
	"Azumen\t67\tMagician\tBox/Alt\tA\t02/13/21\t\tYasstiny\t\toff\toff\t1254\t01/06/21\tYasstiny\t\n" +
	"Bagdamagus\t70\tWarrior\tBox/Alt\tA\t01/16/21\t\tParadigmm\t\toff\toff\t0\t\tParadigmm\t\n" +
	"Banchie\t65\tCleric\tRetired\t\t07/01/20\t\t\t\toff\toff\t0\t\t\t\n" +
	"Bardarsed\t70\tBard\tBox/Alt\tA\t01/16/21\t\tRivix\t\toff\toff\t0\t\tRivix\t\n" +
	"Barks\t65\tDruid\tRetired\t\t05/13/20\t\t\t\toff\toff\t0\t\t\t\n" +
	"Barlain\t65\tRanger\tRetired\t\t09/09/20\t\t\t\toff\toff\t0\t\t\t\n" +
	"Beasttiren\t70\tBeastlord\tRaider\t\t02/15/21\t\t\t\ton\toff\t7276\t01/27/21\t\t\n" +
	"Beatn\t70\tBard\tRetired\t\t11/19/20\t\tNeedsbuffs\t\toff\toff\t0\t\tNeedsbuffs\t\n" +
	"Beatzyck\t70\tBard\tBox/Alt\tA\t02/16/21\tGuild Hall\tMorzyck\t\toff\toff\t0\t\tMorzyck\t\n" +
	"Beechaayas\t70\tMonk\tRaider\t\t02/15/21\t\t\t\toff\toff\t0\t\t\t\n" +
	"Beppan\t54\tDruid\tOfficer Box/Alt\tA\t02/15/21\t\tRektor\t\toff\toff\t3360\t02/03/21\tRektor\t\n" +
	"Bewda\t65\tDruid\tRetired\t\t02/05/21\t\t\t\toff\toff\t0\t\t\t\n" +
	"Bigold\t65\tShadow Knight\tRetired\t\t10/25/20\t\tNaecoeht\t\toff\toff\t0\t\tNaecoeht\t\n" +
	"Bitesize\t66\tWarrior\tRetired\t\t01/22/21\t\t\t\toff\toff\t0\t\t\t\n" +
	"Blacktears\t65\tEnchanter\tRetired\t\t01/28/21\t\t\t\toff\toff\t0\t\t\t\n" +
	"Blaen\t70\tShaman\tBox/Alt\tA\t02/07/21\t\tMarm\t\toff\toff\t0\t\tMarm\t\n" +
	"Blinq\t60\tDruid\tOfficer Box/Alt\t\t02/16/21\t\tBonq\t\toff\toff\t0\t\tBonq\t\n" +
	"Bluediamond\t13\tNecromancer\tBox/Alt\tA\t02/08/21\t\tMeazzu\t\toff\toff\t0\t\tMeazzu\t\n" +
	"Bodand\t66\tCleric\tRetired\t\t10/12/20\t\tCronx\t\toff\toff\t0\t\tCronx\t\n" +
	"Bogs\t65\tRogue\tRetired\t\t05/22/20\t\t\t\toff\toff\t0\t\t\t\n" +
	"Bolanur\t70\tCleric\tRetired\t\t09/19/20\t\t\t\toff\toff\t0\t\t\t\n" +
	"Bonq\t70\tWarrior\tOfficer\t\t02/16/21\t\t\t\toff\toff\t0\t\t\t\n" +
	"Boomboomchaka\t70\tWizard\tBox/Alt\tA\t02/14/21\t\tNeanu\t\toff\toff\t0\t\tNeanu\t\n"

func TestParseGuildDump(t *testing.T) {
	file := strings.NewReader(guildDump)
	members, err := ParseGuildDump(file)
	if err!=nil {
		t.Fatal(err)
	}
	for _, v := range members {
		fmt.Println(v)
	}
}
