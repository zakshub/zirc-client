package irc

import (
	"fmt"
	"strings"
)

type Message struct {
	Prefix string
	Command string
	Params []string
	Trailing string
}

func ParseLine(line string) Message {
	line = strings.TrimRight(line, "\r\n")
	m := Message{}
	if strings.HasPrefix(line, ":") {
		if i := strings.IndexByte(line, ' '); i >= 0 {
			m.Prefix, line = line[1:i], strings.TrimSpace(line[i+1:])
		}
	}
	if i := strings.Index(line, " :"); i >= 0 {
		m.Trailing, line = line[i+2:], line[:i]
	}
	fields := strings.Fields(line)
	if len(fields) > 0 { m.Command = strings.ToUpper(fields[0]); m.Params = fields[1:] }
	return m
}

func Format(command string, params ...string) string {
	var b strings.Builder
	b.WriteString(strings.ToUpper(command))
	for i, p := range params {
		b.WriteByte(' ')
		if i == len(params)-1 && (strings.ContainsAny(p, " :") || p == "") { b.WriteByte(':') }
		b.WriteString(p)
	}
	return b.String()
}

func Validate(line string) error {
	if strings.ContainsAny(line, "\r\n") { return fmt.Errorf("IRC command contains newline") }
	if len(line) > 510 { return fmt.Errorf("IRC command exceeds 510 bytes") }
	return nil
}
