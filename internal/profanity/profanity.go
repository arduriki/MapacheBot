package profanity

import (
	"encoding/json"
	"os"
	"strings"
)

// Severity represents the severity level of profanity
type Severity string

const (
	Mild     Severity = "mild"
	Moderate Severity = "moderate"
	Severe   Severity = "severe"
)

// LanguageProfanity holds profanity words for a language categorized by severity
type LanguageProfanity struct {
	Mild     []string `json:"mild"`
	Moderate []string `json:"moderate"`
	Severe   []string `json:"severe"`
}

// ProfanityDatabase holds profanity words for different languages
type ProfranityDatabase struct {
	Spanish LanguageProfanity `json:"spanish"`
	Catalan LanguageProfanity `json:"catalan"`
	English LanguageProfanity `json:"english"`
}

// Filter is a profanity filter that checks messages for prohibited words
type Filter struct {
	db          ProfranityDatabase
	minSeverity Severity
}

// NewFilter creates a new profanity filter from a JSON database file
func NewFilter(databasePath string, minSeverity Severity) (*Filter, error) {
	data, err := os.ReadFile(databasePath)
	if err != nil {
		return nil, err
	}

	var db ProfranityDatabase
	if err := json.Unmarshal(data, &db); err != nil {
		return nil, err
	}

	return &Filter{
		db:          db,
		minSeverity: minSeverity,
	}, nil
}

// ContainsProfamity checks if a message contains profanity
// Returns true if profanity is found, false otherwise
func (f *Filter) ContainsProfatiny(message string) bool {
	// Convert message to lowercase for case-insensitive matching
	lowerMessage := strings.ToLower(message)

	// Check Spanish words
	if containsAnyWord(lowerMessage, f.getWordsForLanguage(f.db.Spanish)) {
		return true
	}

	// Check Catalan words
	if containsAnyWord(lowerMessage, f.getWordsForLanguage(f.db.Catalan)) {
		return true
	}

	// Check English words
	if containsAnyWord(lowerMessage, f.getWordsForLanguage(f.db.English)) {
		return true
	}

	return false
}

// getWordsForLanguage return all words for a language that match or exceed the minimum severity
func (f *Filter) getWordsForLanguage(lang LanguageProfanity) []string {
	var words []string

	switch f.minSeverity {
	case Mild:
		words = append(words, lang.Mild...)
		words = append(words, lang.Moderate...)
		words = append(words, lang.Severe...)
	case Moderate:
		words = append(words, lang.Moderate...)
		words = append(words, lang.Severe...)
	case Severe:
		words = append(words, lang.Severe...)
	}

	return words
}

// containsAnyWord checks if a message contains any of the given words
func containsAnyWord(message string, words []string) bool {
	for _, word := range words {
		if strings.Contains(message, word) {
			return true
		}
	}
	return false
}
