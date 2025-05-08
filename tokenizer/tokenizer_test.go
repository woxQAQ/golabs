package tokenizer

import (
	"reflect"
	"testing"
)

type testcase struct {
	name     string
	input    string
	expected []string
}

func TestBasicSQLSplit(t *testing.T) {
	tests := []testcase{
		{
			name:     "single statement",
			input:    "SELECT * FROM users;",
			expected: []string{"SELECT * FROM users;"},
		},
		{
			name:     "multiple statements",
			input:    "SELECT * FROM users; SELECT * FROM orders;",
			expected: []string{"SELECT * FROM users;", "SELECT * FROM orders;"},
		},
		{
			name:     "with quoted strings",
			input:    `SELECT * FROM users WHERE name = 'John;Doe';`,
			expected: []string{`SELECT * FROM users WHERE name = 'John;Doe';`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			run(t, tt.input, tt)
		})
	}
}

func TestDelimiterChange(t *testing.T) {
	tests := []testcase{
		{
			name:     "change delimiter",
			input:    "DELIMITER //\nCREATE PROCEDURE test() BEGIN SELECT 1; END//\nDELIMITER ;",
			expected: []string{"CREATE PROCEDURE test() BEGIN SELECT 1; END"},
		},
	}

	for _, tt := range tests {
		run(t, tt.input, tt)
	}
}

func TestCommentsHandling(t *testing.T) {
	tests := []testcase{
		{
			name:     "single line comments",
			input:    "SELECT 1; -- comment\nSELECT 2; # another comment",
			expected: []string{"SELECT 1;", "SELECT 2;"},
		},
		{
			name:     "single line comments withoud trailing space",
			input:    "SELECT 1;-- comment\nSELECT 2;# another comment",
			expected: []string{"SELECT 1;", "SELECT 2;"},
		},
		{
			name: "single line comments with delimiter changed",
			input: `DELIMITER //
SELECT 1;-- comment
SELECT 2;# another comment
//`,
			expected: []string{"SELECT 1;-- comment\nSELECT 2;# another comment\n//"},
		},
		{
			name:     "multi line comments",
			input:    "SELECT /* comment */ 1; SELECT /*\nmulti\nline\n*/ 2;",
			expected: []string{"SELECT /* comment */ 1;", "SELECT /*\nmulti\nline\n*/ 2;"},
		},
	}

	for _, tt := range tests {
		run(t, tt.input, tt)
	}
}

func run(t *testing.T, input string, tt testcase) {
	t.Run(tt.name, func(t *testing.T) {
		tokenizer := &Tokenizer{text: tt.input}
		got := tokenizer.tokenize()
		if !reflect.DeepEqual(got, tt.expected) {
			if len(got) != len(tt.expected) {
				t.Errorf("expected %d statements, got %d", len(tt.expected), len(got))
			}
			t.Errorf("expected:\n %#v, got:\n %#v", tt.expected, got)
		}
	})
}
