package rendersvg

import (
	"fmt"
	"strings"

	"github.com/SCKelemen/color"
)

// GradientStop represents a color stop in a gradient
type GradientStop struct {
	Offset  string // percentage or decimal (e.g., "0%", "50%", "1.0")
	Color   string // color value (hex, rgb, etc.)
	Opacity float64
}

// GradientUnits defines the coordinate system for gradients
type GradientUnits string

const (
	GradientUnitsUserSpaceOnUse GradientUnits = "userSpaceOnUse"
	GradientUnitsObjectBoundingBox GradientUnits = "objectBoundingBox"
)

// GradientSpreadMethod defines how gradient fills outside its bounds
type GradientSpreadMethod string

const (
	GradientSpreadPad     GradientSpreadMethod = "pad"
	GradientSpreadReflect GradientSpreadMethod = "reflect"
	GradientSpreadRepeat  GradientSpreadMethod = "repeat"
)

// LinearGradientDef represents a linear gradient definition
type LinearGradientDef struct {
	ID           string
	X1, Y1       string // Start point (can be percentage or absolute)
	X2, Y2       string // End point (can be percentage or absolute)
	Stops        []GradientStop
	Units        GradientUnits
	SpreadMethod GradientSpreadMethod
}

// RadialGradientDef represents a radial gradient definition
type RadialGradientDef struct {
	ID           string
	CX, CY       string // Center point (can be percentage or absolute)
	R            string // Radius (can be percentage or absolute)
	FX, FY       string // Focal point (optional, defaults to center)
	FR           string // Focal radius (optional)
	Stops        []GradientStop
	Units        GradientUnits
	SpreadMethod GradientSpreadMethod
}

// LinearGradient creates a linear gradient definition (for use in <defs>)
func LinearGradient(def LinearGradientDef) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf(`<linearGradient id="%s"`, def.ID))

	if def.X1 != "" {
		b.WriteString(fmt.Sprintf(` x1="%s"`, def.X1))
	}
	if def.Y1 != "" {
		b.WriteString(fmt.Sprintf(` y1="%s"`, def.Y1))
	}
	if def.X2 != "" {
		b.WriteString(fmt.Sprintf(` x2="%s"`, def.X2))
	}
	if def.Y2 != "" {
		b.WriteString(fmt.Sprintf(` y2="%s"`, def.Y2))
	}
	if def.Units != "" {
		b.WriteString(fmt.Sprintf(` gradientUnits="%s"`, string(def.Units)))
	}
	if def.SpreadMethod != "" {
		b.WriteString(fmt.Sprintf(` spreadMethod="%s"`, string(def.SpreadMethod)))
	}

	b.WriteString(">")
	b.WriteString("\n")

	// Add gradient stops
	for _, stop := range def.Stops {
		b.WriteString(fmt.Sprintf(`  <stop offset="%s" stop-color="%s"`, stop.Offset, stop.Color))
		if stop.Opacity > 0 && stop.Opacity < 1 {
			b.WriteString(fmt.Sprintf(` stop-opacity="%.2f"`, stop.Opacity))
		}
		b.WriteString(`/>`)
		b.WriteString("\n")
	}

	b.WriteString(`</linearGradient>`)
	return b.String()
}

// RadialGradient creates a radial gradient definition (for use in <defs>)
func RadialGradient(def RadialGradientDef) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf(`<radialGradient id="%s"`, def.ID))

	if def.CX != "" {
		b.WriteString(fmt.Sprintf(` cx="%s"`, def.CX))
	}
	if def.CY != "" {
		b.WriteString(fmt.Sprintf(` cy="%s"`, def.CY))
	}
	if def.R != "" {
		b.WriteString(fmt.Sprintf(` r="%s"`, def.R))
	}
	if def.FX != "" {
		b.WriteString(fmt.Sprintf(` fx="%s"`, def.FX))
	}
	if def.FY != "" {
		b.WriteString(fmt.Sprintf(` fy="%s"`, def.FY))
	}
	if def.FR != "" {
		b.WriteString(fmt.Sprintf(` fr="%s"`, def.FR))
	}
	if def.Units != "" {
		b.WriteString(fmt.Sprintf(` gradientUnits="%s"`, string(def.Units)))
	}
	if def.SpreadMethod != "" {
		b.WriteString(fmt.Sprintf(` spreadMethod="%s"`, string(def.SpreadMethod)))
	}

	b.WriteString(">")
	b.WriteString("\n")

	// Add gradient stops
	for _, stop := range def.Stops {
		b.WriteString(fmt.Sprintf(`  <stop offset="%s" stop-color="%s"`, stop.Offset, stop.Color))
		if stop.Opacity > 0 && stop.Opacity < 1 {
			b.WriteString(fmt.Sprintf(` stop-opacity="%.2f"`, stop.Opacity))
		}
		b.WriteString(`/>`)
		b.WriteString("\n")
	}

	b.WriteString(`</radialGradient>`)
	return b.String()
}

// GradientURL creates a url() reference to a gradient for use in fill or stroke
func GradientURL(id string) string {
	return fmt.Sprintf("url(#%s)", id)
}

// SimpleLinearGradient creates a simple two-color linear gradient
// angle is in degrees (0 = left to right, 90 = bottom to top)
func SimpleLinearGradient(id string, startColor, endColor string, angle float64) string {
	// Convert angle to x1,y1,x2,y2 coordinates
	// 0° = left to right (0,0 to 1,0)
	// 90° = bottom to top (0,1 to 0,0)
	// 180° = right to left (1,0 to 0,0)
	// 270° = top to bottom (0,0 to 0,1)

	x1, y1, x2, y2 := "0%", "0%", "100%", "0%"

	switch angle {
	case 0:
		x1, y1, x2, y2 = "0%", "0%", "100%", "0%"
	case 90:
		x1, y1, x2, y2 = "0%", "100%", "0%", "0%"
	case 180:
		x1, y1, x2, y2 = "100%", "0%", "0%", "0%"
	case 270:
		x1, y1, x2, y2 = "0%", "0%", "0%", "100%"
	case 45:
		x1, y1, x2, y2 = "0%", "100%", "100%", "0%"
	case 135:
		x1, y1, x2, y2 = "100%", "100%", "0%", "0%"
	case 225:
		x1, y1, x2, y2 = "100%", "0%", "0%", "100%"
	case 315:
		x1, y1, x2, y2 = "0%", "0%", "100%", "100%"
	}

	return LinearGradient(LinearGradientDef{
		ID: id,
		X1: x1,
		Y1: y1,
		X2: x2,
		Y2: y2,
		Stops: []GradientStop{
			{Offset: "0%", Color: startColor, Opacity: 1.0},
			{Offset: "100%", Color: endColor, Opacity: 1.0},
		},
	})
}

// SimpleRadialGradient creates a simple two-color radial gradient
func SimpleRadialGradient(id string, centerColor, edgeColor string) string {
	return RadialGradient(RadialGradientDef{
		ID: id,
		CX: "50%",
		CY: "50%",
		R:  "50%",
		Stops: []GradientStop{
			{Offset: "0%", Color: centerColor, Opacity: 1.0},
			{Offset: "100%", Color: edgeColor, Opacity: 1.0},
		},
	})
}

// InterpolatedLinearGradient creates a linear gradient with color interpolation in the specified color space
// This creates smooth, natural-looking gradients by interpolating in the chosen color space
// and converting to sRGB for SVG compatibility. Defaults to OKLCH for perceptually uniform gradients.
//
// colorSpace options:
//   - color.GradientRGB: Interpolates in RGB space (fast but not perceptually uniform)
//   - color.GradientHSL: Interpolates in HSL space
//   - color.GradientLAB: Interpolates in CIE LAB space
//   - color.GradientOKLAB: Interpolates in OKLAB space (perceptually uniform)
//   - color.GradientLCH: Interpolates in CIE LCH space
//   - color.GradientOKLCH: Interpolates in OKLCH space (perceptually uniform, recommended)
func InterpolatedLinearGradient(id string, startColor, endColor string, angle float64, steps int, colorSpace color.GradientSpace) (string, error) {
	if steps < 2 {
		steps = 2
	}

	// Parse start and end colors
	start, err := color.ParseColor(startColor)
	if err != nil {
		return "", fmt.Errorf("invalid start color: %w", err)
	}
	end, err := color.ParseColor(endColor)
	if err != nil {
		return "", fmt.Errorf("invalid end color: %w", err)
	}

	// Generate interpolated stops in the specified color space
	colors := color.GradientInSpace(start, end, steps, colorSpace)
	stops := make([]GradientStop, steps)
	for i := 0; i < steps; i++ {
		t := float64(i) / float64(steps-1)

		// Convert to hex for SVG
		hexColor := color.RGBToHex(colors[i])

		stops[i] = GradientStop{
			Offset:  fmt.Sprintf("%.1f%%", t*100),
			Color:   hexColor,
			Opacity: 1.0,
		}
	}

	// Convert angle to x1,y1,x2,y2 coordinates
	x1, y1, x2, y2 := angleToCoordinates(angle)

	return LinearGradient(LinearGradientDef{
		ID:    id,
		X1:    x1,
		Y1:    y1,
		X2:    x2,
		Y2:    y2,
		Stops: stops,
	}), nil
}

// OKLCHLinearGradient creates a perceptually uniform linear gradient using OKLCH interpolation
// This is a convenience wrapper around InterpolatedLinearGradient with color.GradientOKLCH
func OKLCHLinearGradient(id string, startColor, endColor string, angle float64, steps int) (string, error) {
	return InterpolatedLinearGradient(id, startColor, endColor, angle, steps, color.GradientOKLCH)
}

// InterpolatedRadialGradient creates a radial gradient with color interpolation in the specified color space
// This creates smooth, natural-looking gradients by interpolating in the chosen color space
// and converting to sRGB for SVG compatibility.
//
// colorSpace options: same as InterpolatedLinearGradient
func InterpolatedRadialGradient(id string, centerColor, edgeColor string, steps int, colorSpace color.GradientSpace) (string, error) {
	if steps < 2 {
		steps = 2
	}

	// Parse start and end colors
	center, err := color.ParseColor(centerColor)
	if err != nil {
		return "", fmt.Errorf("invalid center color: %w", err)
	}
	edge, err := color.ParseColor(edgeColor)
	if err != nil {
		return "", fmt.Errorf("invalid edge color: %w", err)
	}

	// Generate interpolated stops in the specified color space
	colors := color.GradientInSpace(center, edge, steps, colorSpace)
	stops := make([]GradientStop, steps)
	for i := 0; i < steps; i++ {
		t := float64(i) / float64(steps-1)

		// Convert to hex for SVG
		hexColor := color.RGBToHex(colors[i])

		stops[i] = GradientStop{
			Offset:  fmt.Sprintf("%.1f%%", t*100),
			Color:   hexColor,
			Opacity: 1.0,
		}
	}

	return RadialGradient(RadialGradientDef{
		ID:    id,
		CX:    "50%",
		CY:    "50%",
		R:     "50%",
		Stops: stops,
	}), nil
}

// OKLCHRadialGradient creates a perceptually uniform radial gradient using OKLCH interpolation
// This is a convenience wrapper around InterpolatedRadialGradient with color.GradientOKLCH
func OKLCHRadialGradient(id string, centerColor, edgeColor string, steps int) (string, error) {
	return InterpolatedRadialGradient(id, centerColor, edgeColor, steps, color.GradientOKLCH)
}

// angleToCoordinates converts an angle in degrees to SVG gradient coordinates
func angleToCoordinates(angle float64) (x1, y1, x2, y2 string) {
	x1, y1, x2, y2 = "0%", "0%", "100%", "0%"
	switch angle {
	case 0:
		x1, y1, x2, y2 = "0%", "0%", "100%", "0%"
	case 90:
		x1, y1, x2, y2 = "0%", "100%", "0%", "0%"
	case 180:
		x1, y1, x2, y2 = "100%", "0%", "0%", "0%"
	case 270:
		x1, y1, x2, y2 = "0%", "0%", "0%", "100%"
	case 45:
		x1, y1, x2, y2 = "0%", "100%", "100%", "0%"
	case 135:
		x1, y1, x2, y2 = "100%", "100%", "0%", "0%"
	case 225:
		x1, y1, x2, y2 = "100%", "0%", "0%", "100%"
	case 315:
		x1, y1, x2, y2 = "0%", "0%", "100%", "100%"
	}
	return
}
