package parser

import (
	"errors"
	"io"
	"strings"
	"testing"
	"wb-sort/internal/config"
)

func TestParser_Parse(t *testing.T) {
	cfg := &config.Config{}

	tests := []struct {
		name    string
		source  io.Reader
		want    []string
		wantErr error
	}{
		{
			name:   "Normal reading from source",
			source: strings.NewReader("line1\nline2\nline3\n"),
			want:   []string{"line1", "line2", "line3"},
		},
		{
			name:    "Source is nil",
			source:  nil,
			wantErr: ErrSource,
		},
		{
			name:   "Empty source",
			source: strings.NewReader(""),
			want:   []string{},
		},
		{
			name:   "Source with trailing newline",
			source: strings.NewReader("line1\nline2\n"),
			want:   []string{"line1", "line2"},
		},
		{
			name:    "Scanner error",
			source:  &errorReader{err: errors.New("mock scanner error")},
			wantErr: errors.New("mock scanner error"),
		},
		{
			name:   "Large input",
			source: strings.NewReader(strings.Repeat("long line\n", 1000)),
			want:   generateLargeInput(1000),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(cfg, tt.source)
			got, err := p.Parse()
			if tt.wantErr != nil {
				if err == nil || err.Error() != tt.wantErr.Error() {
					t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("Parse() unexpected error = %v", err)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("Parse() got len = %d, want %d", len(got), len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("Parse() got[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

type errorReader struct {
	err error
}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, e.err
}

func generateLargeInput(n int) []string {
	var res []string
	for i := 0; i < n; i++ {
		res = append(res, "long line")
	}
	return res
}
