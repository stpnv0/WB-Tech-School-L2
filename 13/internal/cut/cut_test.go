package cut

import (
	"bytes"
	"strings"
	"testing"
	"wb-cut/internal/config"
)

func TestCut_Run(t *testing.T) {
	tests := []struct {
		name     string
		cfg      config.Config
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "Basic field 1",
			cfg:      config.Config{Fields: "1", Delimiter: "\t", Separated: false},
			input:    "a\tb\tc\n1\t2\t3",
			expected: "a\n1\n",
		},
		{
			name:     "Fields 1 and 3",
			cfg:      config.Config{Fields: "1,3", Delimiter: "\t", Separated: false},
			input:    "a\tb\tc\n1\t2\t3",
			expected: "a\tc\n1\t3\n",
		},
		{
			name:     "Range 2-3",
			cfg:      config.Config{Fields: "2-3", Delimiter: "\t", Separated: false},
			input:    "a\tb\tc\n1\t2\t3",
			expected: "b\tc\n2\t3\n",
		},
		{
			name:     "Mixed 1,3-4 (out of bound field 4 ignored)",
			cfg:      config.Config{Fields: "1,3-4", Delimiter: "\t", Separated: false},
			input:    "a\tb\tc\n1\t2\t3",
			expected: "a\tc\n1\t3\n",
		},
		{
			name:     "Custom delimiter :",
			cfg:      config.Config{Fields: "2", Delimiter: ":", Separated: false},
			input:    "a:b:c\n1:2:3",
			expected: "b\n2\n",
		},
		{
			name:     "Separated flag - skip no delimiter",
			cfg:      config.Config{Fields: "1", Delimiter: ":", Separated: true},
			input:    "no_delimiter\nwith:delimiter",
			expected: "with\n",
		},
		{
			name:     "No Separated flag - print no delimiter line as is",
			cfg:      config.Config{Fields: "1", Delimiter: ":", Separated: false},
			input:    "no_delimiter\nwith:delimiter",
			expected: "no_delimiter\nwith\n",
		},
		{
			name:    "Invalid Range",
			cfg:     config.Config{Fields: "3-1", Delimiter: "\t"},
			input:   "a\tb",
			wantErr: true,
		},
		{
			name:    "Invalid Field",
			cfg:     config.Config{Fields: "a", Delimiter: "\t"},
			input:   "a\tb",
			wantErr: true,
		},
		{
			name:     "Complex combination",
			cfg:      config.Config{Fields: "1,4", Delimiter: ",", Separated: false},
			input:    "1,2,3,4,5\na,b,c",
			expected: "1,4\na\n", // Line 2 has only 3 parts, so field 4 is ignored -> "a"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cut := New(&tt.cfg)
			src := strings.NewReader(tt.input)
			dst := &bytes.Buffer{}
			err := cut.Run(src, dst)

			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got := dst.String(); got != tt.expected {
					t.Errorf("Run() got = %q, want %q", got, tt.expected)
				}
			}
		})
	}
}
