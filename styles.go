package rendersvg

import "strings"

// StyleSheet represents a collection of CSS styles for SVG rendering
type StyleSheet struct {
	Rules []StyleRule
}

// StyleRule represents a single CSS rule
type StyleRule struct {
	Selector   string
	Properties map[string]string
}

// DefaultStyleSheet returns a sensible default stylesheet for SVG rendering
func DefaultStyleSheet() *StyleSheet {
	return &StyleSheet{
		Rules: []StyleRule{
			{
				Selector: ".sans",
				Properties: map[string]string{
					"font-family": `-apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji"`,
				},
			},
			{
				Selector: ".mono",
				Properties: map[string]string{
					"font-family":    `ui-monospace, SFMono-Regular, "SF Mono", Menlo, Consolas, "Liberation Mono", monospace`,
					"font-size":      "12px",
					"letter-spacing": "-0.5px",
				},
			},
			{
				Selector: ".bold",
				Properties: map[string]string{
					"font-weight": "500",
				},
			},
			{
				Selector: ".medium",
				Properties: map[string]string{
					"font-size": "16px",
				},
			},
			{
				Selector: ".small",
				Properties: map[string]string{
					"font-size": "14px",
				},
			},
			{
				Selector: ".smaller",
				Properties: map[string]string{
					"font-size": "12px",
				},
			},
			{
				Selector: ".pre",
				Properties: map[string]string{
					"white-space": "pre",
				},
			},
			{
				Selector: ".glow",
				Properties: map[string]string{
					"paint-order": "stroke",
				},
			},
		},
	}
}

// ToSVG converts the stylesheet to SVG <style> element
func (ss *StyleSheet) ToSVG() string {
	var b strings.Builder
	b.WriteString("<style>")

	for _, rule := range ss.Rules {
		b.WriteString("\n    ")
		b.WriteString(rule.Selector)
		b.WriteString(" {")

		for prop, value := range rule.Properties {
			b.WriteString("\n        ")
			b.WriteString(prop)
			b.WriteString(": ")
			b.WriteString(value)
			b.WriteString(";")
		}

		b.WriteString("\n    }")
	}

	b.WriteString("\n</style>")
	return b.String()
}

// AddRule adds a custom CSS rule to the stylesheet
func (ss *StyleSheet) AddRule(selector string, properties map[string]string) {
	ss.Rules = append(ss.Rules, StyleRule{
		Selector:   selector,
		Properties: properties,
	})
}
