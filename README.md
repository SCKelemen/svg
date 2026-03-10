# svg

A clean, efficient SVG rendering library that integrates with the [layout](https://github.com/SCKelemen/layout) engine to produce beautiful SVG graphics from layout trees. Also provides standalone SVG path construction and marker support for data visualization.

## Features

- **Layout Integration**: Seamlessly renders `layout.Node` trees to SVG
- **PathBuilder**: Fluent API for constructing complex SVG paths with bezier curves
- **Markers**: Define custom markers for path endpoints and vertices (arrows, circles, diamonds, etc.)
- **Smooth Curves**: Bezier curve interpolation with tension control
- **Basic Shapes**: Rect, Circle, Ellipse, Polygon, Polyline, Line, Text
- **Transform Support**: Full 2D transform support (translate, rotate, scale, skew)
- **Styling System**: Colors, borders, backgrounds, shadows
- **Text Rendering**: SVG text elements with proper positioning
- **ClipPath Management**: Thread-safe unique ID generation for clipping
- **Gradient Support**: Linear and radial gradients with multiple color spaces (OKLCH, OKLAB, sRGB, Display P3)
- **Design Tokens**: Themeable styling system

## Installation

```bash
go get github.com/SCKelemen/svg
```

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "github.com/SCKelemen/layout"
    "github.com/SCKelemen/svg"
)

func main() {
    // Create a layout tree
    root := &layout.Node{
        Style: layout.Style{
            Display: layout.DisplayFlex,
            FlexDirection: layout.FlexDirectionRow,
            Width: layout.Px(400),
            Height: layout.Px(200),
        },
        Children: []*layout.Node{
            {Style: layout.Style{Width: layout.Px(100), Height: layout.Px(100)}},
            {Style: layout.Style{Width: layout.Px(100), Height: layout.Px(100)}},
        },
    }

    // Perform layout
    constraints := layout.Loose(800, 600)
    layout.Layout(root, constraints, nil)

    // Render to SVG
    output := svg.RenderToSVG(root, svg.Options{
        Width: 400,
        Height: 200,
    })

    fmt.Println(output)
}
```

### With Styling

```go
// Create a styled renderer
renderer := svg.NewRenderer(svg.Options{
    Width: 800,
    Height: 600,
    StyleSheet: svg.DefaultStyleSheet(),
})

// Render with custom styles
output := renderer.Render(root)
```

### Gradients

```go
// Create a simple two-stop linear gradient definition
gradientDef := svg.LinearGradient(svg.LinearGradientDef{
    ID: "myGradient",
    X1: "0%", Y1: "0%",
    X2: "100%", Y2: "0%",
    Stops: []svg.GradientStop{
        {Offset: "0%", Color: "#3B82F6"},
        {Offset: "100%", Color: "#8B5CF6"},
    },
})

// Apply gradient to an element
svgElement := fmt.Sprintf(`<rect fill="%s" x="0" y="0" width="100" height="50"/>`, svg.GradientURL("myGradient"))
_ = gradientDef
```

### PathBuilder - Fluent API for Paths

Create complex SVG paths using a chainable API:

```go
// Build a smooth curved path
points := []svg.Point{
    {X: 50, Y: 200},
    {X: 100, Y: 100},
    {X: 200, Y: 150},
    {X: 300, Y: 50},
}

// Simple polyline
path := svg.PolylinePath(points)

// Smooth bezier curve with tension control
smoothPath := svg.SmoothLinePath(points, 0.3)

// Area chart path (filled region)
areaPath := svg.AreaPath(points, 200) // baseline at y=200

// Smooth area chart
smoothArea := svg.SmoothAreaPath(points, 200, 0.3)

// Manual path construction
pb := svg.NewPathBuilder().
    MoveTo(10, 10).
    LineTo(90, 10).
    CurveTo(120, 10, 120, 40, 90, 40).
    Close()
```

### Markers - Path Decorations

Add markers (arrows, dots, shapes) to path endpoints:

```go
// Create predefined markers
arrow := svg.ArrowMarker("arrow-blue", "#3B82F6")
circle := svg.CircleMarker("dot-green", "#10B981")
diamond := svg.DiamondMarker("diamond-purple", "#8B5CF6")

// Apply markers to a path
pathData := svg.SmoothLinePath(points, 0.3)
decoratedPath := svg.PathWithMarkers(
    pathData,
    svg.Style{Stroke: "#3B82F6", StrokeWidth: 2, Fill: "none"},
    svg.MarkerURL("dot-green"),   // start marker
    "",                            // mid markers (optional)
    svg.MarkerURL("arrow-blue"),  // end marker
)

// Build an SVG manually with marker definitions in <defs>
output := fmt.Sprintf(`<svg width="400" height="300" viewBox="0 0 400 300" xmlns="http://www.w3.org/2000/svg">
<defs>
%s
%s
%s
</defs>
%s
</svg>`, arrow, circle, diamond, decoratedPath)
```

Available marker types:
- `ArrowMarker` - Directional arrow
- `CircleMarker` - Filled circle
- `SquareMarker` - Filled square
- `DiamondMarker` - Diamond shape
- `TriangleMarker` - Triangle shape
- `CrossMarker` - Plus/cross symbol
- `XMarker` - X symbol
- `DotMarker` - Small dot (customizable radius)

## Design Philosophy

This library focuses on:

1. **Simplicity**: Clean API with sensible defaults
2. **Integration**: First-class support for the layout engine
3. **Extensibility**: Easy to add custom rendering logic
4. **Performance**: Efficient string building and minimal allocations

## Testing

The library includes comprehensive unit tests:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. -benchmem
```

Coverage: ~60%+ of statements with tests for PathBuilder, markers, rendering, and export.

## Performance

Typical performance on modern hardware:

| Operation | Time/op | Allocations |
|-----------|---------|-------------|
| PolylinePath (100 points) | ~2-3 μs | Minimal |
| SmoothLinePath (100 points) | ~15-20 μs | Moderate |
| Marker generation | ~200-500 ns | Low |
| Full SVG render | ~10-50 μs | Context-dependent |

## Related Projects

- [layout](https://github.com/SCKelemen/layout) - CSS Grid/Flexbox layout engine
- [text](https://github.com/SCKelemen/text) - Unicode text handling
- [color](https://github.com/SCKelemen/color) - Color space handling
- [cli](https://github.com/SCKelemen/cli) - Terminal rendering
