/*
   Copyright 2020 Matthew Boeh

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package zpl

import (
	"fmt"
	"unicode"
)

type State int

const (
	READY State = iota
	CARET_COMMAND
	TILDE_COMMAND
	ARGUMENTS
	ERROR
)

type Parser struct {
	State    State
	Error    string
	Commands []Command

	row        int
	column     int
	caret      rune
	tilde      rune
	delimiter  rune
	currentCmd string
	currentArg string
	args       []string
	knownCmds  map[string][]formatToken
}

func newParser() *Parser {
	return &Parser{
		row:        1,
		column:     1,
		caret:      '^',
		tilde:      '~',
		delimiter:  ',',
		State:      READY,
		Error:      "",
		Commands:   make([]Command, 0),
		currentCmd: "",
		currentArg: "",
		args:       make([]string, 0),
		knownCmds:  commands(),
	}
}

func (p *Parser) feedString(pgm string) bool {
	for _, c := range pgm {
		if !p.feed(c) {
			return false
		}
	}
	return true
}

func (p *Parser) feed(char rune) bool {
	if char == '\n' {
		p.row += 1
		p.column = 1
		return p.State != ERROR
	}
	p.column += 1
	switch p.State {
	case READY:
		// Ready for a new command
		if unicode.IsSpace(char) {
			return true
		}
		if char == p.caret {
			p.State = CARET_COMMAND
			p.currentCmd = "^"
			p.currentArg = ""
		} else if char == p.tilde {
			p.State = TILDE_COMMAND
			p.currentCmd = "~"
			p.currentArg = ""
		} else {
			p.fail("Expected caret (" + string(p.caret) + ") or tilde (" + string(p.tilde) + ")")
			return false
		}
	case TILDE_COMMAND:
		fallthrough
	case CARET_COMMAND:
		// Trying to complete a command
		p.currentCmd += string(char)
		cmd, found := p.knownCmds[p.currentCmd]
		if found {
			if len(cmd) > 0 {
				p.State = ARGUMENTS
			} else {
				// Immediately start reading another command
				p.emit()
			}
		} else if len(p.currentCmd) > 4 {
			p.fail("Invalid command: " + p.currentCmd)
		}
		return true
	case ARGUMENTS:
		cmdFmt := p.currentFmt()
		switch cmdFmt[len(p.args)] {
		case FMT_BYTE:
			p.currentArg = string(char)
			p.finishArg()
		case FMT_PARAM:
			if char == p.delimiter {
				p.finishArg()
				return true
			}
			fallthrough
		case FMT_TEXT:
			if char == p.caret || char == p.tilde {
				p.finishArg()
				return p.feed(char)
			} else {
				p.currentArg += string(char)
			}
		case FMT_NONE:
			p.fail("too many arguments")
			return false
		case FMT_DELIMITER:
			if char != p.delimiter {
				p.fail("expected delimiter " + string(p.delimiter) + " got " + string(char))
				return false
			}
		}
		return true
	case ERROR:
		// An error has occurred, no more commands are accepted
		return false
	}
	return true
}

func (p *Parser) finishArg() {
	p.args = append(p.args, p.currentArg)
	p.currentArg = ""
	fmt := p.currentFmt()
	if len(p.args) == len(fmt) {
		p.emit()
	}
}

func (p *Parser) emit() {
	fmt := p.currentFmt()
	if len(fmt) > len(p.args) {
		p.fail("too few arguments: " + p.currentCmd)
		return
	}

	// Special handling for Change Caret, Change Tilde, and Change Delimiter
	if p.currentCmd == "^CC" || p.currentCmd == "~CC" {
		p.caret = rune(p.args[0][0])
	} else if p.currentCmd == "~CT" || p.currentCmd == "^CT" {
		p.tilde = rune(p.args[0][0])
	} else if p.currentCmd == "^CD" || p.currentCmd == "~CD" {
		p.delimiter = rune(p.args[0][0])
	}
	p.Commands = append(p.Commands, Command{
		Command:   p.currentCmd,
		Arguments: p.args,
	})
	p.currentCmd = ""
	p.currentArg = ""
	p.args = make([]string, 0)
	p.State = READY
}

func (p *Parser) fail(err string) {
	p.State = ERROR
	p.Error = fmt.Sprintf("%d:%d:%s", p.row, p.column, err)
}

func (p *Parser) currentFmt() []formatToken {
	fmt, ok := p.knownCmds[p.currentCmd]
	if !ok {
		panic("called currentFmt() outside valid command")
	}
	return fmt
}
