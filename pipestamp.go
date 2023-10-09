package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

// ParseISO8601Timestamp tries to parse the given timestamp string in ISO 8601 format
func ParseISO8601Timestamp(s string) (time.Time, error) {
	// List of ISO 8601 variants that we'll attempt to parse
	var layouts = []string{
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05.999999Z",
		"2006-01-02T15:04:05.999999-07:00",
		"2006-01-02T15:04:05-07:00",
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, s)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("couldn't parse timestamp: %s", s)
}

// TimeAgo converts a given time to a relative "time ago" format
func TimeAgo(t time.Time) string {
	seconds := time.Since(t).Seconds()

	switch {
	case seconds < 60:
		return fmt.Sprintf("%.0f seconds ago", seconds)
	case seconds < 3600:
		return fmt.Sprintf("%.0f minutes ago", seconds/60)
	case seconds < 86400:
		return fmt.Sprintf("%.0f hours ago", seconds/3600)
	case seconds < 2592000: // 30 days
		return fmt.Sprintf("%.0f days ago", seconds/86400)
	default:
		return fmt.Sprintf("%.0f months ago", seconds/2592000)
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// Regex to identify ISO 8601-like timestamps
	timestampPattern := `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d{1,6})?(?:Z|[-+]\d{2}:\d{2})`
	regex := regexp.MustCompile(timestampPattern)

	for scanner.Scan() {
		line := scanner.Text()

		matches := regex.FindAllString(line, -1)

		for _, match := range matches {
			parsedTime, err := ParseISO8601Timestamp(match)
			if err != nil {
				// Couldn't parse the timestamp; continue with the next match
				continue
			}

			line = strings.Replace(line, match, TimeAgo(parsedTime), 1)
		}

		fmt.Println(line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
	}
}
