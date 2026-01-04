// Package rendersvg provides efficient SVG rendering for layout trees.
//
// This package integrates with github.com/SCKelemen/layout to render
// computed layout trees as SVG graphics. It supports transforms, styling,
// text rendering, and clipPath management.
//
// Basic usage:
//
//	root := &layout.Node{...}
//	layout.Layout(root, constraints, nil)
//	svg := rendersvg.RenderToSVG(root, rendersvg.Options{
//		Width: 800,
//		Height: 600,
//	})
//
package rendersvg
