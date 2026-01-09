package svg

import (
	"strings"
	"testing"
)

func TestMarkerDef(t *testing.T) {
	tests := []struct {
		name     string
		marker   MarkerDef
		contains []string
	}{
		{
			name: "Basic marker definition",
			marker: MarkerDef{
				ID:           "test-marker",
				ViewBox:      "0 0 10 10",
				RefX:         5,
				RefY:         5,
				MarkerWidth:  6,
				MarkerHeight: 6,
				Orient:       MarkerOrientAuto,
				MarkerUnits:  MarkerUnitsStrokeWidth,
				Content:      `<circle cx="5" cy="5" r="2" fill="red"/>`,
			},
			contains: []string{
				`id="test-marker"`,
				`viewBox="0 0 10 10"`,
				`refX="5.00"`,
				`refY="5.00"`,
				`markerWidth="6.00"`,
				`markerHeight="6.00"`,
				`orient="auto"`,
				`markerUnits="strokeWidth"`,
				`<circle cx="5" cy="5" r="2" fill="red"/>`,
			},
		},
		{
			name: "Marker with custom orientation",
			marker: MarkerDef{
				ID:           "angle-marker",
				ViewBox:      "0 0 10 10",
				RefX:         0,
				RefY:         5,
				MarkerWidth:  10,
				MarkerHeight: 10,
				Orient:       "45",
				MarkerUnits:  MarkerUnitsUserSpaceOnUse,
				Content:      `<path d="M 0 0 L 10 5 L 0 10 z" fill="blue"/>`,
			},
			contains: []string{
				`id="angle-marker"`,
				`orient="45"`,
				`markerUnits="userSpaceOnUse"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Marker(tt.marker)
			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("Expected marker to contain %q, got:\n%s", substr, result)
				}
			}
		})
	}
}

func TestArrowMarker(t *testing.T) {
	result := ArrowMarker("arrow-red", "#FF0000")

	expectedContains := []string{
		`id="arrow-red"`,
		`#FF0000`,
		`<path`,
		`fill="#FF0000"`,
	}

	for _, substr := range expectedContains {
		if !strings.Contains(result, substr) {
			t.Errorf("ArrowMarker should contain %q, got:\n%s", substr, result)
		}
	}
}

func TestCircleMarker(t *testing.T) {
	result := CircleMarker("circle-blue", "#0000FF")

	expectedContains := []string{
		`id="circle-blue"`,
		`#0000FF`,
		`<circle`,
		`fill="#0000FF"`,
	}

	for _, substr := range expectedContains {
		if !strings.Contains(result, substr) {
			t.Errorf("CircleMarker should contain %q, got:\n%s", substr, result)
		}
	}
}

func TestSquareMarker(t *testing.T) {
	result := SquareMarker("square-green", "#00FF00")

	expectedContains := []string{
		`id="square-green"`,
		`#00FF00`,
		`<rect`,
		`fill="#00FF00"`,
	}

	for _, substr := range expectedContains {
		if !strings.Contains(result, substr) {
			t.Errorf("SquareMarker should contain %q, got:\n%s", substr, result)
		}
	}
}

func TestDiamondMarker(t *testing.T) {
	result := DiamondMarker("diamond-purple", "#800080")

	expectedContains := []string{
		`id="diamond-purple"`,
		`#800080`,
		`<path`,
		`fill="#800080"`,
	}

	for _, substr := range expectedContains {
		if !strings.Contains(result, substr) {
			t.Errorf("DiamondMarker should contain %q, got:\n%s", substr, result)
		}
	}
}

func TestTriangleMarker(t *testing.T) {
	result := TriangleMarker("triangle-orange", "#FFA500")

	expectedContains := []string{
		`id="triangle-orange"`,
		`#FFA500`,
		`<path`,
		`fill="#FFA500"`,
	}

	for _, substr := range expectedContains {
		if !strings.Contains(result, substr) {
			t.Errorf("TriangleMarker should contain %q, got:\n%s", substr, result)
		}
	}
}

func TestCrossMarker(t *testing.T) {
	result := CrossMarker("cross-black", "#000000", 1.5)

	expectedContains := []string{
		`id="cross-black"`,
		`#000000`,
		`stroke="#000000"`,
		`stroke-width="1.5"`,
		`<path`,
	}

	for _, substr := range expectedContains {
		if !strings.Contains(result, substr) {
			t.Errorf("CrossMarker should contain %q, got:\n%s", substr, result)
		}
	}
}

func TestXMarker(t *testing.T) {
	result := XMarker("x-red", "#FF0000", 2.0)

	expectedContains := []string{
		`id="x-red"`,
		`#FF0000`,
		`stroke="#FF0000"`,
		`stroke-width="2.0"`,
		`<path`,
	}

	for _, substr := range expectedContains {
		if !strings.Contains(result, substr) {
			t.Errorf("XMarker should contain %q, got:\n%s", substr, result)
		}
	}

	// X marker should have diagonal paths in the d attribute
	if !strings.Contains(result, "M 2 2 L 8 8 M 8 2 L 2 8") {
		t.Errorf("XMarker should contain diagonal path commands")
	}
}

func TestDotMarker(t *testing.T) {
	result := DotMarker("dot-cyan", "#00FFFF", 2.0)

	expectedContains := []string{
		`id="dot-cyan"`,
		`#00FFFF`,
		`<circle`,
		`fill="#00FFFF"`,
	}

	for _, substr := range expectedContains {
		if !strings.Contains(result, substr) {
			t.Errorf("DotMarker should contain %q, got:\n%s", substr, result)
		}
	}
}

func TestPathWithMarkers(t *testing.T) {
	pathData := "M 0 0 L 10 10 L 20 5"
	style := Style{
		Stroke:      "#000000",
		StrokeWidth: 2,
		Fill:        "none",
	}

	tests := []struct {
		name        string
		markerStart string
		markerMid   string
		markerEnd   string
		contains    []string
	}{
		{
			name:        "Only end marker",
			markerStart: "",
			markerMid:   "",
			markerEnd:   "url(#arrow-end)",
			contains: []string{
				`marker-end="url(#arrow-end)"`,
				pathData,
			},
		},
		{
			name:        "Start and end markers",
			markerStart: "url(#circle-start)",
			markerMid:   "",
			markerEnd:   "url(#arrow-end)",
			contains: []string{
				`marker-start="url(#circle-start)"`,
				`marker-end="url(#arrow-end)"`,
			},
		},
		{
			name:        "All markers",
			markerStart: "url(#circle-start)",
			markerMid:   "url(#dot-mid)",
			markerEnd:   "url(#arrow-end)",
			contains: []string{
				`marker-start="url(#circle-start)"`,
				`marker-mid="url(#dot-mid)"`,
				`marker-end="url(#arrow-end)"`,
			},
		},
		{
			name:        "No markers",
			markerStart: "",
			markerMid:   "",
			markerEnd:   "",
			contains: []string{
				pathData,
				`stroke="#000000"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PathWithMarkers(pathData, style, tt.markerStart, tt.markerMid, tt.markerEnd)
			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("Expected path to contain %q, got:\n%s", substr, result)
				}
			}

			// Should not contain empty marker attributes
			if tt.markerStart == "" && strings.Contains(result, `marker-start=""`) {
				t.Errorf("Should not contain empty marker-start attribute")
			}
			if tt.markerMid == "" && strings.Contains(result, `marker-mid=""`) {
				t.Errorf("Should not contain empty marker-mid attribute")
			}
			if tt.markerEnd == "" && strings.Contains(result, `marker-end=""`) {
				t.Errorf("Should not contain empty marker-end attribute")
			}
		})
	}
}

func TestMarkerOrient(t *testing.T) {
	tests := []struct {
		orient   MarkerOrient
		expected string
	}{
		{MarkerOrientAuto, "auto"},
		{MarkerOrientAutoStart, "auto-start-reverse"},
		{"0", "0"},
		{"45", "45"},
		{"90", "90"},
		{"180", "180"},
		{"-45", "-45"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			marker := MarkerDef{
				ID:      "test",
				ViewBox: "0 0 10 10",
				Orient:  tt.orient,
				Content: "<circle/>",
			}
			result := Marker(marker)
			expectedAttr := `orient="` + tt.expected + `"`
			if !strings.Contains(result, expectedAttr) {
				t.Errorf("Expected orient=%q, got:\n%s", tt.expected, result)
			}
		})
	}
}

func TestMarkerUnits(t *testing.T) {
	tests := []struct {
		units    MarkerUnits
		expected string
	}{
		{MarkerUnitsStrokeWidth, "strokeWidth"},
		{MarkerUnitsUserSpaceOnUse, "userSpaceOnUse"},
	}

	for _, tt := range tests {
		t.Run(string(tt.units), func(t *testing.T) {
			marker := MarkerDef{
				ID:          "test",
				ViewBox:     "0 0 10 10",
				MarkerUnits: tt.units,
				Content:     "<circle/>",
			}
			result := Marker(marker)
			expectedAttr := `markerUnits="` + tt.expected + `"`
			if !strings.Contains(result, expectedAttr) {
				t.Errorf("Expected markerUnits=%q, got:\n%s", tt.expected, result)
			}
		})
	}
}

func BenchmarkArrowMarker(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ArrowMarker("arrow", "#FF0000")
	}
}

func BenchmarkCircleMarker(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = CircleMarker("circle", "#0000FF")
	}
}

func BenchmarkPathWithMarkers(b *testing.B) {
	pathData := "M 0 0 L 10 10 L 20 5 L 30 15"
	style := Style{Stroke: "#000", StrokeWidth: 2}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = PathWithMarkers(pathData, style, "url(#start)", "url(#mid)", "url(#end)")
	}
}
