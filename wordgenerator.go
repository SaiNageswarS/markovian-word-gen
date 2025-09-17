package main

import (
	"strings"
)

// WordGenerator generates strings with character distribution as per input
type WordGenerator struct {
	charDist      map[string]*Sampler[rune]
	contextLength int
}

// NewWordGenerator creates a new WordGenerator
func NewWordGenerator() *WordGenerator {
	return &WordGenerator{
		charDist:      make(map[string]*Sampler[rune]),
		contextLength: 3,
	}
}

// Build builds the character distribution from a list of words
func (wg *WordGenerator) Build(words []string) {
	// Preprocess words: add start and end markers, convert to lowercase
	processedWords := make([]string, len(words))
	for i, word := range words {
		processedWords[i] = "^" + strings.ToLower(word) + "$"
	}

	for _, word := range processedWords {
		wg.addWord(word)
	}
}

// addWord adds a word to the character distribution
func (wg *WordGenerator) addWord(word string) {
	for i := 0; i < len(word)-1; i++ {
		for j := i + 1; j <= i+wg.contextLength && j < len(word); j++ {
			ctx := word[i:j]

			sampler, exists := wg.charDist[ctx]
			if !exists {
				sampler = NewSampler[rune](29) // a-z 26 and ^, $ = 28
				wg.charDist[ctx] = sampler
			}
			sampler.Add(rune(word[j]), 1)
		}
	}
}

// Generate generates a word with the specified maximum length
func (wg *WordGenerator) Generate(maxLength int) string {
	var result strings.Builder
	result.WriteRune('^')
	nextChar := '^'

	for nextChar != '$' && result.Len() < maxLength {
		var ctx string
		if result.Len() < wg.contextLength {
			ctx = result.String()
		} else {
			ctx = result.String()[result.Len()-wg.contextLength:]
		}

		// decrease context till it exists
		for !wg.hasContext(ctx) {
			if len(ctx) <= 1 {
				break
			}
			ctx = ctx[1:]
		}

		sampler := wg.charDist[ctx]
		if sampler == nil {
			break
		}

		sampled := sampler.Sample()
		if sampled == nil {
			break
		}

		nextChar = *sampled
		result.WriteRune(nextChar)
	}

	return result.String()
}

// hasContext checks if the given context exists in the character distribution
func (wg *WordGenerator) hasContext(ctx string) bool {
	_, exists := wg.charDist[ctx]
	return exists
}
