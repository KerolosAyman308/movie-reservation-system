package files

import (
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// Matches any character that is NOT alphanumeric, a hyphen, underscore, or period
	illegalChars = regexp.MustCompile(`[^a-zA-Z0-9.\-_]`)
	// Matches multiple consecutive underscores
	multiUnderscores = regexp.MustCompile(`_{2,}`)
)

func RefineS3Filename(filename string) string {
	// 1. Separate name and extension to avoid refining the dot
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)

	// 2. Convert to lowercase for consistency (optional but recommended)
	name = strings.ToLower(name)

	// 3. Replace spaces and tabs with underscores
	name = strings.ReplaceAll(name, " ", "_")

	// 4. Remove all illegal characters
	name = illegalChars.ReplaceAllString(name, "")

	// 5. Clean up: Remove multiple underscores in a row and trim edges
	name = multiUnderscores.ReplaceAllString(name, "_")
	name = strings.Trim(name, "_-")

	// 6. Fallback for empty names (if the input was only illegal characters)
	if name == "" {
		name = "unnamed_file"
	}

	return name + strings.ToLower(ext)
}
