package figlet4go

import "errors"

// Parser stores some output specific stuff
type Parser struct {
	// Name used for switching in colors.go
	Name string
	// Parser prefix
	Prefix string
	// Parser suffix
	Suffix     string
	LinePrefix string
	LineSuffix string
	// Things to be replaced (f.e. \n to <br>)
	Replaces map[string]string
}

var parsers map[string]Parser = map[string]Parser{
	// Default terminal parser
	"terminal": {"terminal", "", "", "", "\n", nil},
	// Parser for HTML code
	"html": {"html", "<code>", "</code>", "", "<br>", map[string]string{" ": "&nbsp;"}},
	// Parser for SVG
	"svg": {"svg", "", "", "<tspan x=\"10\" dy=\"1em\">", "</tspan>", map[string]string{" ": "&#160;"}},
}

// GetParser returns a parser by its key
func GetParser(key string) (*Parser, error) {
	parser, ok := parsers[key]
	if !ok {
		return nil, errors.New("Invalid parser key: " + key)
	}
	return &parser, nil
}
