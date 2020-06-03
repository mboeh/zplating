package zpl

import (
	"io"
)

func Parse(rd io.Reader) (parser *Parser, err error) {
	buf := make([]byte, 1000)
	parser = newParser()
	for {
		n, err := rd.Read(buf)
		if err != nil {
			break
		}
		if !parser.feedString(string(buf[:n])) {
			break
		}
	}
	return
}
