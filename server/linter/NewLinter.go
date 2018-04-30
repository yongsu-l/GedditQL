package linter

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// NewLinter ...
func NewLinter(filename string) func(string, string) bool {
	reMap := map[string]string{}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("ERROR: File does not exist")
	}

	defer file.Close()
	b, _ := ioutil.ReadAll(file)
	rules := strings.Split(string(b), "\n")

	for _, rule := range rules {
		if rule == "" {
			continue
		}
		name, body := parseRule(rule)

		reMap[name] = parseBody(body, reMap)
	}

	return func(s string, entry string) bool {
		re := reMap[entry]
		linter := regexp.MustCompile("^" + re + "$")
		if !linter.MatchString(strings.ToLower(s)) {
			return false
		}
		return true
	}
}

func parseRule(s string) (string, string) {
	split := strings.Split(s, ":=")
	name := strings.TrimSpace(split[0])
	return name[1 : len(name)-1], strings.TrimSpace(split[1])
}

func parseBody(s string, rmap map[string]string) string {
	is := bytes.NewBufferString(s)

	var r bytes.Buffer
	for {
		c, err := is.ReadByte()
		if err != nil {
			break
		}
		switch c {
		case '\\':
			d, _ := is.ReadByte()
			r.WriteByte(d)
			break
		case '<':
			name := ""
			for {
				d, _ := is.ReadByte()
				if d == '>' {
					break
				}
				name += string(d)
			}
			r.WriteString("(" + rmap[name] + ")")
		default:
			r.WriteByte(c)
			break
		}
	}

	return r.String()
}
