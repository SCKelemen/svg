module github.com/SCKelemen/svg

go 1.25.4

require (
	github.com/SCKelemen/color v1.0.5
	github.com/SCKelemen/units v1.2.0
	golang.org/x/image v0.35.0
)

require (
	github.com/SCKelemen/layout v1.1.3
	github.com/SCKelemen/text v1.1.3 // indirect
	github.com/SCKelemen/unicode v1.1.1 // indirect
)

// Exclude problematic test-only dependency (used only in layout tests)
exclude github.com/SCKelemen/wpt-test-gen v0.0.0-00010101000000-000000000000

exclude github.com/SCKelemen/wpt-test-gen v0.0.0-20251213153317-6265321ae2a1

// Use replace for layout to avoid wpt-test-gen issues
replace github.com/SCKelemen/layout => github.com/SCKelemen/layout v1.1.1
