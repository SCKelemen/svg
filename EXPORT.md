# SVG Export

The svg library now includes native export functionality to convert SVG to raster formats (PNG, JPEG) using only standard Go libraries.

## Features

- **No external dependencies**: Uses only `golang.org/x/image` and standard library
- **Multiple formats**: SVG (passthrough), PNG, and JPEG
- **Configurable**: Width, height, quality, and DPI settings
- **DPI-aware units**: Supports physical SVG units (`in`, `cm`, `mm`, `pt`, `pc`, `q`)
- **Safe unsupported handling**: Returns an error for unsupported renderable elements by default (configurable)
- **Basic shape support**: Rectangles, circles, and lines with fill and stroke colors

## Usage

### Basic Export

```go
import "github.com/SCKelemen/svg"

// Generate or load SVG data
svgData := `<svg width="200" height="200">
    <rect x="50" y="50" width="100" height="100" fill="#ff0000"/>
</svg>`

// Export to PNG
opts := svg.ExportOptions{
    Format: svg.FormatPNG,
    Width:  200,
    Height: 200,
}

pngData, err := svg.Export(svgData, opts)
if err != nil {
    log.Fatal(err)
}

// Write to file
os.WriteFile("output.png", pngData, 0644)
```

### Export Formats

```go
// SVG (passthrough - just returns the SVG data)
opts := svg.ExportOptions{
    Format: svg.FormatSVG,
}

// PNG
opts := svg.ExportOptions{
    Format: svg.FormatPNG,
    Width:  800,
    Height: 600,
}

// JPEG with quality setting
opts := svg.ExportOptions{
    Format:  svg.FormatJPEG,
    Width:   800,
    Height:  600,
    Quality: 90, // 0-100, default 90
}

// Ignore unsupported renderable elements (e.g. <text>, <path>) instead of failing
opts := svg.ExportOptions{
    Format:            svg.FormatPNG,
    IgnoreUnsupported: true,
}
```

### Default Options

```go
opts := svg.DefaultExportOptions()
// Returns:
// - Format: FormatSVG
// - Quality: 90
// - DPI: 96
```

### Helper Functions

```go
// Get MIME type for HTTP responses
mimeType := svg.GetMimeType(svg.FormatPNG) // "image/png"

// Get file extension
ext := svg.GetFileExtension(svg.FormatPNG) // ".png"

// Parse format from string
format, err := svg.ParseFormat("jpeg") // FormatJPEG
```

## Supported SVG Elements

The current implementation supports basic SVG shapes:

- ✅ `<rect>` - Rectangles with fill color
- ✅ `<circle>` - Circles with fill color and antialiasing
- ✅ `<line>` - Lines with stroke color
- ✅ `<g>` - Groups (renders children)
- ✅ Color parsing: hex colors (`#RGB`, `#RRGGBB`), named colors (red, blue, etc.)
- ❌ `<text>` - Not yet implemented (requires font support)
- ❌ `<path>` - Not yet implemented (requires path parsing)

## Implementation Details

### Architecture

The export system consists of three main components:

1. **SVG Parser** (`parseSVG`): Uses `encoding/xml` to parse SVG into a tree structure
2. **Rasterizer** (`rasterize`): Uses `golang.org/x/image/vector` for antialiased rendering
3. **Encoders**: Uses standard `image/png` and `image/jpeg` encoders

### Color Support

- Hex colors: `#RGB`, `#RRGGBB`
- Named colors: `white`, `black`, `red`, `green`, `blue`
- Transparent: `none` or empty string

### Rendering Strategy

- PNG exports preserve transparency by default
- JPEG exports use a white background
- Antialiased circles using 32-segment approximation
- Rectangles rendered directly to image
- Lines use a width-aware pixel brush stroke renderer

## Limitations

1. **Text rendering**: Not yet implemented (requires font support from `golang.org/x/image/font`)
2. **Path elements**: Not yet implemented (requires SVG path parser)
3. **Transforms**: Not yet supported (translate, rotate, scale)
4. **Gradients**: Not yet supported
5. **Advanced features**: Filters, masks, patterns not supported

## Future Enhancements

- [ ] Text rendering with font support
- [ ] SVG path parsing and rendering
- [ ] Transform support (translate, rotate, scale)
- [ ] Gradient fills (linear, radial)
- [ ] Stroke width and dash arrays
- [ ] Opacity and blend modes
- [ ] Advanced shapes (ellipse, polygon, polyline)

## Performance

The native implementation provides good performance for simple shapes:

- SVG passthrough: O(1) - no processing
- PNG export: O(n) where n = number of shapes
- JPEG export: Similar to PNG with compression overhead

For complex SVGs with many elements, performance is limited by the vector rasterizer.

## Examples

### Export Chart to PNG

```go
// Generate an SVG string (from RenderToSVG or hand-built markup)
svgData := `<svg width="800" height="600">
  <rect x="0" y="0" width="800" height="600" fill="#ffffff"/>
</svg>`

// Export to PNG
opts := svg.ExportOptions{
    Format: svg.FormatPNG,
    Width:  800,
    Height: 600,
}

pngData, _ := svg.Export(svgData, opts)
os.WriteFile("chart.png", pngData, 0644)
```

### HTTP Response

```go
http.HandleFunc("/chart.png", func(w http.ResponseWriter, r *http.Request) {
    svgData := generateChart()

    opts := svg.ExportOptions{
        Format: svg.FormatPNG,
        Width:  800,
        Height: 600,
    }

    pngData, err := svg.Export(svgData, opts)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    w.Header().Set("Content-Type", svg.GetMimeType(svg.FormatPNG))
    w.Write(pngData)
})
```

## Testing

Run the export tests:

```bash
go test -v -run TestExport
```

All tests should pass, verifying:
- SVG passthrough
- PNG export with signature validation
- JPEG export with signature validation
- Circle rendering
- Format parsing and utility functions
