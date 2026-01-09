package svg

import (
	"strings"
	"testing"
)

func TestPathBuilder_Basic(t *testing.T) {
	tests := []struct {
		name     string
		builder  func() *PathBuilder
		expected string
	}{
		{
			name: "MoveTo and LineTo",
			builder: func() *PathBuilder {
				return NewPathBuilder().
					MoveTo(10, 20).
					LineTo(30, 40)
			},
			expected: "M 10.00 20.00 L 30.00 40.00",
		},
		{
			name: "HorizontalLineTo and VerticalLineTo",
			builder: func() *PathBuilder {
				return NewPathBuilder().
					MoveTo(0, 0).
					HorizontalLineTo(50).
					VerticalLineTo(50)
			},
			expected: "M 0.00 0.00 H 50.00 V 50.00",
		},
		{
			name: "CurveTo",
			builder: func() *PathBuilder {
				return NewPathBuilder().
					MoveTo(10, 80).
					CurveTo(40, 10, 65, 10, 95, 80)
			},
			expected: "M 10.00 80.00 C 40.00 10.00, 65.00 10.00, 95.00 80.00",
		},
		{
			name: "SmoothCurveTo",
			builder: func() *PathBuilder {
				return NewPathBuilder().
					MoveTo(10, 80).
					CurveTo(40, 10, 65, 10, 95, 80).
					SmoothCurveTo(150, 150, 180, 80)
			},
			expected: "M 10.00 80.00 C 40.00 10.00, 65.00 10.00, 95.00 80.00 S 150.00 150.00, 180.00 80.00",
		},
		{
			name: "QuadraticCurveTo",
			builder: func() *PathBuilder {
				return NewPathBuilder().
					MoveTo(10, 80).
					QuadraticCurveTo(52.5, 10, 95, 80)
			},
			expected: "M 10.00 80.00 Q 52.50 10.00, 95.00 80.00",
		},
		{
			name: "ArcTo",
			builder: func() *PathBuilder {
				return NewPathBuilder().
					MoveTo(10, 50).
					ArcTo(40, 40, 0, 0, 1, 90, 50)
			},
			expected: "M 10.00 50.00 A 40.00 40.00 0.00 0 1 90.00 50.00",
		},
		{
			name: "Close",
			builder: func() *PathBuilder {
				return NewPathBuilder().
					MoveTo(10, 10).
					LineTo(90, 10).
					LineTo(90, 90).
					LineTo(10, 90).
					Close()
			},
			expected: "M 10.00 10.00 L 90.00 10.00 L 90.00 90.00 L 10.00 90.00 Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strings.TrimSpace(tt.builder().String())
			expected := strings.TrimSpace(tt.expected)
			if result != expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
			}
		})
	}
}

func TestPathBuilder_EmptyPath(t *testing.T) {
	pb := NewPathBuilder()
	result := pb.String()
	if result != "" {
		t.Errorf("Expected empty string, got: %s", result)
	}
}

func TestPolylinePath(t *testing.T) {
	points := []Point{
		{X: 0, Y: 0},
		{X: 10, Y: 20},
		{X: 30, Y: 15},
		{X: 40, Y: 40},
	}

	result := PolylinePath(points)
	expected := "M 0.00 0.00 L 10.00 20.00 L 30.00 15.00 L 40.00 40.00"

	if strings.TrimSpace(result) != strings.TrimSpace(expected) {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestPolylinePath_EmptyPoints(t *testing.T) {
	result := PolylinePath([]Point{})
	if result != "" {
		t.Errorf("Expected empty string for empty points, got: %s", result)
	}
}

func TestPolylinePath_SinglePoint(t *testing.T) {
	points := []Point{{X: 10, Y: 20}}
	result := PolylinePath(points)
	expected := "M 10.00 20.00"

	if strings.TrimSpace(result) != strings.TrimSpace(expected) {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestPolygonPath(t *testing.T) {
	points := []Point{
		{X: 50, Y: 0},
		{X: 100, Y: 50},
		{X: 50, Y: 100},
		{X: 0, Y: 50},
	}

	result := PolygonPath(points)

	// Should start with M, have L commands, and end with Z
	if !strings.HasPrefix(result, "M ") {
		t.Errorf("Expected path to start with 'M ', got: %s", result)
	}
	if !strings.HasSuffix(strings.TrimSpace(result), "Z") {
		t.Errorf("Expected path to end with 'Z', got: %s", result)
	}
}

func TestSmoothLinePath(t *testing.T) {
	points := []Point{
		{X: 0, Y: 100},
		{X: 50, Y: 50},
		{X: 100, Y: 80},
		{X: 150, Y: 20},
	}

	result := SmoothLinePath(points, 0.3)

	// Should start with M and contain C (cubic bezier) commands
	if !strings.HasPrefix(result, "M ") {
		t.Errorf("Expected path to start with 'M ', got: %s", result)
	}
	if !strings.Contains(result, "C ") {
		t.Errorf("Expected path to contain 'C ' (cubic bezier), got: %s", result)
	}
}

func TestSmoothLinePath_TwoPoints(t *testing.T) {
	points := []Point{
		{X: 0, Y: 0},
		{X: 100, Y: 100},
	}

	result := SmoothLinePath(points, 0.3)
	// With only 2 points, should fall back to simple line
	expected := "M 0.00 0.00 L 100.00 100.00"

	if strings.TrimSpace(result) != strings.TrimSpace(expected) {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestSmoothLinePath_TensionZero(t *testing.T) {
	points := []Point{
		{X: 0, Y: 100},
		{X: 50, Y: 50},
		{X: 100, Y: 80},
	}

	// Tension of 0 should still produce a valid path
	result := SmoothLinePath(points, 0)

	if !strings.HasPrefix(result, "M ") {
		t.Errorf("Expected path to start with 'M ', got: %s", result)
	}
}

func TestAreaPath(t *testing.T) {
	points := []Point{
		{X: 0, Y: 50},
		{X: 50, Y: 20},
		{X: 100, Y: 80},
	}
	baselineY := 100.0

	result := AreaPath(points, baselineY)

	// Should start with M to baseline
	if !strings.HasPrefix(result, "M ") {
		t.Errorf("Expected path to start with 'M ', got: %s", result)
	}
	// Should contain L commands for the points
	if !strings.Contains(result, "L ") {
		t.Errorf("Expected path to contain 'L ' commands, got: %s", result)
	}
	// Should end with Z to close the path
	if !strings.HasSuffix(strings.TrimSpace(result), "Z") {
		t.Errorf("Expected path to end with 'Z', got: %s", result)
	}
	// Should contain the baseline Y value
	if !strings.Contains(result, "100.00") {
		t.Errorf("Expected path to contain baseline Y value, got: %s", result)
	}
}

func TestSmoothAreaPath(t *testing.T) {
	points := []Point{
		{X: 0, Y: 50},
		{X: 50, Y: 20},
		{X: 100, Y: 80},
		{X: 150, Y: 40},
	}
	baselineY := 100.0
	tension := 0.3

	result := SmoothAreaPath(points, baselineY, tension)

	// Should start with M to baseline
	if !strings.HasPrefix(result, "M ") {
		t.Errorf("Expected path to start with 'M ', got: %s", result)
	}
	// Should contain C (cubic bezier) commands
	if !strings.Contains(result, "C ") {
		t.Errorf("Expected path to contain 'C ' commands, got: %s", result)
	}
	// Should end with Z to close the path
	if !strings.HasSuffix(strings.TrimSpace(result), "Z") {
		t.Errorf("Expected path to end with 'Z', got: %s", result)
	}
}

func BenchmarkPolylinePath(b *testing.B) {
	points := make([]Point, 100)
	for i := range points {
		points[i] = Point{X: float64(i), Y: float64(i * 2)}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = PolylinePath(points)
	}
}

func BenchmarkSmoothLinePath(b *testing.B) {
	points := make([]Point, 100)
	for i := range points {
		points[i] = Point{X: float64(i), Y: float64(i * 2)}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SmoothLinePath(points, 0.3)
	}
}
