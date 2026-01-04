package rendersvg

import (
	"strings"
	"testing"

	"github.com/SCKelemen/layout"
)

func TestRenderToSVG_Basic(t *testing.T) {
	root := &layout.Node{
		Style: layout.Style{
			Width:  layout.Px(100),
			Height: layout.Px(100),
		},
		Rect: layout.Rect{
			X:      0,
			Y:      0,
			Width:  100,
			Height: 100,
		},
	}

	opts := DefaultOptions()
	opts.Width = 200
	opts.Height = 200

	svg := RenderToSVG(root, opts)

	// Check SVG structure
	if !strings.Contains(svg, "<svg") {
		t.Error("Expected SVG tag")
	}
	if !strings.Contains(svg, `width="200"`) {
		t.Error("Expected width attribute")
	}
	if !strings.Contains(svg, `height="200"`) {
		t.Error("Expected height attribute")
	}
	if !strings.Contains(svg, "</svg>") {
		t.Error("Expected closing SVG tag")
	}
}

func TestRenderToSVG_WithChildren(t *testing.T) {
	root := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionRow,
			Width:         layout.Px(200),
			Height:        layout.Px(100),
		},
		Children: []*layout.Node{
			{
				Style: layout.Style{Width: layout.Px(100), Height: layout.Px(100)},
				Rect:  layout.Rect{X: 0, Y: 0, Width: 100, Height: 100},
			},
			{
				Style: layout.Style{Width: layout.Px(100), Height: layout.Px(100)},
				Rect:  layout.Rect{X: 100, Y: 0, Width: 100, Height: 100},
			},
		},
		Rect: layout.Rect{X: 0, Y: 0, Width: 200, Height: 100},
	}

	// Perform layout
	ctx := &layout.LayoutContext{
		ViewportWidth:  800,
		ViewportHeight: 600,
		RootFontSize:   16,
	}
	layout.Layout(root, layout.Loose(800, 600), ctx)

	opts := DefaultOptions()
	opts.Width = 200
	opts.Height = 100

	svg := RenderToSVG(root, opts)

	// Should contain multiple rectangles
	rectCount := strings.Count(svg, "<rect")
	if rectCount < 2 {
		t.Errorf("Expected at least 2 rectangles, got %d", rectCount)
	}
}

func TestRenderToSVG_WithStyleSheet(t *testing.T) {
	root := &layout.Node{
		Style: layout.Style{
			Width:  layout.Px(100),
			Height: layout.Px(100),
		},
		Rect: layout.Rect{
			X:      0,
			Y:      0,
			Width:  100,
			Height: 100,
		},
	}

	opts := DefaultOptions()
	opts.StyleSheet = DefaultStyleSheet()

	svg := RenderToSVG(root, opts)

	// Should contain style tag
	if !strings.Contains(svg, "<style>") {
		t.Error("Expected style tag")
	}
	if !strings.Contains(svg, ".sans") {
		t.Error("Expected .sans class in stylesheet")
	}
}

func TestRenderToSVG_WithBackground(t *testing.T) {
	root := &layout.Node{
		Style: layout.Style{
			Width:  layout.Px(100),
			Height: layout.Px(100),
		},
		Rect: layout.Rect{
			X:      0,
			Y:      0,
			Width:  100,
			Height: 100,
		},
	}

	opts := DefaultOptions()
	opts.BackgroundColor = "#ffffff"

	svg := RenderToSVG(root, opts)

	// Should contain background rectangle
	if !strings.Contains(svg, `fill="#ffffff"`) {
		t.Error("Expected background color")
	}
}

func TestClipPathManager(t *testing.T) {
	manager := NewClipPathManager()

	// Test rounded rect
	id1 := manager.AddRoundedRect(0, 0, 100, 100, 10)
	if id1 == "" {
		t.Error("Expected non-empty clipPath ID")
	}

	// Test rect
	id2 := manager.AddRect(0, 0, 50, 50)
	if id2 == "" {
		t.Error("Expected non-empty clipPath ID")
	}

	// IDs should be unique
	if id1 == id2 {
		t.Error("Expected unique clipPath IDs")
	}

	// Should generate SVG defs
	defs := manager.ToSVGDefs()
	if !strings.Contains(defs, id1) {
		t.Error("Expected first clipPath ID in defs")
	}
	if !strings.Contains(defs, id2) {
		t.Error("Expected second clipPath ID in defs")
	}
}

func TestElements(t *testing.T) {
	style := Style{
		Fill:   "#ff0000",
		Stroke: "#000000",
	}

	// Test Rect
	rect := Rect(10, 20, 100, 50, style)
	if !strings.Contains(rect, `x="10.00"`) {
		t.Error("Expected x coordinate")
	}
	if !strings.Contains(rect, `fill="#ff0000"`) {
		t.Error("Expected fill color")
	}

	// Test Circle
	circle := Circle(50, 50, 25, style)
	if !strings.Contains(circle, `cx="50.00"`) {
		t.Error("Expected cx coordinate")
	}

	// Test Text
	text := Text("Hello", 10, 20, style)
	if !strings.Contains(text, ">Hello</text>") {
		t.Error("Expected text content")
	}
}

func TestXMLEscape(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"<test>", "&lt;test&gt;"},
		{"&amp;", "&amp;amp;"},
		{`"quoted"`, "&quot;quoted&quot;"},
		{"normal text", "normal text"},
	}

	for _, tt := range tests {
		result := escapeXML(tt.input)
		if result != tt.expected {
			t.Errorf("escapeXML(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}
