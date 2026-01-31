package sorter

import (
	"sort"
	"strconv"
	"strings"
	"wb-sort/internal/config"
)

// ParserInterface defines the interface for parsing input data.
type ParserInterface interface {
	Parse() ([]string, error)
}

// Sorter sorts lines according to the given configuration.
type Sorter struct {
	cfg    *config.Config
	parser ParserInterface
}

// NewSorter creates a new Sorter with the given configuration and parser.
func NewSorter(cfg *config.Config, p ParserInterface) *Sorter {
	return &Sorter{
		cfg:    cfg,
		parser: p,
	}
}

// Sort sorts the lines according to the configuration and returns the result.
func (s *Sorter) Sort() ([]string, error) {
	lines, err := s.parser.Parse()
	if err != nil {
		return nil, err
	}
	if len(lines) == 0 {
		return lines, nil
	}
	k := s.cfg.Column - 1
	if k < 0 {
		k = 0
	}
	sort.SliceStable(lines, func(i, j int) bool {
		cmp := s.compare(lines[i], lines[j], k)
		if s.cfg.IsReverse {
			return cmp > 0
		}
		return cmp < 0
	})

	if s.cfg.IsUnique {
		lines = uniqSort(lines)
	}

	return lines, nil
}

func (s *Sorter) compare(a, b string, k int) int {
	acol := ""
	bcol := ""
	as := strings.Split(a, "\t")
	bs := strings.Split(b, "\t")
	if k < len(as) {
		acol = as[k]
	}
	if k < len(bs) {
		bcol = bs[k]
	}

	if s.cfg.IsNumeric {
		af, aerr := strconv.ParseFloat(acol, 64)
		bf, berr := strconv.ParseFloat(bcol, 64)
		if aerr == nil && berr == nil {
			if af < bf {
				return -1
			}
			if af > bf {
				return 1
			}
			return 0
		} else if aerr != nil && berr != nil {
			if acol < bcol {
				return -1
			}
			if acol > bcol {
				return 1
			}
			return 0
		} else if aerr != nil {
			if 0 < bf {
				return -1
			}
			if 0 > bf {
				return 1
			}
			return 0
		} else {
			if af < 0 {
				return -1
			}
			if af > 0 {
				return 1
			}
			return 0
		}
	}
	// Fallback string
	if acol < bcol {
		return -1
	}
	if acol > bcol {
		return 1
	}
	return 0
}

func uniqSort(s []string) []string {
	if len(s) == 0 {
		return nil
	}

	j := 0
	for i := 1; i < len(s); i++ {
		if s[i] != s[j] {
			j++
			s[j] = s[i]
		}
	}

	return s[:j+1]
}
