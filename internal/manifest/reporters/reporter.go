package reporters

import (
	"fmt"

	"github.com/thisismeamir/hepsw/internal/manifest"
)

// ReportFormat specifies the output format for reports
type ReportFormat string

const (
	FormatText     ReportFormat = "text"
	FormatMarkdown ReportFormat = "markdown"
	FormatJSON     ReportFormat = "json"
)

// GenerateReport creates a comprehensive report of a manifest
func GenerateReport(m *manifest.Manifest, format ReportFormat) (string, error) {
	switch format {
	case FormatText:
		return generateTextReport(m), nil
	case FormatMarkdown:
		return generateMarkdownReport(m), nil
	case FormatJSON:
		return generateJSONReport(m)
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}
