package lines

import (
	"fmt"
	"strings"
)

// Sprintf returns fmt.Sprintf(format, a...) split on newlines.
func Sprintf(format string, a ...any) []string {
	return strings.Split(fmt.Sprintf(format, a...), "\n")
}

// Indent returns the lines with the given prefix prepended. If skip is
// positive, then that many lines will be skipped before indenting the
// remainder.
func Indent(lines []string, prefix string, skip int) (indented []string) {
	indented = make([]string, len(lines))

	for i, line := range lines {
		if skip <= 0 {
			indented[i] = prefix + line
			continue
		}

		indented[i] = line
		skip--
	}

	return indented
}

// TrimTrailing removes trailing lines that are empty after applying strings.Trim.
func TrimTrailing(lines []string, cutset string) (trimmed []string) {
	trimmed = make([]string, 0, len(lines))

	for i := len(lines) - 1; i >= 0; i-- {
		if len(strings.Trim(lines[i], cutset)) == 0 {
			continue
		}

		trimmed = append(trimmed, lines[:i+1]...)
		break
	}

	return trimmed
}
