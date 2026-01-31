package main

import (
	"fmt"
	"sort"
	"strings"
)

type Group struct {
	firstW string
	words  []string
}

func findAnagrams(data []string) map[string][]string {
	groups := make(map[string]*Group)

	for _, w := range data {
		word := strings.ToLower(w)
		key := getSorted(word)

		if _, ok := groups[key]; !ok {
			groups[key] = &Group{
				firstW: word,
			}
		}

		groups[key].words = append(groups[key].words, word)

	}

	res := make(map[string][]string, len(groups))
	for _, g := range groups {
		if len(g.words) < 2 {
			continue
		}

		sort.Strings(g.words)
		res[g.firstW] = g.words
	}

	return res
}

func getSorted(word string) string {
	runes := []rune(word)
	sort.Slice(runes, func(i, j int) bool {
		return runes[i] < runes[j]
	})
	return string(runes)
}

func main() {
	input := []string{
		"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол",
	}

	res := findAnagrams(input)

	for k, v := range res {
		fmt.Printf("-%q: %v\n", k, v)
	}
}
