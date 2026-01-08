package rendersvg

// Options configures SVG rendering behavior
type Options struct {
	// Width of the output SVG
	Width float64

	// Height of the output SVG
	Height float64

	// ViewBox specifies a custom viewBox attribute (optional)
	// If empty, uses "0 0 {Width} {Height}"
	ViewBox string

	// StyleSheet to include in the SVG (optional)
	StyleSheet *StyleSheet

	// IncludeXMLDeclaration includes <?xml...?> declaration
	IncludeXMLDeclaration bool

	// Namespace adds xmlns attribute to <svg>
	Namespace bool

	// PreserveAspectRatio sets the preserveAspectRatio attribute
	PreserveAspectRatio string

	// BackgroundColor sets a background rectangle (optional)
	BackgroundColor string

	// StyleFunc allows custom styling per node
	// Called for each node with the node and its depth in the tree
	StyleFunc func(node interface{}, depth int) Style

	// RenderFunc allows custom rendering per node type
	// If it returns non-empty string, that's used instead of default rendering
	RenderFunc func(node interface{}, depth int) string
}

// DefaultOptions returns sensible default options
func DefaultOptions() Options {
	return Options{
		Width:                 800,
		Height:                600,
		StyleSheet:            DefaultStyleSheet(),
		IncludeXMLDeclaration: false,
		Namespace:             true,
		PreserveAspectRatio:   "xMidYMid meet",
	}
}

// WithSize creates options with specified dimensions
func WithSize(width, height float64) Options {
	opts := DefaultOptions()
	opts.Width = width
	opts.Height = height
	return opts
}

// WithStyleSheet creates options with a custom stylesheet
func WithStyleSheet(stylesheet *StyleSheet) Options {
	opts := DefaultOptions()
	opts.StyleSheet = stylesheet
	return opts
}
