package main

import (
	"fmt"
	"os"

	zpl "github.com/mboeh/zplating/internal"
)

func printZPL(pgm []zpl.Command) {
	prefixes := map[rune]rune{
		'^': '^',
		'~': '~',
	}
	for i := range pgm {
		cmd := pgm[i]
		prefix := prefixes[rune(cmd.Command[0])]
		rest := cmd.Command[1:]
		fmt.Printf("%s%s%s\n", string(prefix), rest, cmd.Argument)
		if rest == "CC" {
			// Change caret for future output
			prefixes['^'] = rune(cmd.Argument[0])
		} else if rest == "CT" {
			// Change tilde for future output
			prefixes['~'] = rune(cmd.Argument[0])
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
		fmt.Fprintf(os.Stderr, "ERROR: %s", parser.Error)
		os.Exit(1)
	} else {
		printZPL(parser.Commands)
	}
}
