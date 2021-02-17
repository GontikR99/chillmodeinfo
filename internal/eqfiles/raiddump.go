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

func ParseRaidDump(in io.Reader) ([]record.Member, error) {
	var records []record.Member
	scan := bufio.NewScanner(in)
	for scan.Scan() {
		line := scan.Text()
		fields := strings.Split(line, "\t")

		if len(fields) != 9 {
			return nil, errors.New("Not a valid raid dump: should have 9 fields per entry, found an entry with " + strconv.Itoa(len(fields)))
		}

		_, err := strconv.Atoi(fields[0])
		if err!=nil {
			return nil, errors.New("Not a valid raid dump: first field should be a number")
		}

		name := fields[1]
		if name=="" {
			return nil, errors.New("Not a valid raid dump: second field should be a name")
		}

		level, err := strconv.Atoi(fields[2])
		if err!=nil || level<=0 || level>=150 {
			return nil, errors.New("Not a valid raid dump: third field should be a level (number)")
		}

		class := fields[3]
		if _, ok := ClassMap[class]; !ok {
			return nil, errors.New("Not a valid raid dump: unrecognized class "+class)
		}
		records=append(records, &record.BasicMember{
			Name:       name,
			Class:      class,
			Level:      int16(level),
			Alt:        false,
			DKP:        0,
			LastActive: time.Time{},
			Owner:      "",
		})
	}
	return records, nil
}