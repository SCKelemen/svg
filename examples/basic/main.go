package main

import (
	"fmt"

	"github.com/SCKelemen/layout"
	rendersvg "github.com/SCKelemen/render-svg"
)

func main() {
	// Create a flexbox layout with two children
	root := &layout.Node{
		Style: layout.Style{
			Display:        layout.DisplayFlex,
			FlexDirection:  layout.FlexDirectionRow,
			JustifyContent: layout.JustifyContentSpaceBetween,
			AlignItems:     layout.AlignItemsCenter,
			Width:          layout.Px(400),
			Height:         layout.Px(200),
			Padding:        layout.Uniform(layout.Px(20)),
			FlexGap:        layout.Px(10),
		},
		Children: []*layout.Node{
			{
				Style: layout.Style{
					Width:  layout.Px(100),
					Height: layout.Px(100),
				},
			},
			{
				Style: layout.Style{
					Width:  layout.Px(100),
					Height: layout.Px(100),
				},
			},
		},
	}

	// Perform layout
	constraints := layout.Loose(800, 600)
	ctx := &layout.LayoutContext{
		ViewportWidth:  800,
		ViewportHeight: 600,
		RootFontSize:   16,
	}
	layout.Layout(root, constraints, ctx)

	// Render to SVG with custom styling
	opts := rendersvg.DefaultOptions()
	opts.Width = 400
	opts.Height = 200
	opts.BackgroundColor = "#f8f9fa"
	opts.StyleFunc = func(node interface{}, depth int) rendersvg.Style {
		if depth == 0 {
			// Root node - transparent
			return rendersvg.Style{
				Fill:   "none",
				Stroke: "#dee2e6",
			}
		}
		// Children - colored boxes
		return rendersvg.Style{
			Fill:        "#6366f1",
			Stroke:      "#4f46e5",
			StrokeWidth: 2,
		}
	}

	svg := rendersvg.RenderToSVG(root, opts)

	fmt.Println(svg)
}
