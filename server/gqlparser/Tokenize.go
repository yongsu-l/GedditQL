package gqlparser

import (
	"bytes"
)

// Tokenize returns an array of tokens
func Tokenize(s string) []string {
	var buffer bytes.Buffer

	c := CountTokens(s)
	r := make([]string, c)
	j := 0
	ReadBuffer := func() {
		if buffer.String() != "" {
			r[j] = buffer.String()
			buffer.Reset()
			j++
		}
	}

	for i := 0; i < len(s); i++ {
		switch nb := s[i]; nb {
		case ' ', '\n':
			ReadBuffer()
			break
		case ';', ',', '.', '(', ')', '+', '-', '*', '/', '%', '=':
			ReadBuffer()
			r[j] = string(nb)
			j++
			break
		case '<':
			ReadBuffer()
			buffer.WriteByte(nb)
			if i+1 < len(s) && (s[i+1] == '>' || s[i+1] == '=') {
				i++
				buffer.WriteByte(s[i])
			}
			ReadBuffer()
			break
		case '>', '!':
			buffer.WriteByte(nb)
			if i+1 < len(s) && s[i+1] == '=' {
				i++
				buffer.WriteByte(s[i])
			}
			ReadBuffer()
			break
		default:
			buffer.WriteByte(nb)
		}
	}

	if j < c {
		ReadBuffer()
	}

	return r
}

// CountTokens returns the number of tokens in a string
func CountTokens(s string) int {
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
		default:
			b = true
		}
	}

	return c
}
