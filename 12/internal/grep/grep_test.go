package grep

import (
	"bytes"
	"strings"
	"testing"
	"wb-grep/internal/config"
)

func TestGrep_Run(t *testing.T) {
	tests := []struct {
		name     string
		cfg      config.Config
		input    string
		expected string
	}{
		{
			name: "Basic match",
			cfg: config.Config{
				Pattern: "apple",
			},
			input:    "apple\nbanana\napple pie",
			expected: "apple\napple pie\n",
		},
		{
			name: "Fixed string match (-F)",
			cfg: config.Config{
				Pattern: ".",
				Fixed:   true,
			},
			input:    "apple\n.\nbanana",
			expected: ".\n",
		},
		{
			name: "Regex match",
			cfg: config.Config{
				Pattern: "a.p",
			},
			input:    "apple\naxp\nbanana",
			expected: "apple\naxp\n",
		},
		{
			name: "Ignore case (-i)",
			cfg: config.Config{
				Pattern:    "APPLE",
				IgnoreCase: true,
			},
			input:    "apple\nBanana\nApple Pie",
			expected: "apple\nApple Pie\n",
		},
		{
			name: "Invert match (-v)",
			cfg: config.Config{
				Pattern: "apple",
				Invert:  true,
			},
			input:    "apple\nbanana\napple pie",
			expected: "banana\n",
		},
		{
			name: "Count only (-c)",
			cfg: config.Config{
				Pattern:   "apple",
				CountOnly: true,
			},
			input:    "apple\nbanana\napple pie",
			expected: "2\n",
		},
		{
			name: "Line numbers (-n)",
			cfg: config.Config{
				Pattern: "apple",
				LineNum: true,
			},
			input:    "apple\nbanana\napple pie",
			expected: "1:apple\n3:apple pie\n",
		},
		{
			name: "After context (-A)",
			cfg: config.Config{
				Pattern: "apple",
				After:   1,
			},
			input:    "apple\nbanana\ncherry\napple\ndate",
			expected: "apple\nbanana\n--\napple\ndate\n",
		},
		{
			name: "Before context (-B)",
			cfg: config.Config{
				Pattern: "apple",
				Before:  1,
			},
			input:    "banana\napple\ncherry\ndate\napple",
			expected: "banana\napple\n--\ndate\napple\n",
		},
		{
			name: "Context (-C) / (-A -B)",
			cfg: config.Config{
				Pattern: "apple",
				After:   1,
				Before:  1,
			},
			input:    "1\napple\n2\n3\napple\n4",
			expected: "1\napple\n2\n3\napple\n4\n",
		},
		{
			name: "Context with line numbers",
			cfg: config.Config{
				Pattern: "apple",
				After:   1,
				Before:  1,
				LineNum: true,
			},
			input:    "hue\napple\njack",
			expected: "1-hue\n2:apple\n3-jack\n",
		},
		{
			name: "Overlapping context",
			cfg: config.Config{
				Pattern: "A",
				After:   1,
				Before:  1,
			},
			input:    "1\nA\n2\nA\n3",
			expected: "1\nA\n2\nA\n3\n",
		},
		{
			name: "No separator without context",
			cfg: config.Config{
				Pattern: "A",
				After:   0,
				Before:  0,
			},
			input:    "A\nB\nC\nA",
			expected: "A\nA\n",
		},
		{
			name: "Separator with context",
			cfg: config.Config{
				Pattern: "A",
				After:   1,
				Before:  1,
			},
			input:    "A\nB\nC\nD\nA",
			expected: "A\nB\n--\nD\nA\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := strings.NewReader(tt.input)
			grep := New(&tt.cfg, src)

			var dst bytes.Buffer
			err := grep.Run(&dst)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if dst.String() != tt.expected {
				t.Errorf("Expected output:\n%q\nGot:\n%q", tt.expected, dst.String())
			}
		})
	}
}
