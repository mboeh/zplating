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
		"^ADN": p3,
		"^BC":  p6,
		"^BY":  p3,
		"^CC":  pbyte,
		"~CC":  pbyte,
		"^CD":  pbyte,
		"~CD":  pbyte,
		"^CF":  p3,
		"^CT":  pbyte,
		"~CT":  pbyte,
		"^FD":  p1,
		"^FO":  p2,
		"^FR":  p0,
		"^FS":  p0,
		"^FX":  p1,
		"^GB":  p5,
		"^XA":  p0,
		"^XZ":  p0,
	}
}
