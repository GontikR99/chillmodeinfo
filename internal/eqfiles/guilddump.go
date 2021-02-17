package eqfiles

import (
	"bufio"
	"errors"
	"github.com/GontikR99/chillmodeinfo/internal/record"
	"io"
	"strconv"
	"strings"
	"time"
)

func ParseGuildDump(re io.Reader) ([]record.Member, error) {
	var records []record.Member
	scan:=bufio.NewScanner(re)
	for scan.Scan() {
		line := scan.Text()
		elements := strings.Split(line, "\t")

		if len(elements)!=15 {
			return nil, errors.New("Not a valid guild dump: should have 15 fields per entry, found an entry with "+strconv.Itoa(len(elements)))
		}
		name := elements[0]
		if len(name)==0 {
			return nil, errors.New("Not a valid guild dump: first field should be a character name")
		}
		level, err := strconv.Atoi(elements[1])
		if err!=nil {
			return nil, errors.New("Not a valid guild dump: second field isn't a level number")
		}
		class := elements[2]
		if _, ok := ClassMap[class]; !ok {
			return nil, errors.New("Not a valid guild dump: unknown class "+class)
		}
		rank := elements[3]
		if rank=="" {
			return nil, errors.New("Not a valid guild dump: missing rank field")
		}
		altFlagStr := elements[4]
		var altFlag bool
		if altFlagStr=="" {
			altFlag = false
		} else if altFlagStr=="A" {
			altFlag = true
		} else {
			return nil, errors.New("Not a valid guild dump: Mangled alt flag")
		}
		lastOnlineParts := strings.Split(elements[5], "/")
		if len(lastOnlineParts)!=3 {
			return nil, errors.New("Not a valid guild dump: failed to parse last online date")
		}
		lastOnlineMonth, err := strconv.Atoi(lastOnlineParts[0])
		if err!=nil || 0>=lastOnlineMonth || 12<lastOnlineMonth {
			return nil, errors.New("Not a valid guild dump: failed to parse last online date")
		}
		lastOnlineDay, err := strconv.Atoi(lastOnlineParts[1])
		if err!=nil || 0>=lastOnlineDay || 31<lastOnlineDay {
			return nil, errors.New("Not a valid guild dump: failed to parse last online date")
		}
		lastOnlineYear, err := strconv.Atoi(lastOnlineParts[2])
		if err!=nil || 0>lastOnlineYear || 99<lastOnlineYear {
			return nil, errors.New("Not a valid guild dump: failed to parse last online date")
		}
		lastOnline := time.Date(2000+lastOnlineYear, time.Month(lastOnlineMonth), lastOnlineDay, 12, 0, 0, 0, time.UTC)

		comment := elements[7]
		records=append(records, &record.BasicMember{
			Name:       name,
			Class:      class,
			Level:      int16(level),
			Rank:       rank,
			Alt:        altFlag,
			DKP:        0,
			LastActive: lastOnline,
			Owner:      comment,
		})
	}
	return records, nil
}