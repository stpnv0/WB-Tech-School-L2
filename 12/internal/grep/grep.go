// Package grep implements the text filtering logic.
package grep

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
	"wb-grep/internal/config"
)

type Grep struct {
	cfg    *config.Config
	source io.Reader
}

func New(cfg *config.Config, source io.Reader) *Grep {
	return &Grep{
		cfg:    cfg,
		source: source,
	}
}

func (g *Grep) Run(dst io.Writer) error {
	scanner := bufio.NewScanner(g.source)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	matches, err := g.findMatches(lines)
	if err != nil {
		return err
	}

	if g.cfg.CountOnly {
		count := 0
		for _, m := range matches {
			if m {
				count++
			}
		}

		fmt.Fprintf(dst, "%d\n", count)
		return nil
	}

	toOutput := g.calculateOutputLines(matches, len(lines))

	g.outputLines(dst, lines, toOutput, matches)

	return nil
}

func (g *Grep) findMatches(lines []string) ([]bool, error) {
	matches := make([]bool, len(lines))
	p := g.cfg.Pattern

	var matcher func(string) bool

	if g.cfg.Fixed {
		if g.cfg.IgnoreCase {
			p = strings.ToLower(p)
			matcher = func(s string) bool {
				return strings.Contains(strings.ToLower(s), p)
			}
		} else {
			matcher = func(s string) bool {
				return strings.Contains(s, p)
			}
		}
	} else {
		if g.cfg.IgnoreCase {
			p = "(?i)" + p
		}
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		matcher = func(s string) bool {
			return re.MatchString(s)
		}
	}

	for i, line := range lines {
		matched := matcher(line)
		if g.cfg.Invert {
			matched = !matched
		}
		matches[i] = matched
	}

	return matches, nil
}

func (g *Grep) calculateOutputLines(matches []bool, totalLines int) []bool {
	toOutput := make([]bool, totalLines)

	for i := 0; i < totalLines; i++ {
		if matches[i] {
			toOutput[i] = true

			//before context
			for j := i - g.cfg.Before; j < i; j++ {
				if j >= 0 {
					toOutput[j] = true
				}
			}

			//after context
			for j := i + 1; j <= i+g.cfg.After; j++ {
				if j < totalLines {
					toOutput[j] = true
				}
			}
		}
	}

	return toOutput
}

func (g *Grep) outputLines(dst io.Writer, lines []string, toOutput []bool, matches []bool) {
	lastOutputIndex := -1
	sep := "--"
	useSep := g.cfg.After > 0 || g.cfg.Before > 0

	for i := 0; i < len(lines); i++ {
		if !toOutput[i] {
			continue
		}

		if useSep && lastOutputIndex != -1 && i > lastOutputIndex+1 {
			fmt.Fprintln(dst, sep)
		}

		prefix := ""
		if g.cfg.LineNum {
			if matches[i] {
				prefix = fmt.Sprintf("%d:", i+1)
			} else {
				prefix = fmt.Sprintf("%d-", i+1)
			}
		}

		fmt.Fprintf(dst, "%s%s\n", prefix, lines[i])
		lastOutputIndex = i
	}
}
