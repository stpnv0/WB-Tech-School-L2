package unpacker

import (
	"testing"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectedErr bool
	}{
		{
			name:     "basic expansion",
			input:    "a4bc2d5e",
			expected: "aaaabccddddde",
		},
		{
			name:     "no digits",
			input:    "abcd",
			expected: "abcd",
		},
		{
			name:     "zero removes previous rune",
			input:    "ab0c",
			expected: "ac",
		},
		{
			name:     "escape makes digit literal",
			input:    `qwe\4\5`,
			expected: "qwe45",
		},
		{
			name:     "escaped digit can be repeated",
			input:    `qwe\45`,
			expected: "qwe44444",
		},
		{
			name:     "escape then repeat after escape",
			input:    `a\3b2`,
			expected: "a3bb",
		},
		{
			name:     "double escape slash",
			input:    `qwe\\5`,
			expected: `qwe\\\\\`,
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},

		// ==== ошибки ====

		{
			name:        "starts with digit",
			input:       "45",
			expectedErr: true,
		},
		{
			name:        "digit without previous symbol",
			input:       "3abc",
			expectedErr: true,
		},
		{
			name:        "trailing escape",
			input:       `abc\`,
			expectedErr: true,
		},
		{
			name:        "multi-digit number",
			input:       "a12",
			expectedErr: true,
		},
		{
			name:        "digit repeated but no previous symbol after removal",
			input:       "a0 2",
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Unpack(tt.input)

			if tt.expectedErr {
				if err == nil {
					t.Errorf("expected error but got none, output=%q", got)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if got != tt.expected {
				t.Errorf("Unpack(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}
