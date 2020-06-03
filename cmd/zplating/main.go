package main

import (
	"fmt"
	"os"
	"strings"

	zpl "github.com/mboeh/zplating/pkg/zpl"
)

func printZPL(pgm []zpl.Command) {
	delimiter := ","
	prefixes := map[rune]rune{
		'^': '^',
		'~': '~',
	}
	for i := range pgm {
		cmd := pgm[i]
		prefix := prefixes[rune(cmd.Command[0])]
		rest := cmd.Command[1:]
		fmt.Printf("%s%s%s\n", string(prefix), rest, strings.Join(cmd.Arguments, delimiter))
		if rest == "CC" {
			// Change caret for future output
			prefixes['^'] = rune(cmd.Arguments[0][0])
		} else if rest == "CT" {
			// Change tilde for future output
			prefixes['~'] = rune(cmd.Arguments[0][0])
		} else if rest == "CD" {
			// Changd delimiter for future output
			delimiter = cmd.Arguments[0]
		}
	}
}
func main() {
	filename := os.Args[1]
	f, err := os.Open(filename)
	if err != nil {
		panic("whoops: " + err.Error())
	}
	parser, err := zpl.Parse(f)
	if err != nil {
		panic("burp: " + err.Error())
	}
	if parser.State == zpl.ERROR {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", parser.Error)
		os.Exit(1)
	} else {
		printZPL(parser.Commands)
	}
}
