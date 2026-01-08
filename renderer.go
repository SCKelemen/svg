package rendersvg

import (
	"fmt"
	"strings"

	"github.com/SCKelemen/layout"
)

// Renderer renders layout trees to SVG
type Renderer struct {
	options      Options
	clipPath     *ClipPathManager
	builder      strings.Builder
	defaultStyle Style
}

// NewRenderer creates a new SVG renderer with the given options
func NewRenderer(opts Options) *Renderer {
	return &Renderer{
		options:  opts,
		clipPath: NewClipPathManager(),
		defaultStyle: Style{
			Fill:   "#e0e0e0",
			Stroke: "#333",
		},
	}
}

// RenderToSVG renders a layout tree to an SVG string
func RenderToSVG(root *layout.Node, opts Options) string {
	renderer := NewRenderer(opts)
	return renderer.Render(root)
}

// Render renders the layout tree to SVG
func (r *Renderer) Render(root *layout.Node) string {
	r.builder.Reset()

	// XML declaration
	if r.options.IncludeXMLDeclaration {
		r.builder.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
		r.builder.WriteString("\n")
	}

	// Start SVG tag
	r.builder.WriteString("<svg")

	// Width and height
	r.builder.WriteString(fmt.Sprintf(` width="%.0f" height="%.0f"`, r.options.Width, r.options.Height))

	// ViewBox
	viewBox := r.options.ViewBox
	if viewBox == "" {
		viewBox = fmt.Sprintf("0 0 %.0f %.0f", r.options.Width, r.options.Height)
	}
	r.builder.WriteString(fmt.Sprintf(` viewBox="%s"`, viewBox))

	// Namespace
	if r.options.Namespace {
		r.builder.WriteString(` xmlns="http://www.w3.org/2000/svg"`)
	}

	// PreserveAspectRatio
	if r.options.PreserveAspectRatio != "" {
		r.builder.WriteString(fmt.Sprintf(` preserveAspectRatio="%s"`, r.options.PreserveAspectRatio))
	}

	r.builder.WriteString(">")
	r.builder.WriteString("\n")

	// Defs section
	r.builder.WriteString("<defs>")
	r.builder.WriteString("\n")

	// Stylesheet
	if r.options.StyleSheet != nil {
		r.builder.WriteString(r.options.StyleSheet.ToSVG())
		r.builder.WriteString("\n")
	}

	// Render nodes (this may add clipPaths)
	content := r.renderNode(root, 0)

	// ClipPaths (added during rendering)
	if clipDefs := r.clipPath.ToSVGDefs(); clipDefs != "" {
		r.builder.WriteString("    ")
		r.builder.WriteString(clipDefs)
	}

	r.builder.WriteString("</defs>")
	r.builder.WriteString("\n")

	// Background
	if r.options.BackgroundColor != "" {
		r.builder.WriteString(fmt.Sprintf(`<rect width="%.0f" height="%.0f" fill="%s"/>`,
			r.options.Width, r.options.Height, r.options.BackgroundColor))
		r.builder.WriteString("\n")
	}

	// Content
	r.builder.WriteString(content)

	// End SVG tag
	r.builder.WriteString("</svg>")

	return r.builder.String()
}

// renderNode recursively renders a layout node and its children
func (r *Renderer) renderNode(node *layout.Node, depth int) string {
	if node == nil {
		return ""
	}

	// Allow custom rendering
	if r.options.RenderFunc != nil {
		if custom := r.options.RenderFunc(node, depth); custom != "" {
			return custom
		}
	}

	var b strings.Builder
	rect := node.Rect

	// Get style for this node
	style := r.defaultStyle
	if r.options.StyleFunc != nil {
		style = r.options.StyleFunc(node, depth)
	}

	// Get transform
	transform := GetTransformFromNode(node)

	// Start group if there's a transform or children
	hasTransform := transform != ""
	hasChildren := len(node.Children) > 0

	if hasTransform || hasChildren {
		if hasTransform {
			b.WriteString(fmt.Sprintf(`<g transform="%s">`, transform))
		} else {
			b.WriteString("<g>")
		}
		b.WriteString("\n")
	}

	// Render the node itself as a rectangle
	// Only render if it has non-zero dimensions
	if rect.Width > 0 && rect.Height > 0 {
		indent := strings.Repeat("  ", depth+1)
		b.WriteString(indent)
		b.WriteString(Rect(rect.X, rect.Y, rect.Width, rect.Height, style))
		b.WriteString("\n")
	}

	// Render children
	for _, child := range node.Children {
		childContent := r.renderNode(child, depth+1)
		if childContent != "" {
			indent := strings.Repeat("  ", depth+1)
			b.WriteString(indent)
			b.WriteString(childContent)
		}
	}

	// End group
	if hasTransform || hasChildren {
		indent := strings.Repeat("  ", depth)
		b.WriteString(indent)
		b.WriteString("</g>")
		b.WriteString("\n")
	}

	return b.String()
}

// GetClipPathManager returns the clipPath manager for custom clipPath creation
func (r *Renderer) GetClipPathManager() *ClipPathManager {
	return r.clipPath
}

// SetDefaultStyle sets the default style for rendered nodes
func (r *Renderer) SetDefaultStyle(style Style) {
	r.defaultStyle = style
}

// RenderNodes renders multiple layout nodes at their computed positions
// This is useful when you have a collection of already-positioned nodes
func RenderNodes(nodes []*layout.Node, opts Options) string {
	renderer := NewRenderer(opts)

	var b strings.Builder

	// Start SVG
	b.WriteString(fmt.Sprintf(`<svg width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f"`,
		opts.Width, opts.Height, opts.Width, opts.Height))

	if opts.Namespace {
		b.WriteString(` xmlns="http://www.w3.org/2000/svg"`)
	}

	b.WriteString(">")
	b.WriteString("\n")

	// Defs
	b.WriteString("<defs>")
	b.WriteString("\n")

	if opts.StyleSheet != nil {
		b.WriteString(opts.StyleSheet.ToSVG())
		b.WriteString("\n")
	}

	b.WriteString("</defs>")
	b.WriteString("\n")

	// Background
	if opts.BackgroundColor != "" {
		b.WriteString(fmt.Sprintf(`<rect width="%.0f" height="%.0f" fill="%s"/>`,
			opts.Width, opts.Height, opts.BackgroundColor))
		b.WriteString("\n")
	}

	// Render each node
	for _, node := range nodes {
		content := renderer.renderNode(node, 0)
		b.WriteString(content)
	}

	b.WriteString("</svg>")

	return b.String()
}
