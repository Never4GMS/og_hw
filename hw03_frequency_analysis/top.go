package hw03frequencyanalysis

import (
	"sort"
	"strings"
	"unicode"
)

const (
	hyphen       rune   = '-'
	hyphenString string = "-"
)

func Top10(text string) []string {
	return top10(createRating(createWordMap(text)))
}

func createWordMap(text string) map[string]uint16 {
	counters := map[string]uint16{}
	sb := strings.Builder{}

	for _, symbol := range text {
		if unicode.IsLetter(symbol) || symbol == hyphen {
			sb.WriteRune(symbol)
		} else if sb.Len() > 0 {
			word := sb.String()

			if word != hyphenString {
				counters[strings.ToLower(sb.String())]++
			}

			sb.Reset()
		}
	}

	return counters
}

func createRating(counters map[string]uint16) []string {
	rating := make([]string, 0, len(counters))
	for word := range counters {
		rating = append(rating, word)
	}
	sort.Slice(rating, func(i, j int) bool {
		left := rating[i]
		leftCount := counters[left]
		right := rating[j]
		rightCount := counters[right]

		if leftCount == rightCount {
			return strings.Compare(left, right) < 0
		}

		return leftCount > rightCount
	})

	return rating
}

func top10(rating []string) []string {
	limit := 10
	if len(rating) < 10 {
		limit = len(rating)
	}

	return rating[:limit]
}
