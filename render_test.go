package figlet4go

import (
	"strings"
	"testing"
)

func TestSVGRendering(t *testing.T) {
	ascii := NewAsciiRender()

	tests := []struct {
		name     string
		input    string
		fontName string
	}{
		{
			name:     "simple text with standard font",
			input:    "Hi",
			fontName: "standard",
		},
		{
			name:     "single character",
			input:    "A",
			fontName: "standard",
		},
		{
			name:     "text with spaces",
			input:    "A B",
			fontName: "standard",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := NewRenderOptions()
			parser, err := GetParser("svg")
			if err != nil {
				t.Fatalf("Failed to get SVG parser: %v", err)
			}
			options.Parser = *parser
			options.FontName = tt.fontName

			result, err := ascii.RenderOpts(tt.input, options)
			if err != nil {
				t.Fatalf("RenderOpts() error = %v", err)
			}

			// Verify SVG structure
			if !strings.HasPrefix(result, "<text>") {
				t.Errorf("SVG output should start with <text>, got: %s", result[:min(20, len(result))])
			}
			if !strings.HasSuffix(result, "</text>") {
				t.Errorf("SVG output should end with </text>, got: %s", result[max(0, len(result)-20):])
			}

			// Verify line breaks are SVG-style
			if strings.Contains(result, "\n") {
				// The result will contain newlines between lines
				if !strings.Contains(result, "<br/>") {
					t.Errorf("SVG output should contain <br/> tags")
				}
			}

			// Verify spaces are entity-encoded
			if strings.Contains(tt.input, " ") {
				if !strings.Contains(result, "&#160;") {
					t.Errorf("SVG output should contain &#160; for spaces, got: %s", result)
				}
			}
		})
	}
}

func TestCompareParserOutputs(t *testing.T) {
	ascii := NewAsciiRender()
	input := "AB"

	// Render with terminal parser
	termOptions := NewRenderOptions()
	termOptions.FontName = "standard"
	termResult, err := ascii.RenderOpts(input, termOptions)
	if err != nil {
		t.Fatalf("Terminal render failed: %v", err)
	}

	// Render with SVG parser
	svgOptions := NewRenderOptions()
	svgParser, err := GetParser("svg")
	if err != nil {
		t.Fatalf("Failed to get SVG parser: %v", err)
	}
	svgOptions.Parser = *svgParser
	svgOptions.FontName = "standard"
	svgResult, err := ascii.RenderOpts(input, svgOptions)
	if err != nil {
		t.Fatalf("SVG render failed: %v", err)
	}

	// Render with HTML parser
	htmlOptions := NewRenderOptions()
	htmlParser, err := GetParser("html")
	if err != nil {
		t.Fatalf("Failed to get HTML parser: %v", err)
	}
	htmlOptions.Parser = *htmlParser
	htmlOptions.FontName = "standard"
	htmlResult, err := ascii.RenderOpts(input, htmlOptions)
	if err != nil {
		t.Fatalf("HTML render failed: %v", err)
	}

	// SVG and HTML should be different due to different prefixes/suffixes
	if svgResult == htmlResult {
		t.Error("SVG and HTML output should differ")
	}

	// SVG should have <text> tags
	if !strings.Contains(svgResult, "<text>") || !strings.Contains(svgResult, "</text>") {
		t.Error("SVG output missing <text> tags")
	}

	// HTML should have <code> tags
	if !strings.Contains(htmlResult, "<code>") || !strings.Contains(htmlResult, "</code>") {
		t.Error("HTML output missing <code> tags")
	}

	// Terminal should have neither
	if strings.Contains(termResult, "<") {
		t.Error("Terminal output should not contain HTML/SVG tags")
	}
}

func TestSVGRenderWithColors(t *testing.T) {
	ascii := NewAsciiRender()

	options := NewRenderOptions()
	parser, err := GetParser("svg")
	if err != nil {
		t.Fatalf("Failed to get SVG parser: %v", err)
	}
	options.Parser = *parser
	options.FontName = "standard"
	options.FontColor = []Color{
		ColorRed,
		ColorGreen,
		ColorBlue,
	}

	result, err := ascii.RenderOpts("RGB", options)
	if err != nil {
		t.Fatalf("RenderOpts() with colors error = %v", err)
	}

	// Should still have SVG structure
	if !strings.HasPrefix(result, "<text>") {
		t.Error("Colored SVG output should start with <text>")
	}
	if !strings.HasSuffix(result, "</text>") {
		t.Error("Colored SVG output should end with </text>")
	}

	// Should contain color codes (ANSI or similar)
	// The exact format depends on the Color implementation
	if len(result) < len("<text></text>") {
		t.Error("Colored SVG output seems too short")
	}
}

func TestEmptyStringRender(t *testing.T) {
	ascii := NewAsciiRender()

	parsers := []string{"terminal", "html", "svg"}

	for _, parserKey := range parsers {
		t.Run("empty_string_"+parserKey, func(t *testing.T) {
			options := NewRenderOptions()
			parser, err := GetParser(parserKey)
			if err != nil {
				t.Fatalf("Failed to get parser: %v", err)
			}
			options.Parser = *parser

			result, err := ascii.RenderOpts("", options)
			if err != nil {
				t.Fatalf("RenderOpts() error = %v", err)
			}

			// Should still have proper structure
			if !strings.HasPrefix(result, parser.Prefix) {
				t.Errorf("Output should start with prefix %q", parser.Prefix)
			}
			if !strings.HasSuffix(result, parser.Suffix) {
				t.Errorf("Output should end with suffix %q", parser.Suffix)
			}
		})
	}
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
