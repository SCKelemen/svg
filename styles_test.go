package svg

import (
	"strings"
	"testing"
)

func TestStyleSheetToSVGDeterministicPropertyOrder(t *testing.T) {
	ss := &StyleSheet{
		Rules: []StyleRule{
			{
				Selector: ".x",
				Properties: map[string]string{
					"z-index": "9",
					"color":   "red",
					"border":  "1px solid black",
				},
			},
		},
	}

	first := ss.ToSVG()
	for i := 0; i < 10; i++ {
		if got := ss.ToSVG(); got != first {
			t.Fatalf("expected deterministic stylesheet output")
		}
	}

	colorIdx := strings.Index(first, "color:")
	borderIdx := strings.Index(first, "border:")
	zIdx := strings.Index(first, "z-index:")
	if !(borderIdx < colorIdx && colorIdx < zIdx) {
		t.Fatalf("expected sorted property order, got: %s", first)
	}
}
