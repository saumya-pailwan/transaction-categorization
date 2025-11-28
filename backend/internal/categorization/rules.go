package categorization

import (
	"context"
	"regexp"
	"strings"

	"autonomoustx/internal/db"
)

func MatchRule(description string) (string, bool) {
	ctx := context.Background()

	// Fetch all rules (in a real app, cache this)
	rows, err := db.Pool.Query(ctx, "SELECT pattern, category FROM rules ORDER BY priority DESC")
	if err != nil {
		return "", false
	}
	defer rows.Close()

	for rows.Next() {
		var pattern, category string
		if err := rows.Scan(&pattern, &category); err != nil {
			continue
		}

		// Simple contains check or regex
		// assumption: simple case-insensitive contains for simplicity,
		// or if it starts with "regex:", treat as regex.

		if strings.HasPrefix(pattern, "regex:") {
			regexPattern := strings.TrimPrefix(pattern, "regex:")
			matched, _ := regexp.MatchString(regexPattern, description)
			if matched {
				return category, true
			}
		} else {
			if strings.Contains(strings.ToLower(description), strings.ToLower(pattern)) {
				return category, true
			}
		}
	}

	return "", false
}
