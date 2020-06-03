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

type Command struct {
	Command   string
	Arguments []string
}

type formatToken int

const (
	FMT_NONE formatToken = iota
	// Read a single byte as a single character string
	FMT_BYTE
	// Read up to the next control character as a string
	FMT_TEXT
	// Read up to the next delimiter as a string
	FMT_PARAM
	// Read and discard a delimiter
	FMT_DELIMITER
)

// Commands

func commands() map[string][]formatToken {
	p0 := []formatToken{}
	p1 := []formatToken{FMT_TEXT}
	p2 := []formatToken{FMT_PARAM, FMT_PARAM}
	p3 := []formatToken{FMT_PARAM, FMT_PARAM, FMT_PARAM}
	//p4 := []formatToken{FMT_PARAM,FMT_PARAM,FMT_PARAM,FMT_PARAM}
	p5 := []formatToken{FMT_PARAM, FMT_PARAM, FMT_PARAM, FMT_PARAM, FMT_PARAM}
	p6 := []formatToken{FMT_PARAM, FMT_PARAM, FMT_PARAM, FMT_PARAM, FMT_PARAM}
	pbyte := []formatToken{FMT_BYTE}

	return map[string][]formatToken{
		"^A":  {FMT_BYTE, FMT_BYTE, FMT_DELIMITER, FMT_PARAM, FMT_PARAM},
		"^BC": p6,
		"^BY": p3,
		"^CC": pbyte,
		"~CC": pbyte,
		"^CD": pbyte,
		"~CD": pbyte,
		"^CF": p3,
		"^CT": pbyte,
		"~CT": pbyte,
		"^FD": p1,
		"^FO": p2,
		"^FR": p0,
		"^FS": p0,
		"^FX": p1,
		"^GB": p5,
		"^XA": p0,
		"^XZ": p0,
	}
}
