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

type ParsingState int

const (
	// Awaiting a command to parse.
	READY ParsingState = iota
	// The caret (^) was found; reading bytes until a command name is found.
	CARET_COMMAND
	// The tilde (~) was found; reading bytes until a command name is found.
	TILDE_COMMAND
	// A command name was found; parsing the arguments according to that command's format.
	ARGUMENTS
	// An error occurred and parsing cannot continue.
	ERROR
)

func stateStr(state ParsingState) string {
	switch state {
	case READY:
		return "READY"
	case CARET_COMMAND:
		return "CARET_COMMAND"
	case TILDE_COMMAND:
		return "TILDE_COMMAND"
	case ARGUMENTS:
		return "ARGUMENTS"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

type Parser struct {
	State    ParsingState
	Error    string
	Commands []Command

	// Source position, for diagnostics.
	row    int
	column int

	// Special characters, which can be modified on the fly by commands.
	caret     rune
	tilde     rune
	delimiter rune

	// State for parsing a single command.
	currentCmd string
	currentArg string
	args       []string
	argn       int

	// Lexicon of known (not necessarily fully supported) commands.
	knownCmds map[string][]formatToken
}

// Create a new Parser, ready to start receiving commands.
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
		argn:       0,
		knownCmds:  commands(),
	}
}

// Parse a text string (byte by byte), returning false on error.
// If this returns false, the error is in the Parser's Error field.
func (p *Parser) feedString(pgm string) bool {
	for _, c := range pgm {
		if !p.feed(c) {
			return false
		}
	}
	return true
}

// Parse a single byte, returning false on error.
// If this returns false, the error is in the Parser's Error field.
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
		cmdFmt := p.currentFmt()[p.argn]
		switch cmdFmt {
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
			} else {
				p.argn += 1
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
	p.argn += 1
	p.currentArg = ""
	fmt := p.currentFmt()
	if p.argn == len(fmt) {
		p.emit()
	}
}

func (p *Parser) emit() {
	fmtc := p.currentFmt()
	if len(fmtc) > p.argn {
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
	p.argn = 0
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
