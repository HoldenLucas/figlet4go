package figlet4go

import (
	"strings"
	"testing"
)

// TestAllParsers ensures all parsers work correctly
func TestAllParsers(t *testing.T) {
	ascii := NewAsciiRender()
	testText := "OK"

	parsers := []string{"terminal", "html", "svg"}

	for _, parserName := range parsers {
		t.Run(parserName, func(t *testing.T) {
			options := NewRenderOptions()
			parser, err := GetParser(parserName)
			if err != nil {
				t.Fatalf("GetParser(%s) failed: %v", parserName, err)
			}
			options.Parser = *parser

			result, err := ascii.RenderOpts(testText, options)
			if err != nil {
				t.Fatalf("RenderOpts failed: %v", err)
			}

			// Check prefix and suffix
			if !strings.HasPrefix(result, parser.Prefix) {
				t.Errorf("Result should start with prefix %q", parser.Prefix)
			}
			if !strings.HasSuffix(result, parser.Suffix) {
				t.Errorf("Result should end with suffix %q", parser.Suffix)
			}

			// Result should not be empty
			if len(result) == 0 {
				t.Error("Result should not be empty")
			}
		})
	}
}

// TestParserColorSupport ensures colors work with each parser
func TestParserColorSupport(t *testing.T) {
	ascii := NewAsciiRender()
	testText := "RGB"

	tests := []struct {
		parser        string
		expectedColor string
	}{
		{"terminal", "\x1b["},     // ANSI escape codes
		{"html", "<span style="}, // HTML span tags
		{"svg", "<tspan fill="},  // SVG tspan with fill
	}

	for _, tt := range tests {
		t.Run(tt.parser+"_colors", func(t *testing.T) {
			options := NewRenderOptions()
			parser, err := GetParser(tt.parser)
			if err != nil {
				t.Fatalf("GetParser failed: %v", err)
			}
			options.Parser = *parser
			options.FontColor = []Color{ColorRed, ColorGreen, ColorBlue}

			result, err := ascii.RenderOpts(testText, options)
			if err != nil {
				t.Fatalf("RenderOpts failed: %v", err)
			}

			if !strings.Contains(result, tt.expectedColor) {
				t.Errorf("Result should contain color marker %q for parser %s\nGot: %s",
					tt.expectedColor, tt.parser, result[:min(100, len(result))])
			}
		})
	}
}

// TestSVGSpecificFeatures tests SVG-specific functionality
func TestSVGSpecificFeatures(t *testing.T) {
	ascii := NewAsciiRender()

	t.Run("space_encoding", func(t *testing.T) {
		options := NewRenderOptions()
		parser, _ := GetParser("svg")
		options.Parser = *parser

		result, _ := ascii.RenderOpts("A B", options)
		if !strings.Contains(result, "&#160;") {
			t.Error("SVG should encode spaces as &#160;")
		}
	})

	t.Run("line_breaks", func(t *testing.T) {
		options := NewRenderOptions()
		parser, _ := GetParser("svg")
		options.Parser = *parser

		result, _ := ascii.RenderOpts("X", options)
		// Multi-line ASCII art should have <br/> tags
		if strings.Count(result, "\n") > 1 {
			if !strings.Contains(result, "<br/>") {
				t.Error("Multi-line SVG should contain <br/> tags")
			}
		}
	})

	t.Run("color_format", func(t *testing.T) {
		options := NewRenderOptions()
		parser, _ := GetParser("svg")
		options.Parser = *parser
		options.FontColor = []Color{
			TrueColor{r: 255, g: 100, b: 50},
		}

		result, _ := ascii.RenderOpts("X", options)
		if !strings.Contains(result, "rgb(255,100,50)") {
			t.Error("SVG colors should use rgb(r,g,b) format")
		}
		if !strings.Contains(result, "<tspan") && !strings.Contains(result, "</tspan>") {
			t.Error("SVG colored output should use tspan elements")
		}
	})
}

// TestAnsiColorConversion ensures AnsiColor converts properly for SVG
func TestAnsiColorConversion(t *testing.T) {
	ascii := NewAsciiRender()

	options := NewRenderOptions()
	parser, _ := GetParser("svg")
	options.Parser = *parser
	options.FontColor = []Color{ColorRed} // AnsiColor

	result, _ := ascii.RenderOpts("A", options)

	// AnsiColor should be converted to RGB for SVG
	if !strings.Contains(result, "rgb(") {
		t.Error("AnsiColor should be converted to rgb() for SVG")
	}
	// Should contain the red color from the conversion table
	if !strings.Contains(result, "rgb(255,65,54)") {
		t.Error("ColorRed should convert to rgb(255,65,54) for SVG")
	}
}

// TestCompareHTMLvsSVG ensures SVG and HTML produce different but valid outputs
func TestCompareHTMLvsSVG(t *testing.T) {
	ascii := NewAsciiRender()
	testText := "Hi"

	htmlOptions := NewRenderOptions()
	htmlParser, _ := GetParser("html")
	htmlOptions.Parser = *htmlParser
	htmlOptions.FontColor = []Color{ColorRed}

	svgOptions := NewRenderOptions()
	svgParser, _ := GetParser("svg")
	svgOptions.Parser = *svgParser
	svgOptions.FontColor = []Color{ColorRed}

	htmlResult, _ := ascii.RenderOpts(testText, htmlOptions)
	svgResult, _ := ascii.RenderOpts(testText, svgOptions)

	// They should be different
	if htmlResult == svgResult {
		t.Error("HTML and SVG output should differ")
	}

	// HTML should use <span>
	if !strings.Contains(htmlResult, "<span") {
		t.Error("HTML should use span elements for colors")
	}

	// SVG should use <tspan>
	if !strings.Contains(svgResult, "<tspan") {
		t.Error("SVG should use tspan elements for colors")
	}

	// HTML uses <code> wrapper
	if !strings.Contains(htmlResult, "<code>") {
		t.Error("HTML should have <code> wrapper")
	}

	// SVG uses <text> wrapper
	if !strings.Contains(svgResult, "<text>") {
		t.Error("SVG should have <text> wrapper")
	}
}

// TestEdgeCases tests various edge cases
func TestEdgeCases(t *testing.T) {
	ascii := NewAsciiRender()

	t.Run("single_character", func(t *testing.T) {
		options := NewRenderOptions()
		parser, _ := GetParser("svg")
		options.Parser = *parser

		result, err := ascii.RenderOpts("A", options)
		if err != nil {
			t.Fatalf("Single character render failed: %v", err)
		}
		if len(result) == 0 {
			t.Error("Result should not be empty")
		}
	})

	t.Run("multiple_spaces", func(t *testing.T) {
		options := NewRenderOptions()
		parser, _ := GetParser("svg")
		options.Parser = *parser

		result, _ := ascii.RenderOpts("A  B", options)
		// All spaces should be encoded
		spaceCount := strings.Count("A  B", " ")
		encodedCount := strings.Count(result, "&#160;")
		if encodedCount < spaceCount {
			t.Errorf("Expected at least %d encoded spaces, got %d", spaceCount, encodedCount)
		}
	})

	t.Run("color_cycling", func(t *testing.T) {
		options := NewRenderOptions()
		parser, _ := GetParser("svg")
		options.Parser = *parser
		// Only 2 colors for 4 characters - should cycle
		options.FontColor = []Color{ColorRed, ColorGreen}

		result, _ := ascii.RenderOpts("ABCD", options)
		// Should contain both colors
		if !strings.Contains(result, "<tspan") {
			t.Error("Should contain colored spans")
		}
	})
}
