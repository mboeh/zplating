package zpl

import (
	"io"
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

type Command struct {
	Command  string
	Argument string
}

type Parser struct {
	Caret      rune
	Tilde      rune
	State      State
	Error      string
	Commands   []Command
	currentCmd string
	currentArg string
	args       []string
	knownCmds  map[string]int
}

func NewParser() *Parser {
	return &Parser{
		Caret:      '^',
		Tilde:      '~',
		State:      READY,
		Error:      "",
		Commands:   make([]Command, 0),
		currentCmd: "",
		currentArg: "",
		args:       make([]string, 0),
		knownCmds:  commands(),
	}
}

func Parse(rd io.Reader) (parser *Parser, err error) {
	buf := make([]byte, 1000)
	parser = NewParser()
	for {
		n, err := rd.Read(buf)
		if err != nil {
			break
		}
		if !parser.FeedString(string(buf[:n])) {
			break
		}
	}
	return
}

func (p *Parser) FeedString(pgm string) bool {
	for _, c := range pgm {
		if !p.feed(c) {
			return false
		}
	}
	return true
}

func (p *Parser) feed(char rune) bool {
	if char == '\n' {
		return p.State != ERROR
	}
	switch p.State {
	case READY:
		// Ready for a new command
		if unicode.IsSpace(char) {
			return true
		}
		if char == p.Caret {
			p.State = CARET_COMMAND
			p.currentCmd = "^"
			p.currentArg = ""
		} else if char == p.Tilde {
			p.State = TILDE_COMMAND
			p.currentCmd = "~"
			p.currentArg = ""
		} else {
			p.fail("Expected caret (" + string(p.Caret) + ") or tilde (" + string(p.Tilde) + ")")
			return false
		}
	case TILDE_COMMAND:
		fallthrough
	case CARET_COMMAND:
		// Trying to complete a command
		p.currentCmd += string(char)
		if p.foundCommand() {
			if p.knownCmds[p.currentCmd] > 0 {
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
		// Special handling for Change Caret and Change Tilde, which take a single byte argument
		if p.currentCmd == "^CC" || p.currentCmd == "~CC" {
			p.Caret = char
			p.currentArg = string(char)
			p.emit()
		} else if p.currentCmd == "~CT" || p.currentCmd == "^CT" {
			p.Tilde = char
			p.currentArg = string(char)
			p.emit()
		} else if char == p.Caret || char == p.Tilde {
			p.emit()
			return p.feed(char)
		} else {
			p.currentArg += string(char)
		}
		return true
	case ERROR:
		// An error has occurred, no more commands are accepted
		return false
	}
	return true
}

func (p *Parser) emit() {
	p.Commands = append(p.Commands, Command{
		Command:  p.currentCmd,
		Argument: p.currentArg,
	})
	p.currentCmd = ""
	p.currentArg = ""
	p.args = make([]string, 0)
	p.State = READY
}

func (p *Parser) fail(err string) {
	p.State = ERROR
	p.Error = err
}

func (p *Parser) foundCommand() bool {
	_, ok := p.knownCmds[p.currentCmd]
	return ok
}

// Commands

func commands() map[string]int {
	return map[string]int{
		"^ADN": 3,
		"^BC":  6,
		"^BY":  3,
		"^CC":  1,
		"~CC":  1,
		"^CF":  3,
		"^CT":  1,
		"~CT":  1,
		"^FD":  1,
		"^FO":  2,
		"^FR":  0,
		"^FS":  0,
		"^FX":  1,
		"^GB":  5,
		"^XA":  0,
		"^XZ":  0,
	}
}
