package svg

import (
	"fmt"
	"strings"
)

// MarkerUnits defines the coordinate system for marker dimensions
type MarkerUnits string

const (
	MarkerUnitsStrokeWidth    MarkerUnits = "strokeWidth"
	MarkerUnitsUserSpaceOnUse MarkerUnits = "userSpaceOnUse"
)

// MarkerOrient defines the orientation of a marker
type MarkerOrient string

const (
	MarkerOrientAuto      MarkerOrient = "auto"
	MarkerOrientAutoStart MarkerOrient = "auto-start-reverse"
)

// MarkerDef represents a marker definition
type MarkerDef struct {
	ID           string
	ViewBox      string       // e.g., "0 0 10 10"
	RefX         float64      // Reference point X
	RefY         float64      // Reference point Y
	MarkerWidth  float64      // Width of marker viewport
	MarkerHeight float64      // Height of marker viewport
	Orient       MarkerOrient // auto, auto-start-reverse, or angle
	MarkerUnits  MarkerUnits
	Content      string // SVG content inside the marker
}

// Marker creates a marker definition (for use in <defs>)
func Marker(def MarkerDef) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf(`<marker id="%s"`, escapeAttr(def.ID)))

	if def.ViewBox != "" {
		b.WriteString(fmt.Sprintf(` viewBox="%s"`, escapeAttr(def.ViewBox)))
	}

	b.WriteString(fmt.Sprintf(` refX="%.2f" refY="%.2f"`, def.RefX, def.RefY))

	if def.MarkerWidth > 0 {
		b.WriteString(fmt.Sprintf(` markerWidth="%.2f"`, def.MarkerWidth))
	}
	if def.MarkerHeight > 0 {
		b.WriteString(fmt.Sprintf(` markerHeight="%.2f"`, def.MarkerHeight))
	}

	if def.Orient != "" {
		b.WriteString(fmt.Sprintf(` orient="%s"`, escapeAttr(string(def.Orient))))
	}

	if def.MarkerUnits != "" {
		b.WriteString(fmt.Sprintf(` markerUnits="%s"`, escapeAttr(string(def.MarkerUnits))))
	}

	b.WriteString(">")
	b.WriteString("\n")
	b.WriteString(def.Content)
	b.WriteString("\n")
	b.WriteString("</marker>")

	return b.String()
}

// MarkerURL creates a url() reference to a marker
func MarkerURL(id string) string {
	return fmt.Sprintf("url(#%s)", id)
}

// MarkerStyle represents marker styling that can be added to a Style struct
type MarkerStyle struct {
	MarkerStart string // URL reference to marker for start of path
	MarkerMid   string // URL reference to marker for middle vertices
	MarkerEnd   string // URL reference to marker for end of path
}

// StyleWithMarkers extends a Style with marker references
func StyleWithMarkers(style Style, markers MarkerStyle) Style {
	style.MarkerStart = markers.MarkerStart
	style.MarkerMid = markers.MarkerMid
	style.MarkerEnd = markers.MarkerEnd
	return style
}

// Common marker shapes

// ArrowMarker creates a simple arrow marker pointing right
func ArrowMarker(id string, color string) string {
	content := fmt.Sprintf(`<path d="M 0 0 L 10 5 L 0 10 Z" fill="%s"/>`, escapeAttr(color))
	return Marker(MarkerDef{
		ID:           id,
		ViewBox:      "0 0 10 10",
		RefX:         10,
		RefY:         5,
		MarkerWidth:  6,
		MarkerHeight: 6,
		Orient:       MarkerOrientAuto,
		Content:      content,
	})
}

// CircleMarker creates a circular marker
func CircleMarker(id string, color string) string {
	content := fmt.Sprintf(`<circle cx="5" cy="5" r="4" fill="%s"/>`, escapeAttr(color))
	return Marker(MarkerDef{
		ID:           id,
		ViewBox:      "0 0 10 10",
		RefX:         5,
		RefY:         5,
		MarkerWidth:  5,
		MarkerHeight: 5,
		Orient:       MarkerOrientAuto,
		Content:      content,
	})
}

// SquareMarker creates a square marker
func SquareMarker(id string, color string) string {
	content := fmt.Sprintf(`<rect x="1" y="1" width="8" height="8" fill="%s"/>`, escapeAttr(color))
	return Marker(MarkerDef{
		ID:           id,
		ViewBox:      "0 0 10 10",
		RefX:         5,
		RefY:         5,
		MarkerWidth:  5,
		MarkerHeight: 5,
		Orient:       MarkerOrientAuto,
		Content:      content,
	})
}

// DiamondMarker creates a diamond marker
func DiamondMarker(id string, color string) string {
	content := fmt.Sprintf(`<path d="M 5 1 L 9 5 L 5 9 L 1 5 Z" fill="%s"/>`, escapeAttr(color))
	return Marker(MarkerDef{
		ID:           id,
		ViewBox:      "0 0 10 10",
		RefX:         5,
		RefY:         5,
		MarkerWidth:  5,
		MarkerHeight: 5,
		Orient:       MarkerOrientAuto,
		Content:      content,
	})
}

// TriangleMarker creates a triangle marker
func TriangleMarker(id string, color string) string {
	content := fmt.Sprintf(`<path d="M 5 1 L 9 9 L 1 9 Z" fill="%s"/>`, escapeAttr(color))
	return Marker(MarkerDef{
		ID:           id,
		ViewBox:      "0 0 10 10",
		RefX:         5,
		RefY:         9,
		MarkerWidth:  5,
		MarkerHeight: 5,
		Orient:       MarkerOrientAuto,
		Content:      content,
	})
}

// CrossMarker creates a cross/plus marker
func CrossMarker(id string, color string, strokeWidth float64) string {
	content := fmt.Sprintf(`<path d="M 5 1 L 5 9 M 1 5 L 9 5" stroke="%s" stroke-width="%.1f" stroke-linecap="round"/>`, escapeAttr(color), strokeWidth)
	return Marker(MarkerDef{
		ID:           id,
		ViewBox:      "0 0 10 10",
		RefX:         5,
		RefY:         5,
		MarkerWidth:  5,
		MarkerHeight: 5,
		Orient:       MarkerOrientAuto,
		Content:      content,
	})
}

// XMarker creates an X marker
func XMarker(id string, color string, strokeWidth float64) string {
	content := fmt.Sprintf(`<path d="M 2 2 L 8 8 M 8 2 L 2 8" stroke="%s" stroke-width="%.1f" stroke-linecap="round"/>`, escapeAttr(color), strokeWidth)
	return Marker(MarkerDef{
		ID:           id,
		ViewBox:      "0 0 10 10",
		RefX:         5,
		RefY:         5,
		MarkerWidth:  5,
		MarkerHeight: 5,
		Orient:       MarkerOrientAuto,
		Content:      content,
	})
}

// DotMarker creates a small dot marker (good for data points)
func DotMarker(id string, color string, radius float64) string {
	content := fmt.Sprintf(`<circle cx="5" cy="5" r="%.1f" fill="%s"/>`, radius, escapeAttr(color))
	return Marker(MarkerDef{
		ID:           id,
		ViewBox:      "0 0 10 10",
		RefX:         5,
		RefY:         5,
		MarkerWidth:  4,
		MarkerHeight: 4,
		Orient:       MarkerOrientAuto,
		Content:      content,
	})
}

// Helper to apply markers to path style attributes
func applyMarkers(style Style, markerStart, markerMid, markerEnd string) string {
	if markerStart != "" {
		style.MarkerStart = markerStart
	}
	if markerMid != "" {
		style.MarkerMid = markerMid
	}
	if markerEnd != "" {
		style.MarkerEnd = markerEnd
	}
	return formatStyle(style)
}

// PathWithMarkers renders a path with marker references
func PathWithMarkers(d string, style Style, markerStart, markerMid, markerEnd string) string {
	attrs := applyMarkers(style, markerStart, markerMid, markerEnd)
	return fmt.Sprintf(`<path d="%s"%s/>`, escapeAttr(d), attrs)
}

// LineWithMarkers renders a line with marker references
func LineWithMarkers(x1, y1, x2, y2 float64, style Style, markerStart, markerEnd string) string {
	attrs := applyMarkers(style, markerStart, "", markerEnd)
	return fmt.Sprintf(`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f"%s/>`,
		x1, y1, x2, y2, attrs)
}

// PolylineWithMarkers renders a polyline with marker references
func PolylineWithMarkers(points []Point, style Style, markerStart, markerMid, markerEnd string) string {
	if len(points) == 0 {
		return ""
	}

	var pointsStr strings.Builder
	for i, p := range points {
		if i > 0 {
			pointsStr.WriteString(" ")
		}
		fmt.Fprintf(&pointsStr, "%.2f,%.2f", p.X, p.Y)
	}

	attrs := applyMarkers(style, markerStart, markerMid, markerEnd)
	return fmt.Sprintf(`<polyline points="%s"%s/>`, pointsStr.String(), attrs)
}
