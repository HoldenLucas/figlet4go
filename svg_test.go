package figlet4go

import (
	"fmt"
	"strings"
	"testing"
)

// TestSVGDocumentGeneration shows how to create a complete SVG document
func TestSVGDocumentGeneration(t *testing.T) {
	ascii := NewAsciiRender()

	options := NewRenderOptions()
	parser, err := GetParser("svg")
	if err != nil {
		t.Fatalf("Failed to get SVG parser: %v", err)
	}
	options.Parser = *parser
	options.FontName = "standard"

	result, err := ascii.RenderOpts("SVG", options)
	if err != nil {
		t.Fatalf("RenderOpts() error = %v", err)
	}

	// Wrap in a complete SVG document
	svgDoc := WrapSVG(result, SVGOptions{
		Width:      800,
		Height:     200,
		FontFamily: "monospace",
		FontSize:   14,
	})

	// Basic validation
	if !strings.HasPrefix(svgDoc, "<?xml") && !strings.HasPrefix(svgDoc, "<svg") {
		t.Error("SVG document should start with XML declaration or svg tag")
	}
	if !strings.Contains(svgDoc, "<svg") {
		t.Error("SVG document should contain <svg> tag")
	}
	if !strings.Contains(svgDoc, "</svg>") {
		t.Error("SVG document should contain </svg> closing tag")
	}
	// Check for text elements instead of exact result string
	// since WrapSVG transforms the content
	if !strings.Contains(svgDoc, "<text") {
		t.Error("SVG document should contain <text> elements")
	}
}

func TestSVGWithColors(t *testing.T) {
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
		t.Fatalf("RenderOpts() error = %v", err)
	}

	// Should contain SVG tspan elements with fill attribute
	if !strings.Contains(result, "<tspan fill=") {
		t.Error("Colored SVG should contain <tspan fill=> elements")
	}
	if !strings.Contains(result, "</tspan>") {
		t.Error("Colored SVG should contain </tspan> closing tags")
	}
	if !strings.Contains(result, "rgb(") {
		t.Error("SVG colors should use rgb() format")
	}
}

func TestSVGTrueColor(t *testing.T) {
	ascii := NewAsciiRender()

	options := NewRenderOptions()
	parser, err := GetParser("svg")
	if err != nil {
		t.Fatalf("Failed to get SVG parser: %v", err)
	}
	options.Parser = *parser
	options.FontName = "standard"

	// Create a custom TrueColor
	customColor := TrueColor{r: 255, g: 128, b: 0} // Orange
	options.FontColor = []Color{customColor}

	result, err := ascii.RenderOpts("OK", options)
	if err != nil {
		t.Fatalf("RenderOpts() error = %v", err)
	}

	// Should contain the specific RGB values
	if !strings.Contains(result, "rgb(255,128,0)") {
		t.Error("SVG should contain the custom RGB color rgb(255,128,0)")
	}
}

// SVGOptions for wrapping rendered text in a complete SVG document
type SVGOptions struct {
	Width      int
	Height     int
	FontFamily string
	FontSize   int
}

// WrapSVG wraps the rendered ASCII art in a complete SVG document
func WrapSVG(content string, opts SVGOptions) string {
	// Set defaults
	if opts.Width == 0 {
		opts.Width = 800
	}
	if opts.Height == 0 {
		opts.Height = 200
	}
	if opts.FontFamily == "" {
		opts.FontFamily = "monospace"
	}
	if opts.FontSize == 0 {
		opts.FontSize = 14
	}

	// Calculate approximate line height
	lineHeight := opts.FontSize + 2

	// Split content by newlines to position each line
	lines := strings.Split(content, "\n")

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`,
		opts.Width, opts.Height, opts.Width, opts.Height))
	builder.WriteString("\n  <style>text { font-family: ")
	builder.WriteString(opts.FontFamily)
	builder.WriteString("; font-size: ")
	builder.WriteString(fmt.Sprintf("%dpx", opts.FontSize))
	builder.WriteString("; white-space: pre; }</style>\n")

	// Position each line
	for i, line := range lines {
		if line == "" || line == "<text>" || line == "</text>" {
			continue
		}
		// Remove the <text> and </text> wrappers, keep the content
		line = strings.TrimPrefix(line, "<text>")
		line = strings.TrimSuffix(line, "</text>")
		line = strings.ReplaceAll(line, "<br/>", "")

		if line != "" {
			y := (i + 1) * lineHeight
			builder.WriteString(fmt.Sprintf(`  <text x="10" y="%d">%s</text>`, y, line))
			builder.WriteString("\n")
		}
	}

	builder.WriteString("</svg>")
	return builder.String()
}

// Example test showing usage
func Example_svgGeneration() {
	ascii := NewAsciiRender()

	options := NewRenderOptions()
	parser, _ := GetParser("svg")
	options.Parser = *parser
	options.FontName = "standard"

	result, _ := ascii.RenderOpts("Test", options)

	svgDoc := WrapSVG(result, SVGOptions{
		Width:      600,
		Height:     150,
		FontFamily: "Courier New, monospace",
		FontSize:   12,
	})

	fmt.Println(svgDoc)
}
