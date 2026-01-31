package cut

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"wb-cut/internal/config"
)

type Cut struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Cut {
	return &Cut{cfg: cfg}
}

func (c *Cut) Run(source io.Reader, dest io.Writer) error {
	fieldMap, err := parseFields(c.cfg.Fields)
	if err != nil {
		return err
	}

	delimiter := []byte(c.cfg.Delimiter)
	scanner := bufio.NewScanner(source)

	for scanner.Scan() {
		line := scanner.Bytes()

		if !bytes.Contains(line, delimiter) {
			if c.cfg.Separated {
				continue
			}
			dest.Write(line)
			dest.Write([]byte{'\n'})
			continue
		}

		parts := bytes.Split(line, delimiter)
		var outParts [][]byte

		for i, part := range parts {
			if fieldMap[i+1] {
				outParts = append(outParts, part)
			}
		}

		dest.Write(bytes.Join(outParts, delimiter))
		dest.Write([]byte{'\n'})
	}

	return scanner.Err()
}

func parseFields(fieldStr string) (map[int]bool, error) {
	res := make(map[int]bool)
	parts := strings.Split(fieldStr, ",")
	for _, p := range parts {
		if strings.Contains(p, "-") {
			// Range
			ranges := strings.Split(p, "-")
			if len(ranges) != 2 {
				return nil, fmt.Errorf("invalid range format: %s", p)
			}

			startStr, endStr := ranges[0], ranges[1]

			start, err1 := strconv.Atoi(startStr)
			end, err2 := strconv.Atoi(endStr)

			if err1 != nil || err2 != nil {
				return nil, fmt.Errorf("invalid numbers in range: %s", p)
			}

			if start > end {
				return nil, fmt.Errorf("invalid decreasing range: %s", p)
			}

			for i := start; i <= end; i++ {
				res[i] = true
			}
		} else {
			// Single number
			val, err := strconv.Atoi(p)
			if err != nil {
				return nil, fmt.Errorf("invalid field number: %s", p)
			}
			res[val] = true
		}
	}
	return res, nil
}
