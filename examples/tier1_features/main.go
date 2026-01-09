package main

import (
	"fmt"

	"github.com/SCKelemen/svg"
)

func main() {
	fmt.Println("=== SVG Tier 1 Features Demo ===")
	fmt.Println()

	// 1. PathBuilder
	fmt.Println("1. PathBuilder - Clean path construction")
	pb := svg.NewPathBuilder()
	pb.MoveTo(10, 10).LineTo(90, 10).LineTo(90, 90).LineTo(10, 90).Close()
	fmt.Printf("   Square path: %s\n", pb.String())

	// Smooth curve
	points := []svg.Point{
		{X: 10, Y: 50},
		{X: 30, Y: 20},
		{X: 50, Y: 80},
		{X: 70, Y: 30},
		{X: 90, Y: 60},
	}
	smoothPath := svg.SmoothLinePath(points, 0.3)
	fmt.Printf("   Smooth curve: %s\n", smoothPath[:50])
	fmt.Println()

	// 2. Polygon & Polyline
	fmt.Println("2. Polygon & Polyline - Multi-point shapes")
	trianglePoints := []svg.Point{
		{X: 50, Y: 10},
		{X: 90, Y: 90},
		{X: 10, Y: 90},
	}
	triangle := svg.Polygon(trianglePoints, svg.Style{
		Fill:   "#3B82F6",
		Stroke: "#1E40AF",
	})
	fmt.Printf("   Triangle: %s\n", triangle[:60])

	zigzag := []svg.Point{
		{X: 0, Y: 50},
		{X: 25, Y: 20},
		{X: 50, Y: 80},
		{X: 75, Y: 30},
		{X: 100, Y: 70},
	}
	polyline := svg.Polyline(zigzag, svg.Style{
		Fill:        "none",
		Stroke:      "#EF4444",
		StrokeWidth: 2,
	})
	fmt.Printf("   Zigzag: %s\n", polyline[:60])
	fmt.Println()

	// 3. Markers
	fmt.Println("3. Markers - Arrowheads and data point shapes")
	arrowMarker := svg.ArrowMarker("arrow", "#3B82F6")
	fmt.Printf("   Arrow marker definition:\n   %s\n", arrowMarker[:80])

	circleMarker := svg.CircleMarker("dot", "#EF4444")
	fmt.Printf("   Circle marker definition:\n   %s\n", circleMarker[:80])

	// Line with arrow at end
	lineWithArrow := svg.LineWithMarkers(10, 50, 90, 50, svg.Style{
		Stroke:      "#3B82F6",
		StrokeWidth: 2,
	}, "", svg.MarkerURL("arrow"))
	fmt.Printf("   Line with arrow: %s\n", lineWithArrow[:80])
	fmt.Println()

	// 4. Area Chart Path
	fmt.Println("4. Area Chart - Using AreaPath")
	areaPoints := []svg.Point{
		{X: 0, Y: 30},
		{X: 20, Y: 15},
		{X: 40, Y: 40},
		{X: 60, Y: 25},
		{X: 80, Y: 50},
		{X: 100, Y: 20},
	}
	areaPath := svg.AreaPath(areaPoints, 100) // baseline at y=100
	fmt.Printf("   Area path: %s\n", areaPath[:60])
	fmt.Println()

	// 5. Complete SVG Example
	fmt.Println("5. Complete SVG with all features:")
	fmt.Println(generateCompleteSVG())
}

func generateCompleteSVG() string {
	// Define markers in defs
	markers := svg.ArrowMarker("arrow", "#3B82F6") + "\n" +
		svg.CircleMarker("dot", "#EF4444") + "\n" +
		svg.DiamondMarker("diamond", "#10B981") + "\n" +
		svg.SquareMarker("square", "#F59E0B")

	// Create content with various shapes
	var content string

	// 1. Line graph using PathBuilder
	graphPoints := []svg.Point{
		{X: 50, Y: 200},
		{X: 100, Y: 150},
		{X: 150, Y: 180},
		{X: 200, Y: 120},
		{X: 250, Y: 160},
		{X: 300, Y: 100},
		{X: 350, Y: 140},
	}
	smoothGraph := svg.SmoothLinePath(graphPoints, 0.2)
	content += svg.PathWithMarkers(smoothGraph, svg.Style{
		Fill:           "none",
		Stroke:         "#3B82F6",
		StrokeWidth:    3,
		StrokeLinecap:  svg.StrokeLinecapRound,
		StrokeLinejoin: svg.StrokeLinejoinRound,
	}, "", svg.MarkerURL("dot"), svg.MarkerURL("arrow"))
	content += "\n"

	// 2. Area chart using AreaPath
	areaPoints := []svg.Point{
		{X: 400, Y: 250},
		{X: 450, Y: 220},
		{X: 500, Y: 280},
		{X: 550, Y: 240},
	}
	areaPath := svg.SmoothAreaPath(areaPoints, 350, 0.2)
	content += svg.Path(areaPath, svg.Style{
		Fill:        "#10B981",
		FillOpacity: 0.3,
		Stroke:      "#10B981",
		StrokeWidth: 2,
	})
	content += "\n"

	// 3. Polygon - triangle
	triangle := []svg.Point{
		{X: 100, Y: 300},
		{X: 150, Y: 350},
		{X: 50, Y: 350},
	}
	content += svg.Polygon(triangle, svg.Style{
		Fill:        "#F59E0B",
		Stroke:      "#D97706",
		StrokeWidth: 2,
	})
	content += "\n"

	// 4. Ellipse
	content += svg.Ellipse(500, 100, 40, 25, svg.Style{
		Fill:        "#8B5CF6",
		Stroke:      "#7C3AED",
		StrokeWidth: 2,
	})
	content += "\n"

	// Add title
	content += svg.Text("SVG Tier 1 Features", 300, 30, svg.Style{
		Fill:             "#E5E7EB",
		TextAnchor:       svg.TextAnchorMiddle,
		DominantBaseline: svg.DominantBaselineMiddle,
		Class:            "sans",
	})

	// Build complete SVG
	var fullSVG string
	fullSVG += fmt.Sprintf(`<svg width="600" height="400" viewBox="0 0 600 400" xmlns="http://www.w3.org/2000/svg">`)
	fullSVG += "\n<defs>\n"
	fullSVG += markers
	fullSVG += "\n</defs>\n"
	fullSVG += fmt.Sprintf(`<rect width="600" height="400" fill="#020617"/>`)
	fullSVG += "\n"
	fullSVG += content
	fullSVG += "\n</svg>"

	return fullSVG
}
