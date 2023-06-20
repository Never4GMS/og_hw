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

type wordStat struct {
	word  string
	count uint16
}

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

func createRating(counters map[string]uint16) []wordStat {
	rating := make([]wordStat, 0, len(counters))
	for word, count := range counters {
		rating = append(rating, wordStat{
			word:  word,
			count: count,
		})
	}
	sort.Slice(rating, func(i, j int) bool {
		if rating[i].count == rating[j].count {
			return strings.Compare(rating[i].word, rating[j].word) < 0
		}

		return rating[i].count > rating[j].count
	})

	return rating
}

func top10(rating []wordStat) []string {
	var limit int = 10
	if len(rating) < 10 {
		limit = len(rating)
	}

	top := make([]string, 0, limit)

	for i := 0; i < limit; i++ {
		top = append(top, rating[i].word)
	}

	return top
}
