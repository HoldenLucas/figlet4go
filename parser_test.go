package figlet4go

import (
	"testing"
)

func TestGetParser(t *testing.T) {
	tests := []struct {
		name      string
		parserKey string
		wantErr   bool
	}{
		{
			name:      "terminal parser exists",
			parserKey: "terminal",
			wantErr:   false,
		},
		{
			name:      "html parser exists",
			parserKey: "html",
			wantErr:   false,
		},
		{
			name:      "svg parser exists",
			parserKey: "svg",
			wantErr:   false,
		},
		{
			name:      "invalid parser returns error",
			parserKey: "invalid",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, err := GetParser(tt.parserKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetParser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && parser == nil {
				t.Errorf("GetParser() returned nil parser for valid key %s", tt.parserKey)
			}
		})
	}
}

func TestSVGParser(t *testing.T) {
	parser, err := GetParser("svg")
	if err != nil {
		t.Fatalf("Failed to get SVG parser: %v", err)
	}

	// Verify SVG parser properties
	if parser.Name != "svg" {
		t.Errorf("SVG parser Name = %v, want %v", parser.Name, "svg")
	}
	if parser.Prefix != "<text>" {
		t.Errorf("SVG parser Prefix = %v, want %v", parser.Prefix, "<text>")
	}
	if parser.Suffix != "</text>" {
		t.Errorf("SVG parser Suffix = %v, want %v", parser.Suffix, "</text>")
	}
	if parser.NewLine != "<br/>" {
		t.Errorf("SVG parser NewLine = %v, want %v", parser.NewLine, "<br/>")
	}

	// Verify space replacement
	if parser.Replaces == nil {
		t.Fatal("SVG parser Replaces is nil")
	}
	if replacement, ok := parser.Replaces[" "]; !ok || replacement != "&#160;" {
		t.Errorf("SVG parser space replacement = %v, want %v", replacement, "&#160;")
	}
}

func TestParserReplaces(t *testing.T) {
	tests := []struct {
		name      string
		parserKey string
		input     string
		expected  string
	}{
		{
			name:      "terminal parser no replaces",
			parserKey: "terminal",
			input:     "Hello World",
			expected:  "Hello World",
		},
		{
			name:      "html parser replaces spaces",
			parserKey: "html",
			input:     "Hello World",
			expected:  "Hello&nbsp;World",
		},
		{
			name:      "svg parser replaces spaces",
			parserKey: "svg",
			input:     "Hello World",
			expected:  "Hello&#160;World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, err := GetParser(tt.parserKey)
			if err != nil {
				t.Fatalf("Failed to get parser: %v", err)
			}

			result := handleReplaces(tt.input, *parser)
			if result != tt.expected {
				t.Errorf("handleReplaces() = %v, want %v", result, tt.expected)
			}
		})
	}
}
