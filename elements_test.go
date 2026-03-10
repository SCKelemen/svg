package svg

import (
	"strings"
	"testing"
)

func TestStyleOpacitySetAllowsZero(t *testing.T) {
	style := Style{
		Fill:       "#000000",
		Opacity:    0,
		OpacitySet: true,
	}

	rect := Rect(0, 0, 10, 10, style)
	if !strings.Contains(rect, `opacity="0.00"`) {
		t.Fatalf("expected explicit zero opacity, got: %s", rect)
	}
}

func TestGradientStopOpacitySetAllowsZero(t *testing.T) {
	gradient := LinearGradient(LinearGradientDef{
		ID: "g",
		Stops: []GradientStop{
			{
				Offset:     "0%",
				Color:      "#000000",
				Opacity:    0,
				OpacitySet: true,
			},
			{
				Offset: "100%",
				Color:  "#ffffff",
			},
		},
	})

	if !strings.Contains(gradient, `stop-opacity="0.00"`) {
		t.Fatalf("expected explicit zero stop opacity, got: %s", gradient)
	}
}

func TestPathEscapesAttributeValues(t *testing.T) {
	path := Path(`M 0 0 L 1 1 "`, Style{
		Stroke: `rgb(1,2,3)"onload="alert(1)`,
	})

	if strings.Contains(path, `"onload="`) {
		t.Fatalf("expected path attributes to be escaped, got: %s", path)
	}
	if !strings.Contains(path, "&quot;") {
		t.Fatalf("expected escaped quotes in output, got: %s", path)
	}
}
