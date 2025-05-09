// package tokenizer is designed to tokenize SQL statements for split multiple statements in slice.

// currently, we support mysql syntax
// for mysql, we can easily split by delimiter ";"
// but mysql also alow user to declare a new delimiter, like "//", to support procedure or function.

package tokenizer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Tokenizer struct {
	text   string
	cursor int
	lines  int
}

const (
	KEYWORD_DELIMITER = "delimiter"
	DEFAULT_DELIMITER = ";"
)

func (t *Tokenizer) tokenize() []string {
	delimiter := DEFAULT_DELIMITER
	res := []string{}

	// tokenPos point to the begin pos of last unparsed token in input rune array
	tokenPos := 0
	// stmtPos point to the begin pos of last unparsed statement in input rune array
	stmtPos := 0
	nextfunc := t.nextutf8
	// we assume that the statements have base line number 1
	t.lines = 1

	push := func(s string) {
		// if hasPrefixes(s, []string{"#", "--"}) {}
		// s = strings.TrimSpace(s)
		fmt.Printf("push \"%s\" into res, lines: %d\n", s, t.lines)
		if t.lines == 1 && hasPrefixes(s, []string{"#", "--"}) {
			return
		} else if t.lines > 1 {
			pieces := strings.Split(s, "\n")
			sb := strings.Builder{}
			for _, piece := range pieces {
				if !hasPrefixes(piece, []string{"#", "--"}) {
					if sb.Len() != 0 {
						sb.WriteString("\n")
					}
					sb.WriteString(piece)
				}
			}
			s = sb.String()
		}
		res = append(res, s)
		t.skipBlank()
		stmtPos = t.cursor
		t.lines = 1
	}

	for {
		r := nextfunc()
		switch r {
		case utf8.RuneError:
			if stmtPos < t.cursor {
				// res = append(res, t.getString(stmtPos, t.cursor-1))
				push(t.getString(stmtPos, t.cursor))
			}
			return res
		case '\'', '"', '`':
			t.scanString(r)
			t.skipBlank()
			tokenPos = t.cursor
		case ';':
			if delimiter == DEFAULT_DELIMITER {
				// res = append(res, t.getString(stmtPos, t.cursor))
				push(t.getString(stmtPos, t.cursor))
				tokenPos = stmtPos
			}
		case ' ', '\n':
			// why cursor - 1? the space unicode doesn't need to be parsed
			s := t.getString(tokenPos, t.cursor-1)
			tokenPos = t.cursor
			if r == '\n' {
				t.lines++
			}
			switch {
			// there are not any unparsed token in front of " ", for example
			// select * from a; select * from b;
			//                 ^
			case s == "":
				continue
			// when meeting the "delimiter" token, for example
			// DELIMITER //
			// select * from a;
			// select * from b;
			// //
			// delimiter ;
			case strings.EqualFold(s, KEYWORD_DELIMITER):
				delimiter = t.scanMysqlDelimiter()
				stmtPos = t.cursor
				tokenPos = t.cursor
			case s == delimiter:
				push(t.getString(stmtPos, t.cursor-1))
				tokenPos = t.cursor
			case strings.HasSuffix(s, delimiter):
				temp := t.getString(stmtPos, t.cursor-1)
				if delimiter != DEFAULT_DELIMITER {
					temp = strings.TrimSuffix(temp, delimiter)
				}
				push(temp)
				tokenPos = t.cursor
			}
		case '/':
			nr := nextfunc()
			if nr == '*' {
				t.scanMultComment()
				t.skipBlank()
				tokenPos = t.cursor
			}
		case '#':
			t.scanSingleComment()
			t.skipBlank()
			tokenPos = t.cursor
			// stmtPos = tokenPos
		case '-':
			nr := nextfunc()
			if nr == '-' {
				t.scanSingleComment()
				t.skipBlank()
				tokenPos = t.cursor
				// stmtPos = tokenPos
			}
		default:
		}
	}
}

func (t *Tokenizer) scanString(delimiter rune) {
	for {
		r := t.nextutf8()
		switch r {
		case utf8.RuneError:
			return
		case '\\':
			t.nextutf8()
		case delimiter:
			return
		}
	}
}

func (t *Tokenizer) scanIdentifier(delimiter rune) {
	for {
		r := t.nextutf8()
		switch r {
		case utf8.RuneError:
			return
		case delimiter:
			return
		}
	}
}

func (t *Tokenizer) scanSingleComment() {
	for {
		r := t.nextutf8()
		switch r {
		case utf8.RuneError:
			return
		case '\n':
			t.lines++
			return
		}
	}
}

func (t *Tokenizer) scanMultComment() {
	for {
		r := t.nextutf8()
		switch r {
		case utf8.RuneError:
			return
		case '*':
			nr := t.nextutf8()
			if nr == '/' {
				return
			}
		}
	}
}

func (t *Tokenizer) scanMysqlDelimiter() string {
	pos := t.cursor
	for {
		r := t.nextutf8()
		switch r {
		case utf8.RuneError:
			return t.getString(pos, t.cursor-1)
		case ' ', '\n':
			return t.getString(pos, t.cursor-1)
		}
	}
}

// getString 方法可以增加边界检查
func (t *Tokenizer) getString(startPos, endPos int) string {
	if startPos < 0 || endPos > len(t.text) {
		return ""
	}
	res := t.text[startPos:endPos]
	return res
}

// nextutf8 advances the cursor and returns the next rune
// Returns utf8.RuneError if no valid rune can be decoded
func (t *Tokenizer) nextutf8() rune {
	r, size := utf8.DecodeRuneInString(t.text[t.cursor:])
	t.cursor += size
	return r
}

func (t *Tokenizer) skipBlank() {
	for {
		r, size := utf8.DecodeRuneInString(t.text[t.cursor:])
		if !unicode.IsSpace(r) {
			return
		}
		t.cursor += size
		if r == '\n' {
			t.lines++
		}
	}
}

func hasPrefixes(s string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}
