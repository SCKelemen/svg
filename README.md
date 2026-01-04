# render-svg

A clean, efficient SVG rendering library that integrates with the [layout](https://github.com/SCKelemen/layout) engine to produce beautiful SVG graphics from layout trees.

## Features

- **Layout Integration**: Seamlessly renders `layout.Node` trees to SVG
- **Transform Support**: Full 2D transform support (translate, rotate, scale, skew)
- **Styling System**: Colors, borders, backgrounds, shadows
- **Text Rendering**: SVG text elements with proper positioning
- **ClipPath Management**: Thread-safe unique ID generation for clipping
- **Design Tokens**: Themeable styling system

## Installation

```bash
go get github.com/SCKelemen/render-svg
```

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "github.com/SCKelemen/layout"
    "github.com/SCKelemen/render-svg"
)

func main() {
    // Create a layout tree
    root := &layout.Node{
        Style: layout.Style{
            Display: layout.DisplayFlex,
            FlexDirection: layout.FlexDirectionRow,
            Width: 400,
            Height: 200,
        },
        Children: []*layout.Node{
            {Style: layout.Style{Width: 100, Height: 100}},
            {Style: layout.Style{Width: 100, Height: 100}},
        },
    }

    // Perform layout
    constraints := layout.Loose(800, 600)
    layout.Layout(root, constraints, nil)

    // Render to SVG
    svg := rendersvg.RenderToSVG(root, rendersvg.Options{
        Width: 400,
        Height: 200,
    })

    fmt.Println(svg)
}
```

### With Styling

```go
// Create a styled renderer
renderer := rendersvg.NewRenderer(rendersvg.Options{
    Width: 800,
    Height: 600,
    StyleSheet: rendersvg.DefaultStyles(),
})

// Render with custom styles
svg := renderer.Render(root, func(node *layout.Node, depth int) rendersvg.Style {
    return rendersvg.Style{
        Fill: "#e0e0e0",
        Stroke: "#333",
        StrokeWidth: 1,
    }
})
```

## Design Philosophy

This library focuses on:

1. **Simplicity**: Clean API with sensible defaults
2. **Integration**: First-class support for the layout engine
3. **Extensibility**: Easy to add custom rendering logic
4. **Performance**: Efficient string building and minimal allocations

## Related Projects

- [layout](https://github.com/SCKelemen/layout) - CSS Grid/Flexbox layout engine
- [text](https://github.com/SCKelemen/text) - Unicode text handling
- [color](https://github.com/SCKelemen/color) - Color space handling
- [cli](https://github.com/SCKelemen/cli) - Terminal rendering
