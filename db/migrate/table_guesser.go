package migrate

import "regexp"

var (
	CreatePatterns []*regexp.Regexp
	ChangePatterns []*regexp.Regexp
)

type TableGuesser struct {
}

func NewTableGuesser() *TableGuesser {
	CreatePatterns = append(CreatePatterns,
		regexp.MustCompile(`^create_(\w+)_table$`),
		regexp.MustCompile(`^create_(\w+)$`),
	)
	ChangePatterns = append(ChangePatterns,
		regexp.MustCompile(`_(to|from|in)_(\w+)_table$`),
		regexp.MustCompile(`_(to|from|in)_(\w+)$`),
	)

	return &TableGuesser{}
}

func (guesser *TableGuesser) Guess(migration string) (string, bool) {
	for _, pattern := range CreatePatterns {
		if match := pattern.FindStringSubmatch(migration); len(match) > 0 {
			return match[1], true
		}
	}

	for _, pattern := range ChangePatterns {
		if match := pattern.FindStringSubmatch(migration); len(match) > 0 {
			return match[2], false
		}
	}

	return "", false
}
