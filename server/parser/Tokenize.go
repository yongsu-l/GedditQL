package parser

import (
	"bytes"
	"errors"
)

// Tokenize returns an array of tokens
///////////////////////////////////////////////////////////////////////////////////////////////
func Tokenize(s string) ([]string, error) {
	var buffer bytes.Buffer

	c, err := countTokens(s)
	if err != nil {
		return nil, err
	}

	r := make([]string, c)
	j := 0
	PopBuffer := func() {
		if buffer.String() != "" {
			r[j] = buffer.String()
			buffer.Reset()
			j++
		}
	}

	for i := 0; i < len(s); i++ {
		switch nb := s[i]; nb {
		case ' ', '\n':
			PopBuffer()
			break
		case ';', ',', '.', '(', ')', '+', '-', '*', '/', '%', '=':
			PopBuffer()
			r[j] = string(nb)
			j++
			break
		case '<':
			PopBuffer()
			buffer.WriteByte(nb)
			if i+1 < len(s) && (s[i+1] == '>' || s[i+1] == '=') {
				i++
				buffer.WriteByte(s[i])
			}
			PopBuffer()
			break
		case '>', '!':
			PopBuffer()
			buffer.WriteByte(nb)
			if i+1 < len(s) && s[i+1] == '=' {
				i++
				buffer.WriteByte(s[i])
			}
			PopBuffer()
			break
		case '\'':
			PopBuffer()
			buffer.WriteByte('"')
			for {
				i++
				if s[i] == ';' {
					return nil, errors.New("Missing Closing Quotes")
				}
				if s[i] == '\'' {
					buffer.WriteByte('"')
					break
				}
				buffer.WriteByte(s[i])
			}
			PopBuffer()
			break
		case '"':
			PopBuffer()
			buffer.WriteByte(nb)
			for {
				i++
				if s[i] == ';' {
					return nil, errors.New("Missing Closing Quotes")
				}
				buffer.WriteByte(s[i])
				if s[i] == '"' {
					break
				}
			}
			PopBuffer()
			break
		default:
			buffer.WriteByte(nb)
		}
	}

	if j < c {
		PopBuffer()
	}

	return r, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
func countTokens(s string) (int, error) {
	c := 0
	b := false
	CheckB := func() {
		if b {
			c++
			b = false
		}
	}

	for i := 0; i < len(s); i++ {
		switch nb := s[i]; nb {
		case ' ', '\n':
			CheckB()
			break
		case ';', ',', '.', '(', ')', '+', '-', '*', '/', '%', '=':
			CheckB()
			c++
			break
		case '<':
			if i+1 < len(s) && (s[i+1] == '>' || s[i+1] == '=') {
				i++
			}
			CheckB()
			c++
			break
		case '>', '!':
			if i+1 < len(s) && s[i+1] == '=' {
				i++
			}
			CheckB()
			c++
			break
		case '\'':
			for {
				i++
				if s[i] == ';' {
					return 0, errors.New("Missing Closing Quotes")
				}
				if s[i] == '\'' {
					break
				}
			}
			c++
			break
		case '"':
			for {
				i++
				if s[i] == ';' {
					return 0, errors.New("Missing Closing Quotes")
				}
				if s[i] == '"' {
					break
				}
			}
			c++
			break
		default:
			b = true
		}
	}

	return c, nil
}
