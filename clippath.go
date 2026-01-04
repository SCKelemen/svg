package rendersvg

import (
	"fmt"
	"strings"
	"sync/atomic"
)

// Global counter for unique clipPath IDs across all renderers
// This ensures thread-safe unique ID generation
var clipPathCounter int64

// ClipPathManager manages SVG clipPath definitions and generates unique IDs
type ClipPathManager struct {
	paths []ClipPath
}

// ClipPath represents an SVG clipPath definition
type ClipPath struct {
	ID   string
	Path string // SVG path or shapes to use for clipping
}

// NewClipPathManager creates a new clipPath manager
func NewClipPathManager() *ClipPathManager {
	return &ClipPathManager{
		paths: make([]ClipPath, 0),
	}
}

// GenerateID generates a unique clipPath ID
func (m *ClipPathManager) GenerateID() string {
	id := atomic.AddInt64(&clipPathCounter, 1)
	return fmt.Sprintf("clip-%d", id)
}

// AddRoundedRect adds a rounded rectangle clipPath and returns its ID
func (m *ClipPathManager) AddRoundedRect(x, y, width, height, radius float64) string {
	id := m.GenerateID()
	path := fmt.Sprintf(`<rect x="%.2f" y="%.2f" width="%.2f" height="%.2f" rx="%.2f" ry="%.2f"/>`,
		x, y, width, height, radius, radius)

	m.paths = append(m.paths, ClipPath{
		ID:   id,
		Path: path,
	})

	return id
}

// AddRect adds a rectangle clipPath and returns its ID
func (m *ClipPathManager) AddRect(x, y, width, height float64) string {
	id := m.GenerateID()
	path := fmt.Sprintf(`<rect x="%.2f" y="%.2f" width="%.2f" height="%.2f"/>`,
		x, y, width, height)

	m.paths = append(m.paths, ClipPath{
		ID:   id,
		Path: path,
	})

	return id
}

// AddCircle adds a circle clipPath and returns its ID
func (m *ClipPathManager) AddCircle(cx, cy, r float64) string {
	id := m.GenerateID()
	path := fmt.Sprintf(`<circle cx="%.2f" cy="%.2f" r="%.2f"/>`,
		cx, cy, r)

	m.paths = append(m.paths, ClipPath{
		ID:   id,
		Path: path,
	})

	return id
}

// AddCustom adds a custom clipPath and returns its ID
func (m *ClipPathManager) AddCustom(pathContent string) string {
	id := m.GenerateID()
	m.paths = append(m.paths, ClipPath{
		ID:   id,
		Path: pathContent,
	})
	return id
}

// ToSVGDefs converts all clipPaths to SVG <defs> content
func (m *ClipPathManager) ToSVGDefs() string {
	if len(m.paths) == 0 {
		return ""
	}

	var b strings.Builder

	for _, cp := range m.paths {
		b.WriteString(fmt.Sprintf(`<clipPath id="%s">%s</clipPath>`, cp.ID, cp.Path))
		b.WriteString("\n    ")
	}

	return b.String()
}

// URL returns the CSS url() reference for a clipPath ID
func URL(id string) string {
	return fmt.Sprintf("url(#%s)", id)
}
