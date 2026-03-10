package svg

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"strconv"
	"strings"

	"golang.org/x/image/vector"
)

// ExportFormat represents an export format
type ExportFormat string

const (
	// FormatSVG exports as SVG (passthrough)
	FormatSVG ExportFormat = "svg"
	// FormatPNG exports as PNG
	FormatPNG ExportFormat = "png"
	// FormatJPEG exports as JPEG
	FormatJPEG ExportFormat = "jpeg"
)

const defaultRasterDPI = 96.0

// ExportOptions configures export settings
type ExportOptions struct {
	Format  ExportFormat
	Width   int // For raster formats, 0 = use SVG dimensions
	Height  int // For raster formats, 0 = use SVG dimensions
	Quality int // For JPEG, 0-100 (default 90)
	DPI     int // Dots per inch for physical SVG units like in/cm/mm/pt (default 96)
}

// DefaultExportOptions returns sensible defaults
func DefaultExportOptions() ExportOptions {
	return ExportOptions{
		Format:  FormatSVG,
		Quality: 90,
		DPI:     96,
	}
}

// Export converts SVG to the specified format
func Export(svgData string, opts ExportOptions) ([]byte, error) {
	// For SVG, just return the data
	if opts.Format == FormatSVG {
		return []byte(svgData), nil
	}

	// For raster formats, parse and rasterize
	return rasterize(svgData, opts)
}

// svgElement represents a parsed SVG element
type svgElement struct {
	Tag        string
	Attributes map[string]string
	Children   []*svgElement
	Text       string
}

// parseSVG performs basic SVG parsing for our own generated SVG
func parseSVG(svgData string) (*svgElement, error) {
	decoder := xml.NewDecoder(strings.NewReader(svgData))

	var root *svgElement
	var stack []*svgElement

	for {
		token, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("SVG parse error: %w", err)
		}

		switch t := token.(type) {
		case xml.StartElement:
			elem := &svgElement{
				Tag:        t.Name.Local,
				Attributes: make(map[string]string),
			}

			for _, attr := range t.Attr {
				elem.Attributes[attr.Name.Local] = attr.Value
			}

			if len(stack) > 0 {
				parent := stack[len(stack)-1]
				parent.Children = append(parent.Children, elem)
			} else {
				root = elem
			}

			stack = append(stack, elem)

		case xml.EndElement:
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}

		case xml.CharData:
			if len(stack) > 0 {
				text := strings.TrimSpace(string(t))
				if text != "" {
					stack[len(stack)-1].Text = text
				}
			}
		}
	}

	if root == nil {
		return nil, fmt.Errorf("no SVG root element found")
	}

	return root, nil
}

// rasterize converts SVG to a raster image
func rasterize(svgData string, opts ExportOptions) ([]byte, error) {
	dpi := resolveDPI(opts)

	// Parse SVG
	root, err := parseSVG(svgData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SVG: %w", err)
	}

	// Get dimensions
	width, height, err := getSVGDimensions(root, opts, dpi)
	if err != nil {
		return nil, fmt.Errorf("failed to get SVG dimensions: %w", err)
	}

	// Create image
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Preserve transparency for PNG; JPEG always needs an opaque background.
	if opts.Format == FormatJPEG {
		draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	}

	// Create rasterizer
	rasterizer := vector.NewRasterizer(width, height)

	// Render SVG elements
	if err := renderElement(root, img, rasterizer, width, height, dpi); err != nil {
		return nil, fmt.Errorf("failed to render SVG: %w", err)
	}

	// Encode to target format
	var buf bytes.Buffer
	switch opts.Format {
	case FormatPNG:
		encoder := png.Encoder{CompressionLevel: png.DefaultCompression}
		if err := encoder.Encode(&buf, img); err != nil {
			return nil, fmt.Errorf("failed to encode PNG: %w", err)
		}
	case FormatJPEG:
		quality := opts.Quality
		if quality == 0 {
			quality = 90
		}
		if quality < 1 {
			quality = 1
		}
		if quality > 100 {
			quality = 100
		}
		jpegOpts := &jpeg.Options{Quality: quality}
		if err := jpeg.Encode(&buf, img, jpegOpts); err != nil {
			return nil, fmt.Errorf("failed to encode JPEG: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported format: %s", opts.Format)
	}

	return buf.Bytes(), nil
}

// getSVGDimensions extracts width and height from SVG
func getSVGDimensions(root *svgElement, opts ExportOptions, dpi float64) (int, int, error) {
	width := opts.Width
	height := opts.Height

	// Try to get from attributes
	if width == 0 {
		if w, ok := root.Attributes["width"]; ok {
			width = parseLength(w, dpi)
		}
	}

	if height == 0 {
		if h, ok := root.Attributes["height"]; ok {
			height = parseLength(h, dpi)
		}
	}

	// Try viewBox if dimensions not set
	if width == 0 || height == 0 {
		if viewBox, ok := root.Attributes["viewBox"]; ok {
			parts := strings.Fields(viewBox)
			if len(parts) == 4 {
				if width == 0 {
					width = parseLength(parts[2], dpi)
				}
				if height == 0 {
					height = parseLength(parts[3], dpi)
				}
			}
		}
	}

	// Default dimensions
	if width == 0 {
		width = 800
	}
	if height == 0 {
		height = 600
	}

	return width, height, nil
}

func resolveDPI(opts ExportOptions) float64 {
	if opts.DPI <= 0 {
		return defaultRasterDPI
	}
	return float64(opts.DPI)
}

// parseLength parses a length string into device pixels.
func parseLength(s string, dpi float64) int {
	return int(math.Round(parseLengthFloat(s, dpi)))
}

// parseLengthFloat parses a length string and returns a pixel value.
func parseLengthFloat(s string, dpi float64) float64 {
	s = strings.ToLower(strings.TrimSpace(s))
	if s == "" {
		return 0
	}
	// Percentages require viewport/object context and are not currently supported.
	if strings.HasSuffix(s, "%") {
		return 0
	}

	units := []struct {
		suffix string
		factor float64
	}{
		{"px", 1},
		{"pt", dpi / 72.0},
		{"pc", dpi / 6.0},
		{"in", dpi},
		{"cm", dpi / 2.54},
		{"mm", dpi / 25.4},
		{"q", dpi / 101.6},
	}

	for _, unit := range units {
		if strings.HasSuffix(s, unit.suffix) {
			num := strings.TrimSpace(strings.TrimSuffix(s, unit.suffix))
			val, err := strconv.ParseFloat(num, 64)
			if err != nil {
				return 0
			}
			return val * unit.factor
		}
	}

	if val, err := strconv.ParseFloat(s, 64); err == nil {
		return val
	}

	return 0
}

// renderElement renders an SVG element to the image
func renderElement(elem *svgElement, img *image.RGBA, rasterizer *vector.Rasterizer, width, height int, dpi float64) error {
	switch elem.Tag {
	case "svg":
		// Render children
		for _, child := range elem.Children {
			if err := renderElement(child, img, rasterizer, width, height, dpi); err != nil {
				return err
			}
		}

	case "rect":
		return renderRect(elem, img, rasterizer, dpi)

	case "circle":
		return renderCircle(elem, img, rasterizer, dpi)

	case "line":
		return renderLine(elem, img, dpi)

	case "g":
		// Group - render children
		for _, child := range elem.Children {
			if err := renderElement(child, img, rasterizer, width, height, dpi); err != nil {
				return err
			}
		}

	case "text":
		// Text rendering is complex and requires font support
		// For now, we'll skip text elements
		// TODO: Implement text rendering using golang.org/x/image/font

	case "path":
		// Path rendering is complex
		// TODO: Implement path rendering

	default:
		// Unknown or unsupported element, continue rendering children
		for _, child := range elem.Children {
			if err := renderElement(child, img, rasterizer, width, height, dpi); err != nil {
				return err
			}
		}
	}

	return nil
}

// renderRect renders a rectangle
func renderRect(elem *svgElement, img *image.RGBA, rasterizer *vector.Rasterizer, dpi float64) error {
	_ = rasterizer
	x := parseLengthFloat(elem.Attributes["x"], dpi)
	y := parseLengthFloat(elem.Attributes["y"], dpi)
	w := parseLengthFloat(elem.Attributes["width"], dpi)
	h := parseLengthFloat(elem.Attributes["height"], dpi)

	fillColor := parseColor(elem.Attributes["fill"])
	if isTransparent(fillColor) || w <= 0 || h <= 0 {
		return nil
	}

	left := int(math.Floor(x))
	top := int(math.Floor(y))
	right := int(math.Ceil(x + w))
	bottom := int(math.Ceil(y + h))

	// Draw rectangle
	rect := image.Rect(left, top, right, bottom)
	draw.Draw(img, rect, &image.Uniform{fillColor}, image.Point{}, draw.Over)

	return nil
}

// renderCircle renders a circle
func renderCircle(elem *svgElement, img *image.RGBA, rasterizer *vector.Rasterizer, dpi float64) error {
	cx := parseLengthFloat(elem.Attributes["cx"], dpi)
	cy := parseLengthFloat(elem.Attributes["cy"], dpi)
	r := parseLengthFloat(elem.Attributes["r"], dpi)

	fillColor := parseColor(elem.Attributes["fill"])
	if isTransparent(fillColor) || r <= 0 {
		return nil
	}

	// Use vector rasterizer for smooth circles
	rasterizer.Reset(img.Bounds().Dx(), img.Bounds().Dy())
	rasterizer.DrawOp = draw.Over

	// Draw circle using arc approximation
	drawCircle(rasterizer, float32(cx), float32(cy), float32(r))

	// Rasterize
	src := image.NewUniform(fillColor)
	rasterizer.Draw(img, img.Bounds(), src, image.Point{})

	return nil
}

// renderLine renders a line
func renderLine(elem *svgElement, img *image.RGBA, dpi float64) error {
	x1 := parseLengthFloat(elem.Attributes["x1"], dpi)
	y1 := parseLengthFloat(elem.Attributes["y1"], dpi)
	x2 := parseLengthFloat(elem.Attributes["x2"], dpi)
	y2 := parseLengthFloat(elem.Attributes["y2"], dpi)

	strokeColor := parseColor(elem.Attributes["stroke"])
	strokeWidth := parseLengthFloat(elem.Attributes["stroke-width"], dpi)
	if strokeWidth <= 0 {
		strokeWidth = 1
	}
	if isTransparent(strokeColor) {
		return nil
	}

	drawThickLine(img, x1, y1, x2, y2, strokeWidth, strokeColor)

	return nil
}

func isTransparent(c color.Color) bool {
	_, _, _, a := c.RGBA()
	return a == 0
}

func drawThickLine(img *image.RGBA, x1, y1, x2, y2, width float64, c color.Color) {
	dx := x2 - x1
	dy := y2 - y1
	steps := int(math.Max(math.Abs(dx), math.Abs(dy)))

	brush := int(math.Round(width))
	if brush < 1 {
		brush = 1
	}
	half := brush / 2

	if steps == 0 {
		drawBrush(img, int(math.Round(x1)), int(math.Round(y1)), brush, half, c)
		return
	}

	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		x := x1 + dx*t
		y := y1 + dy*t
		drawBrush(img, int(math.Round(x)), int(math.Round(y)), brush, half, c)
	}
}

func drawBrush(img *image.RGBA, cx, cy, brush, half int, c color.Color) {
	for oy := -half; oy < brush-half; oy++ {
		for ox := -half; ox < brush-half; ox++ {
			x := cx + ox
			y := cy + oy
			if !image.Pt(x, y).In(img.Bounds()) {
				continue
			}
			draw.Draw(img, image.Rect(x, y, x+1, y+1), &image.Uniform{C: c}, image.Point{}, draw.Over)
		}
	}
}

// parseColor parses a color string (hex or named)
func parseColor(s string) color.Color {
	s = strings.TrimSpace(s)

	if s == "" || s == "none" {
		return color.Transparent
	}

	// Handle hex colors
	if strings.HasPrefix(s, "#") {
		s = strings.TrimPrefix(s, "#")

		var r, g, b uint8

		if len(s) == 6 {
			// #RRGGBB
			if val, err := strconv.ParseUint(s[0:2], 16, 8); err == nil {
				r = uint8(val)
			}
			if val, err := strconv.ParseUint(s[2:4], 16, 8); err == nil {
				g = uint8(val)
			}
			if val, err := strconv.ParseUint(s[4:6], 16, 8); err == nil {
				b = uint8(val)
			}
		} else if len(s) == 3 {
			// #RGB (shorthand)
			if val, err := strconv.ParseUint(s[0:1], 16, 8); err == nil {
				r = uint8(val * 17) // 0xF -> 0xFF
			}
			if val, err := strconv.ParseUint(s[1:2], 16, 8); err == nil {
				g = uint8(val * 17)
			}
			if val, err := strconv.ParseUint(s[2:3], 16, 8); err == nil {
				b = uint8(val * 17)
			}
		}

		return color.RGBA{R: r, G: g, B: b, A: 255}
	}

	// Handle named colors
	switch s {
	case "white":
		return color.White
	case "black":
		return color.Black
	case "red":
		return color.RGBA{R: 255, A: 255}
	case "green":
		return color.RGBA{G: 255, A: 255}
	case "blue":
		return color.RGBA{B: 255, A: 255}
	default:
		return color.Black
	}
}

// drawCircle draws a circle using the vector rasterizer
func drawCircle(r *vector.Rasterizer, cx, cy, radius float32) {
	const segments = 32

	r.MoveTo(cx+radius, cy)

	for i := 1; i <= segments; i++ {
		angle := float64(i) * 2.0 * math.Pi / float64(segments)
		x := cx + radius*float32(math.Cos(angle))
		y := cy + radius*float32(math.Sin(angle))
		r.LineTo(x, y)
	}

	r.ClosePath()
}

// GetMimeType returns the MIME type for a format
func GetMimeType(format ExportFormat) string {
	switch format {
	case FormatSVG:
		return "image/svg+xml"
	case FormatPNG:
		return "image/png"
	case FormatJPEG:
		return "image/jpeg"
	default:
		return "application/octet-stream"
	}
}

// GetFileExtension returns the file extension for a format
func GetFileExtension(format ExportFormat) string {
	switch format {
	case FormatSVG:
		return ".svg"
	case FormatPNG:
		return ".png"
	case FormatJPEG:
		return ".jpg"
	default:
		return ".bin"
	}
}

// ParseFormat parses a format string
func ParseFormat(s string) (ExportFormat, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	switch s {
	case "svg":
		return FormatSVG, nil
	case "png":
		return FormatPNG, nil
	case "jpeg", "jpg":
		return FormatJPEG, nil
	default:
		return "", fmt.Errorf("unknown format: %s", s)
	}
}
