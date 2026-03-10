package svg

import (
	"bytes"
	"image/png"
	"strings"
	"testing"
)

func countNonWhitePixelsFromPNG(t *testing.T, pngData []byte) int {
	t.Helper()

	img, err := png.Decode(bytes.NewReader(pngData))
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	b := img.Bounds()
	count := 0
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			if !(r == 0xffff && g == 0xffff && b == 0xffff && a == 0xffff) {
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

	nonWhite := countNonWhitePixelsFromPNG(t, result)
	if nonWhite == 0 {
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

	nonWhite := countNonWhitePixelsFromPNG(t, result)
	if nonWhite == 0 {
		t.Fatal("expected line to render, got fully white image")
	}
	// A 1px diagonal line in 100x100 should not fill most of the canvas.
	if nonWhite > 1000 {
		t.Fatalf("expected line rendering to stay narrow, got %d non-white pixels", nonWhite)
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
