package figlet4go

import (
	"fmt"
	"os"
)

// Example demonstrates basic SVG output
func Example_svgBasic() {
	ascii := NewAsciiRender()

	options := NewRenderOptions()
	parser, _ := GetParser("svg")
	options.Parser = *parser

	result, _ := ascii.RenderOpts("Hello", options)
	fmt.Print(result)
}

// Example demonstrates SVG with colors
func Example_svgWithColors() {
	ascii := NewAsciiRender()

	options := NewRenderOptions()
	parser, _ := GetParser("svg")
	options.Parser = *parser
	options.FontColor = []Color{
		ColorRed,
		ColorGreen,
		ColorBlue,
	}

	result, _ := ascii.RenderOpts("RGB", options)
	fmt.Print(result)
}

// Example demonstrates creating a complete SVG file
func Example_svgCompleteDocument() {
	ascii := NewAsciiRender()

	options := NewRenderOptions()
	parser, _ := GetParser("svg")
	options.Parser = *parser
	options.FontName = "standard"

	red, _ := NewTrueColorFromHexString("FF6B6B")
	teal, _ := NewTrueColorFromHexString("4ECDC4")
	yellow, _ := NewTrueColorFromHexString("FFE66D")

	options.FontColor = []Color{
		red,
		teal,
		yellow,
	}

	result, _ := ascii.RenderOpts("SVG", options)

	// Create a complete SVG document
	svgDoc := createSVGDocument(result, 800, 200)

	// Write to file
	os.WriteFile("output.svg", []byte(svgDoc), 0644)
	fmt.Println("SVG written to output.svg")
}

// Helper function to create a complete SVG document
func createSVGDocument(content string, width, height int) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">
  <rect width="100%%" height="100%%" fill="white"/>
  <g font-family="monospace" font-size="14">
    %s
  </g>
</svg>`, width, height, width, height, content)
}

// Example demonstrates different fonts with SVG
func Example_svgDifferentFonts() {
	ascii := NewAsciiRender()

	fonts := []string{"standard", "larry3d"}

	for _, fontName := range fonts {
		options := NewRenderOptions()
		parser, _ := GetParser("svg")
		options.Parser = *parser
		options.FontName = fontName

		result, err := ascii.RenderOpts("Test", options)
		if err == nil {
			fmt.Printf("Font: %s\n%s\n", fontName, result)
		}
	}
}
