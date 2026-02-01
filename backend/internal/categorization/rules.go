package categorization

import (
	"context"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	"autonomoustx/internal/db"
)

type CachedRule struct {
	Pattern  string
	Category string
	Regex    *regexp.Regexp
	IsRegex  bool
}

var (
	ruleCache []CachedRule
	ruleMutex sync.RWMutex
)

// LoadRules fetches all rules from the DB and caches them in memory.
func LoadRules() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.Pool.Query(ctx, "SELECT pattern, category FROM rules ORDER BY priority DESC")
	if err != nil {
		return err
	}
	defer rows.Close()

	var newRules []CachedRule

	for rows.Next() {
		var pattern, category string
		if err := rows.Scan(&pattern, &category); err != nil {
			log.Printf("Error scanning rule: %v", err)
			continue
		}

		rule := CachedRule{
			Pattern:  pattern,
			Category: category,
		}

		if strings.HasPrefix(pattern, "regex:") {
			rule.IsRegex = true
			regexPat := strings.TrimPrefix(pattern, "regex:")
			// Pre-compile regex
			re, err := regexp.Compile(regexPat)
			if err != nil {
				log.Printf("Error compiling regex rule '%s': %v", pattern, err)
				continue
			}
			rule.Regex = re
		} else {
			// Pre-lowercase for case-insensitive check
			rule.Pattern = strings.ToLower(pattern)
		}

		newRules = append(newRules, rule)
	}

	ruleMutex.Lock()
	ruleCache = newRules
	ruleMutex.Unlock()

	return nil
}

func MatchRule(description string) (string, bool) {
	ruleMutex.RLock()
	defer ruleMutex.RUnlock()

	descLower := strings.ToLower(description)

	for _, rule := range ruleCache {
		if rule.IsRegex {
			if rule.Regex.MatchString(description) {
				return rule.Category, true
			}
		} else {
			if strings.Contains(descLower, rule.Pattern) {
				return rule.Category, true
			}
		}
	}

	return "", false
}
