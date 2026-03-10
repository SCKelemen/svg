package svg

import (
	"bytes"
	"image/color"
	"image/png"
	"math"
	"strings"
	"testing"
)

func countVisiblePixelsFromPNG(t *testing.T, pngData []byte) int {
	t.Helper()

	img, err := png.Decode(bytes.NewReader(pngData))
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	b := img.Bounds()
	count := 0
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a != 0 {
				count++
			}
		}
	}

	return count
}

func TestExportSVG(t *testing.T) {
	svgData := `<svg width="100" height="100"><rect x="10" y="10" width="80" height="80" fill="#ff0000"/></svg>`

	opts := DefaultExportOptions()
	result, err := Export(svgData, opts)

	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	if len(result) == 0 {
		t.Fatal("Export returned empty result")
	}

	// For SVG format, should return the original data
	if string(result) != svgData {
		t.Errorf("SVG export should return original data")
	}
}

func TestExportPNG(t *testing.T) {
	svgData := `<svg width="100" height="100">
		<rect x="10" y="10" width="80" height="80" fill="#ff0000"/>
	</svg>`

	opts := ExportOptions{
		Format: FormatPNG,
		Width:  100,
		Height: 100,
	}

	result, err := Export(svgData, opts)

	if err != nil {
		t.Fatalf("PNG export failed: %v", err)
	}

	if len(result) == 0 {
		t.Fatal("PNG export returned empty result")
	}

	// Check PNG signature
	if !strings.HasPrefix(string(result), "\x89PNG") {
		t.Error("Result does not have PNG signature")
	}
}

func TestExportCircle(t *testing.T) {
	svgData := `<svg width="100" height="100">
		<circle cx="50" cy="50" r="40" fill="#0000ff"/>
	</svg>`

	opts := ExportOptions{
		Format: FormatPNG,
		Width:  100,
		Height: 100,
	}

	result, err := Export(svgData, opts)

	if err != nil {
		t.Fatalf("Circle export failed: %v", err)
	}

	if len(result) == 0 {
		t.Fatal("Circle export returned empty result")
	}
}

func TestExportGroupedElements(t *testing.T) {
	svgData := `<svg width="100" height="100">
		<g>
			<rect x="10" y="10" width="30" height="30" fill="#ff0000"/>
		</g>
	</svg>`

	opts := ExportOptions{
		Format: FormatPNG,
		Width:  100,
		Height: 100,
	}

	result, err := Export(svgData, opts)
	if err != nil {
		t.Fatalf("grouped export failed: %v", err)
	}

	nonTransparent := countVisiblePixelsFromPNG(t, result)
	if nonTransparent == 0 {
		t.Fatal("expected grouped content to render, got fully white image")
	}
}

func TestExportMalformedSVGReturnsError(t *testing.T) {
	svgData := `<svg><rect></svg`

	_, err := Export(svgData, ExportOptions{
		Format: FormatPNG,
		Width:  100,
		Height: 100,
	})
	if err == nil {
		t.Fatal("expected malformed SVG to return an error")
	}
}

func TestExportUnsupportedElementsReturnErrorByDefault(t *testing.T) {
	svgData := `<svg width="100" height="100"><text x="10" y="20">hello</text></svg>`

	_, err := Export(svgData, ExportOptions{
		Format: FormatPNG,
		Width:  100,
		Height: 100,
	})
	if err == nil {
		t.Fatal("expected unsupported element error")
	}

	unsupportedErr, ok := err.(*UnsupportedElementsError)
	if !ok {
		t.Fatalf("expected UnsupportedElementsError, got %T: %v", err, err)
	}
	if len(unsupportedErr.Elements) == 0 || unsupportedErr.Elements[0] != "text" {
		t.Fatalf("expected text to be reported unsupported, got %#v", unsupportedErr.Elements)
	}
}

func TestExportUnsupportedElementsCanBeIgnored(t *testing.T) {
	svgData := `<svg width="100" height="100">
		<text x="10" y="20">hello</text>
		<rect x="0" y="0" width="10" height="10" fill="#ff0000"/>
	</svg>`

	result, err := Export(svgData, ExportOptions{
		Format:            FormatPNG,
		Width:             100,
		Height:            100,
		IgnoreUnsupported: true,
	})
	if err != nil {
		t.Fatalf("expected unsupported elements to be ignored, got error: %v", err)
	}

	if countVisiblePixelsFromPNG(t, result) == 0 {
		t.Fatal("expected supported content to still render")
	}
}

func TestExportDefsDoNotRenderContent(t *testing.T) {
	svgData := `<svg width="100" height="100">
		<defs>
			<rect x="0" y="0" width="100" height="100" fill="#ff0000"/>
		</defs>
	</svg>`

	result, err := Export(svgData, ExportOptions{
		Format: FormatPNG,
		Width:  100,
		Height: 100,
	})
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}

	if visible := countVisiblePixelsFromPNG(t, result); visible != 0 {
		t.Fatalf("expected defs content to be non-rendering, got %d visible pixels", visible)
	}
}

func TestExportUnsupportedInsideDefsDoesNotFail(t *testing.T) {
	svgData := `<svg width="100" height="100">
		<defs>
			<clipPath id="c"><path d="M0 0L10 10"/></clipPath>
		</defs>
	</svg>`

	_, err := Export(svgData, ExportOptions{
		Format: FormatPNG,
		Width:  100,
		Height: 100,
	})
	if err != nil {
		t.Fatalf("expected defs-only unsupported content to be ignored, got: %v", err)
	}
}

func TestExportLineDoesNotFloodCanvas(t *testing.T) {
	svgData := `<svg width="100" height="100">
		<line x1="10" y1="10" x2="90" y2="90" stroke="#ff0000" stroke-width="1"/>
	</svg>`

	result, err := Export(svgData, ExportOptions{
		Format: FormatPNG,
		Width:  100,
		Height: 100,
	})
	if err != nil {
		t.Fatalf("line export failed: %v", err)
	}

	nonTransparent := countVisiblePixelsFromPNG(t, result)
	if nonTransparent == 0 {
		t.Fatal("expected line to render, got fully white image")
	}
	// A 1px diagonal line in 100x100 should not fill most of the canvas.
	if nonTransparent > 1000 {
		t.Fatalf("expected line rendering to stay narrow, got %d non-transparent pixels", nonTransparent)
	}
}

func TestExportPNGPreservesTransparencyByDefault(t *testing.T) {
	svgData := `<svg width="10" height="10"></svg>`

	result, err := Export(svgData, ExportOptions{
		Format: FormatPNG,
		Width:  10,
		Height: 10,
	})
	if err != nil {
		t.Fatalf("PNG export failed: %v", err)
	}

	img, err := png.Decode(bytes.NewReader(result))
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	_, _, _, a := img.At(0, 0).RGBA()
	if a != 0 {
		t.Fatalf("expected transparent background for PNG, got alpha=%d", a)
	}
}

func TestExportPercentageLengthsRender(t *testing.T) {
	svgData := `<svg width="200" height="100">
		<rect x="25%" y="20%" width="50%" height="50%" fill="#ff0000"/>
	</svg>`

	result, err := Export(svgData, ExportOptions{
		Format: FormatPNG,
		Width:  200,
		Height: 100,
	})
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}

	visible := countVisiblePixelsFromPNG(t, result)
	if visible == 0 {
		t.Fatal("expected percentage-sized rectangle to render")
	}
	// Expected area is about 100 * 50 = 5000 px; allow small rasterization variance.
	if visible < 4500 || visible > 5500 {
		t.Fatalf("expected visible pixels around 5000, got %d", visible)
	}
}

func TestExportDPIAffectsPhysicalUnits(t *testing.T) {
	svgData := `<svg width="1in" height="1in"></svg>`

	result, err := Export(svgData, ExportOptions{
		Format: FormatPNG,
		DPI:    192,
	})
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}

	cfg, err := png.DecodeConfig(bytes.NewReader(result))
	if err != nil {
		t.Fatalf("failed to decode PNG config: %v", err)
	}
	if cfg.Width != 192 || cfg.Height != 192 {
		t.Fatalf("expected 192x192 at 192 DPI, got %dx%d", cfg.Width, cfg.Height)
	}
}

func TestParseLengthFloatSupportsUnits(t *testing.T) {
	const dpi = 96.0

	tests := []struct {
		input    string
		expected float64
	}{
		{"10.5", 10.5},
		{"1in", 96},
		{"2.54cm", 96},
		{"25.4mm", 96},
		{"72pt", 96},
		{"6pc", 96},
		{"400q", 96 / 101.6 * 400},
		{"400Q", 96 / 101.6 * 400},
		{"1e2px", 100},
	}

	for _, tt := range tests {
		got := parseLengthFloat(tt.input, dpi)
		if math.Abs(got-tt.expected) > 0.001 {
			t.Fatalf("parseLengthFloat(%q) = %f, expected %f", tt.input, got, tt.expected)
		}
	}
}

func TestParseLengthFloatWithReferenceSupportsPercentages(t *testing.T) {
	got := parseLengthFloatWithReference("37.5%", 96, 200)
	if math.Abs(got-75) > 0.001 {
		t.Fatalf("parseLengthFloatWithReference(37.5%%) = %f, expected 75", got)
	}
}

func TestParseColorUnknownReturnsTransparent(t *testing.T) {
	c := parseColor("totally-unknown-color")
	_, _, _, a := c.RGBA()
	if a != 0 {
		t.Fatalf("expected unknown colors to be transparent, got alpha=%d", a)
	}
}

func TestExportJPEG(t *testing.T) {
	svgData := `<svg width="100" height="100">
		<rect x="0" y="0" width="100" height="100" fill="#00ff00"/>
	</svg>`

	opts := ExportOptions{
		Format:  FormatJPEG,
		Width:   100,
		Height:  100,
		Quality: 90,
	}

	result, err := Export(svgData, opts)

	if err != nil {
		t.Fatalf("JPEG export failed: %v", err)
	}

	if len(result) == 0 {
		t.Fatal("JPEG export returned empty result")
	}

	// Check JPEG signature (FF D8 FF)
	if len(result) < 3 || result[0] != 0xFF || result[1] != 0xD8 || result[2] != 0xFF {
		t.Error("Result does not have JPEG signature")
	}
}

func TestParseFormat(t *testing.T) {
	tests := []struct {
		input    string
		expected ExportFormat
		hasError bool
	}{
		{"svg", FormatSVG, false},
		{"SVG", FormatSVG, false},
		{"png", FormatPNG, false},
		{"PNG", FormatPNG, false},
		{"jpeg", FormatJPEG, false},
		{"jpg", FormatJPEG, false},
		{"JPEG", FormatJPEG, false},
		{"unknown", "", true},
	}

	for _, tt := range tests {
		result, err := ParseFormat(tt.input)

		if tt.hasError {
			if err == nil {
				t.Errorf("ParseFormat(%q) expected error, got nil", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("ParseFormat(%q) unexpected error: %v", tt.input, err)
			}
			if result != tt.expected {
				t.Errorf("ParseFormat(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		}
	}
}

func TestGetMimeType(t *testing.T) {
	tests := []struct {
		format   ExportFormat
		expected string
	}{
		{FormatSVG, "image/svg+xml"},
		{FormatPNG, "image/png"},
		{FormatJPEG, "image/jpeg"},
	}

	for _, tt := range tests {
		result := GetMimeType(tt.format)
		if result != tt.expected {
			t.Errorf("GetMimeType(%q) = %q, expected %q", tt.format, result, tt.expected)
		}
	}
}

func TestGetFileExtension(t *testing.T) {
	tests := []struct {
		format   ExportFormat
		expected string
	}{
		{FormatSVG, ".svg"},
		{FormatPNG, ".png"},
		{FormatJPEG, ".jpg"},
	}

	for _, tt := range tests {
		result := GetFileExtension(tt.format)
		if result != tt.expected {
			t.Errorf("GetFileExtension(%q) = %q, expected %q", tt.format, result, tt.expected)
		}
	}
}

func TestParseColor(t *testing.T) {
	// Basic smoke test
	colors := []string{
		"#ff0000",
		"#f00",
		"red",
		"blue",
		"white",
		"black",
		"none",
	}

	for _, c := range colors {
		result := parseColor(c)
		if result == nil {
			t.Errorf("parseColor(%q) returned nil", c)
		}
	}
}

func TestParseColorAdvancedFormats(t *testing.T) {
	tests := []struct {
		input     string
		want      color.RGBA
		tolerance int
	}{
		{
			input:     "rgb(255, 0, 0)",
			want:      color.RGBA{R: 255, G: 0, B: 0, A: 255},
			tolerance: 0,
		},
		{
			input:     "hsl(120, 100%, 50%)",
			want:      color.RGBA{R: 0, G: 255, B: 0, A: 255},
			tolerance: 1,
		},
		{
			input:     "#336699cc",
			want:      color.RGBA{R: 51, G: 102, B: 153, A: 204},
			tolerance: 0,
		},
		{
			input:     "rgb(255 0 0 / 50%)",
			want:      color.RGBA{R: 255, G: 0, B: 0, A: 128},
			tolerance: 1,
		},
	}

	for _, tt := range tests {
		got := color.RGBAModel.Convert(parseColor(tt.input)).(color.RGBA)
		if colorDistance(got, tt.want) > tt.tolerance {
			t.Fatalf("parseColor(%q) = %#v, expected %#v", tt.input, got, tt.want)
		}
	}
}

func colorDistance(a, b color.RGBA) int {
	dr := int(a.R) - int(b.R)
	if dr < 0 {
		dr = -dr
	}
	dg := int(a.G) - int(b.G)
	if dg < 0 {
		dg = -dg
	}
	db := int(a.B) - int(b.B)
	if db < 0 {
		db = -db
	}
	da := int(a.A) - int(b.A)
	if da < 0 {
		da = -da
	}

	max := dr
	if dg > max {
		max = dg
	}
	if db > max {
		max = db
	}
	if da > max {
		max = da
	}
	return max
}
