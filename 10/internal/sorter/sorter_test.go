package sorter

import (
	"errors"
	"reflect"
	"testing"
	"wb-sort/internal/config"
)

type mockParser struct {
	lines []string
	err   error
}

func (m *mockParser) Parse() ([]string, error) {
	return m.lines, m.err
}

func TestSorter_Sort(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config.Config
		input   []string
		want    []string
		wantErr error
	}{
		{
			name:  "Basic lexical sort (no flags)",
			cfg:   &config.Config{Column: 1},
			input: []string{"c", "a", "b"},
			want:  []string{"a", "b", "c"},
		},
		{
			name:  "Sort by column",
			cfg:   &config.Config{Column: 2},
			input: []string{"apple\tc", "banana\ta", "cherry\tb"},
			want:  []string{"banana\ta", "cherry\tb", "apple\tc"},
		},
		{
			name:  "Numeric sort by column",
			cfg:   &config.Config{Column: 2, IsNumeric: true},
			input: []string{"apple\t3", "banana\t1", "cherry\t2"},
			want:  []string{"banana\t1", "cherry\t2", "apple\t3"},
		},
		{
			name:  "Reverse numeric",
			cfg:   &config.Config{Column: 2, IsNumeric: true, IsReverse: true},
			input: []string{"apple\t3", "banana\t1", "cherry\t2"},
			want:  []string{"apple\t3", "cherry\t2", "banana\t1"},
		},
		{
			name:  "Unique after sort",
			cfg:   &config.Config{IsUnique: true},
			input: []string{"a", "b", "a", "c", "b"},
			want:  []string{"a", "b", "c"},
		},
		{
			name:  "Column out of range",
			cfg:   &config.Config{Column: 3},
			input: []string{"a\tb", "c\td"},
			want:  []string{"a\tb", "c\td"},
		},
		{
			name:  "Invalid numeric (fallback to string)",
			cfg:   &config.Config{IsNumeric: true},
			input: []string{"10", "2", "abc", "3"},
			want:  []string{"abc", "2", "3", "10"},
		},
		{
			name:  "Empty input",
			cfg:   &config.Config{},
			input: []string{},
			want:  []string{},
		},
		{
			name:  "Single line",
			cfg:   &config.Config{},
			input: []string{"only one"},
			want:  []string{"only one"},
		},
		{
			name:    "Parser error",
			cfg:     &config.Config{},
			input:   nil,
			wantErr: errors.New("mock parser error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mp := &mockParser{lines: tt.input, err: tt.wantErr}
			s := NewSorter(tt.cfg, mp)
			got, err := s.Sort()
			if tt.wantErr != nil {
				if err == nil || err.Error() != tt.wantErr.Error() {
					t.Errorf("Sort() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("Sort() unexpected error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sort() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compare(t *testing.T) {
	s := &Sorter{cfg: &config.Config{}}
	tests := []struct {
		name string
		a, b string
		k    int
		want int
	}{
		{"Lexical equal", "a", "a", 0, 0},
		{"Lexical less", "a", "b", 0, -1},
		{"Numeric less", "1", "2", 0, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := s.compare(tt.a, tt.b, tt.k); got != tt.want {
				t.Errorf("compare() = %d, want %d", got, tt.want)
			}
		})
	}
}

func Test_uniqSort(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{"No dups", []string{"a", "b", "c"}, []string{"a", "b", "c"}},
		{"With dups", []string{"a", "a", "b", "c", "c"}, []string{"a", "b", "c"}},
		{"All dups", []string{"a", "a", "a"}, []string{"a"}},
		{"Empty", []string{}, nil},
		{"Single", []string{"a"}, []string{"a"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := uniqSort(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("uniqSort() = %v, want %v", got, tt.want)
			}
		})
	}
}
