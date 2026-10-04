package build

import (
	"bufio"
	"strings"
)

// MaxHeadlineLines caps the structured failure summary shown to the user.
const MaxHeadlineLines = 15

// ExtractHeadlines scans a full Gradle log and returns high-signal failure lines.
func ExtractHeadlines(log string) []string {
	lines := splitLines(log)
	var out []string
	seen := map[string]bool{}

	add := func(s string) {
		s = strings.TrimRight(s, "\r\n")
		s = strings.TrimSpace(s)
		if s == "" || seen[s] {
			return
		}
		seen[s] = true
		out = append(out, s)
	}

	// Pass 1: high-priority multi-line blocks
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		switch {
		case strings.Contains(trimmed, "* What went wrong:"):
			add(trimmed)
			for j := i + 1; j < len(lines) && len(out) < MaxHeadlineLines; j++ {
				t := strings.TrimSpace(lines[j])
				if t == "" || strings.HasPrefix(t, "* Try:") || strings.HasPrefix(t, "BUILD ") {
					break
				}
				add(t)
			}
		case strings.HasPrefix(trimmed, "FAILURE: Build failed"):
			add(trimmed)
		case strings.HasPrefix(trimmed, "Execution failed for task"):
			add(trimmed)
		}
	}

	// Pass 2: compiler / CMake / NDK signal lines
	for _, line := range lines {
		if len(out) >= MaxHeadlineLines {
			break
		}
		trimmed := strings.TrimSpace(line)
		lower := strings.ToLower(trimmed)

		switch {
		case strings.HasPrefix(trimmed, "e: file:///"):
			add(trimmed)
		case strings.Contains(trimmed, ".kt:") && (strings.Contains(lower, "error") || strings.HasPrefix(trimmed, "e:")):
			add(trimmed)
		case strings.Contains(trimmed, ".java:") && strings.Contains(lower, "error"):
			add(trimmed)
		case strings.Contains(trimmed, ".xml:") && strings.Contains(lower, "error"):
			add(trimmed)
		case strings.Contains(trimmed, "CMake Error"):
			add(trimmed)
		case strings.Contains(trimmed, "ninja: error"):
			add(trimmed)
		case strings.HasPrefix(trimmed, "error:") || strings.HasPrefix(trimmed, "Error:"):
			if strings.Contains(trimmed, ".kt:") || strings.Contains(trimmed, ".java:") ||
				strings.Contains(trimmed, ".xml:") || strings.Contains(trimmed, "file:///") {
				add(trimmed)
			}
		}
	}

	// Pass 3: if still empty, take last non-empty lines that look like errors
	if len(out) == 0 {
		for i := len(lines) - 1; i >= 0 && len(out) < 8; i-- {
			t := strings.TrimSpace(lines[i])
			if t == "" {
				continue
			}
			lower := strings.ToLower(t)
			if strings.Contains(lower, "failed") || strings.Contains(lower, "error") ||
				strings.Contains(lower, "exception") {
				out = append([]string{t}, out...)
			}
		}
	}

	if len(out) > MaxHeadlineLines {
		out = out[:MaxHeadlineLines]
	}
	return out
}

func splitLines(s string) []string {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(s))
	buf := make([]byte, 0, 64*1024)
	sc.Buffer(buf, 1024*1024)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	if len(lines) == 0 && s != "" {
		return strings.Split(s, "\n")
	}
	return lines
}
