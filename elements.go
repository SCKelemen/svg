package rendersvg

import (
	"fmt"
	"strings"

	"github.com/SCKelemen/layout"
	"github.com/SCKelemen/units"
)

// StrokeLinecap defines how the end of a stroke is rendered
type StrokeLinecap string

const (
	StrokeLinecapButt   StrokeLinecap = "butt"
	StrokeLinecapRound  StrokeLinecap = "round"
	StrokeLinecapSquare StrokeLinecap = "square"
)

// StrokeLinejoin defines how corners of a stroke are rendered
type StrokeLinejoin string

const (
	StrokeLinejoinMiter StrokeLinejoin = "miter"
	StrokeLinejoinRound StrokeLinejoin = "round"
	StrokeLinejoinBevel StrokeLinejoin = "bevel"
)

// TextAnchor defines horizontal text alignment
type TextAnchor string

const (
	TextAnchorStart  TextAnchor = "start"
	TextAnchorMiddle TextAnchor = "middle"
	TextAnchorEnd    TextAnchor = "end"
)

// DominantBaseline defines vertical text alignment
type DominantBaseline string

const (
	DominantBaselineAuto         DominantBaseline = "auto"
	DominantBaselineMiddle       DominantBaseline = "middle"
	DominantBaselineHanging      DominantBaseline = "hanging"
	DominantBaselineTextTop      DominantBaseline = "text-top"
	DominantBaselineTextBottom   DominantBaseline = "text-bottom"
	DominantBaselineAlphabetic   DominantBaseline = "alphabetic"
	DominantBaselineMathematical DominantBaseline = "mathematical"
)

// FontWeight defines font weight values
type FontWeight string

const (
	FontWeightNormal  FontWeight = "normal"
	FontWeightBold    FontWeight = "bold"
	FontWeightBolder  FontWeight = "bolder"
	FontWeightLighter FontWeight = "lighter"
	FontWeight100     FontWeight = "100"
	FontWeight200     FontWeight = "200"
	FontWeight300     FontWeight = "300"
	FontWeight400     FontWeight = "400"
	FontWeight500     FontWeight = "500"
	FontWeight600     FontWeight = "600"
	FontWeight700     FontWeight = "700"
	FontWeight800     FontWeight = "800"
	FontWeight900     FontWeight = "900"
)

// FontStyle defines font style values
type FontStyle string

const (
	FontStyleNormal  FontStyle = "normal"
	FontStyleItalic  FontStyle = "italic"
	FontStyleOblique FontStyle = "oblique"
)

// Style represents styling attributes for SVG elements
type Style struct {
	Fill             string
	Stroke           string
	StrokeWidth      float64
	StrokeLinecap    StrokeLinecap
	StrokeLinejoin   StrokeLinejoin
	Opacity          float64
	FillOpacity      float64
	StrokeOpacity    float64
	Class            string
	ClipPath         string
	TextAnchor       TextAnchor
	DominantBaseline DominantBaseline
	FontFamily       string
	FontSize         units.Length // Type-safe CSS length with units
	FontWeight       FontWeight
	FontStyle        FontStyle
}

// Rect renders an SVG rectangle
func Rect(x, y, width, height float64, style Style) string {
	attrs := formatStyle(style)
	return fmt.Sprintf(`<rect x="%.2f" y="%.2f" width="%.2f" height="%.2f"%s/>`,
		x, y, width, height, attrs)
}

// RoundedRect renders an SVG rectangle with rounded corners
func RoundedRect(x, y, width, height, rx, ry float64, style Style) string {
	attrs := formatStyle(style)
	if ry == 0 {
		ry = rx // If ry not specified, use rx for both
	}
	return fmt.Sprintf(`<rect x="%.2f" y="%.2f" width="%.2f" height="%.2f" rx="%.2f" ry="%.2f"%s/>`,
		x, y, width, height, rx, ry, attrs)
}

// Circle renders an SVG circle
func Circle(cx, cy, r float64, style Style) string {
	attrs := formatStyle(style)
	return fmt.Sprintf(`<circle cx="%.2f" cy="%.2f" r="%.2f"%s/>`,
		cx, cy, r, attrs)
}

// Line renders an SVG line
func Line(x1, y1, x2, y2 float64, style Style) string {
	attrs := formatStyle(style)
	return fmt.Sprintf(`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f"%s/>`,
		x1, y1, x2, y2, attrs)
}

// Text renders an SVG text element
func Text(content string, x, y float64, style Style) string {
	attrs := formatStyle(style)
	return fmt.Sprintf(`<text x="%.2f" y="%.2f"%s>%s</text>`,
		x, y, attrs, escapeXML(content))
}

// TSpan renders an SVG tspan element (for use inside text elements)
// TSpan allows styling different parts of text independently
func TSpan(content string, style Style, dx, dy float64) string {
	attrs := formatStyle(style)
	posAttrs := ""
	if dx != 0 {
		posAttrs += fmt.Sprintf(` dx="%.2f"`, dx)
	}
	if dy != 0 {
		posAttrs += fmt.Sprintf(` dy="%.2f"`, dy)
	}
	return fmt.Sprintf(`<tspan%s%s>%s</tspan>`, posAttrs, attrs, escapeXML(content))
}

// TextWithSpans renders an SVG text element with multiple styled spans
func TextWithSpans(x, y float64, style Style, spans []string) string {
	attrs := formatStyle(style)
	return fmt.Sprintf(`<text x="%.2f" y="%.2f"%s>%s</text>`,
		x, y, attrs, strings.Join(spans, ""))
}

// TextPath renders text along a path
func TextPath(content string, pathID string, style Style, startOffset string) string {
	attrs := formatStyle(style)
	offsetAttr := ""
	if startOffset != "" {
		offsetAttr = fmt.Sprintf(` startOffset="%s"`, startOffset)
	}
	return fmt.Sprintf(`<textPath href="#%s"%s%s>%s</textPath>`,
		pathID, offsetAttr, attrs, escapeXML(content))
}

// Path renders an SVG path
func Path(d string, style Style) string {
	attrs := formatStyle(style)
	return fmt.Sprintf(`<path d="%s"%s/>`, d, attrs)
}

// Group wraps content in an SVG <g> element with optional transform
func Group(content string, transform string, style Style) string {
	var attrs string
	if transform != "" {
		attrs = fmt.Sprintf(` transform="%s"`, transform)
	}
	attrs += formatStyle(style)

	return fmt.Sprintf(`<g%s>%s</g>`, attrs, content)
}

// GroupWithClipPath wraps content in an SVG <g> element with clipPath
func GroupWithClipPath(content string, clipPathID string, style Style) string {
	style.ClipPath = URL(clipPathID)
	return Group(content, "", style)
}

// formatStyle converts a Style struct to SVG attribute string
func formatStyle(s Style) string {
	var attrs []string

	if s.Fill != "" {
		attrs = append(attrs, fmt.Sprintf(`fill="%s"`, s.Fill))
	}
	if s.Stroke != "" {
		attrs = append(attrs, fmt.Sprintf(`stroke="%s"`, s.Stroke))
	}
	if s.StrokeWidth > 0 {
		attrs = append(attrs, fmt.Sprintf(`stroke-width="%.2f"`, s.StrokeWidth))
	}
	if s.StrokeLinecap != "" {
		attrs = append(attrs, fmt.Sprintf(`stroke-linecap="%s"`, string(s.StrokeLinecap)))
	}
	if s.StrokeLinejoin != "" {
		attrs = append(attrs, fmt.Sprintf(`stroke-linejoin="%s"`, string(s.StrokeLinejoin)))
	}
	if s.Opacity > 0 && s.Opacity < 1 {
		attrs = append(attrs, fmt.Sprintf(`opacity="%.2f"`, s.Opacity))
	}
	if s.FillOpacity > 0 && s.FillOpacity < 1 {
		attrs = append(attrs, fmt.Sprintf(`fill-opacity="%.2f"`, s.FillOpacity))
	}
	if s.StrokeOpacity > 0 && s.StrokeOpacity < 1 {
		attrs = append(attrs, fmt.Sprintf(`stroke-opacity="%.2f"`, s.StrokeOpacity))
	}
	if s.Class != "" {
		attrs = append(attrs, fmt.Sprintf(`class="%s"`, s.Class))
	}
	if s.ClipPath != "" {
		attrs = append(attrs, fmt.Sprintf(`clip-path="%s"`, s.ClipPath))
	}
	if s.TextAnchor != "" {
		attrs = append(attrs, fmt.Sprintf(`text-anchor="%s"`, string(s.TextAnchor)))
	}
	if s.DominantBaseline != "" {
		attrs = append(attrs, fmt.Sprintf(`dominant-baseline="%s"`, string(s.DominantBaseline)))
	}
	if s.FontFamily != "" {
		attrs = append(attrs, fmt.Sprintf(`font-family="%s"`, s.FontFamily))
	}
	if s.FontSize.Value != 0 {
		// Format as "valueunit" (e.g., "16px", "1.5em", "2rem")
		attrs = append(attrs, fmt.Sprintf(`font-size="%s"`, s.FontSize.String()))
	}
	if s.FontWeight != "" {
		attrs = append(attrs, fmt.Sprintf(`font-weight="%s"`, string(s.FontWeight)))
	}
	if s.FontStyle != "" {
		attrs = append(attrs, fmt.Sprintf(`font-style="%s"`, string(s.FontStyle)))
	}

	if len(attrs) == 0 {
		return ""
	}

	return " " + strings.Join(attrs, " ")
}

// escapeXML escapes special XML characters in text content
func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

// GetTransformFromNode extracts SVG transform attribute from a layout node
func GetTransformFromNode(node *layout.Node) string {
	t := node.Style.Transform

	// Check for identity transform (no transformation needed)
	if t.IsIdentity() {
		return ""
	}

	// Check for zero transform (uninitialized) - treat as identity
	if t.A == 0 && t.B == 0 && t.C == 0 && t.D == 0 && t.E == 0 && t.F == 0 {
		return ""
	}

	return t.ToSVGString()
}

// GetRectFromNode extracts the computed rectangle from a layout node
func GetRectFromNode(node *layout.Node) layout.Rect {
	return node.Rect
}
