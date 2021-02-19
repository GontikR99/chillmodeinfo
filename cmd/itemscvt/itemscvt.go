package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/GontikR99/chillmodeinfo/internal/eqspec"
	"os"
	"sort"
	"strings"
)

func main() {
	if len(os.Args)!=3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <listing> <filename>\n", os.Args[0])
	}
	dbFile, err := os.Open(os.Args[2])
	if err!=nil {
		panic(err)
	}
	defer dbFile.Close()
	lineFile, err := gzip.NewReader(dbFile)
	if err!=nil {
		panic(err)
	}
	defer lineFile.Close()
	lineScanner := bufio.NewScanner(lineFile)
	lineScanner.Scan() // get rid of header
	switch os.Args[1] {
	case "listing":
		var itemNames []string
		for lineScanner.Scan() {
			line := lineScanner.Text()
			fields := strings.Split(line,"|")
			if len(fields)<2 {
				continue
			}
			itemName:=fields[1]
			itemNames = append(itemNames, itemName)
		}
		sort.Sort(eqspec.LexOrderIgnoreCase(itemNames))
		fmt.Println("package eqspec")
		fmt.Println()
		fmt.Println("var everquestItems=[]string{")
		for _, name := range itemNames {
			v, _ := json.Marshal(name)
			fmt.Println("    "+string(v)+",")
		}
		fmt.Println("}")

	default:
		fmt.Fprintln(os.Stderr, "Unknown subcommand "+os.Args[1])
	}
}