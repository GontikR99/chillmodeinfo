package eqspec

type ClassStat struct {
	Trigraph string
}

var ClassMap = map[string]ClassStat{
	"Bard":          {"BRD"},
	"Beastlord":     {"BST"},
	"Berserker":     {"BER"},
	"Cleric":        {"CLR"},
	"Druid":         {"DRU"},
	"Enchanter":     {"ENC"},
	"Magician":      {"MAG"},
	"Monk":          {"MNK"},
	"Necromancer":   {"NEC"},
	"Paladin":       {"PAL"},
	"Ranger":        {"RNG"},
	"Rogue":         {"ROG"},
	"Shadow Knight": {"SHD"},
	"Shaman":        {"SHM"},
	"Warrior":       {"WAR"},
	"Wizard":        {"WIZ"},
}
