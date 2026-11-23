package irc

import (
	"strings"
	"testing"
)

func FuzzParseLine(f *testing.F) {
	seeds := []string{
		"PING :server",
		":nick!user@host PRIVMSG #channel :hello world",
		":server 353 nick = #channel :@alice +bob",
		"",
		":",
		"CAP nick LS * :sasl=PLAIN",
		"NOTICE AUTH :*** Looking up your hostname",
	}
	for _, seed := range seeds { f.Add(seed) }
	f.Fuzz(func(t *testing.T, line string) {
		if len(line) > 64*1024 { t.Skip() }
		m := ParseLine(line)
		if strings.ContainsAny(m.Command, " \r\n") { t.Fatalf("invalid command token %q", m.Command) }
		for _, p := range m.Params {
			if strings.ContainsAny(p, " \r\n") { t.Fatalf("invalid parameter token %q", p) }
		}
	})
}

func FuzzFormatValidate(f *testing.F) {
	f.Add("PRIVMSG", "#channel", "hello world")
	f.Add("NICK", "zirc-user", "")
	f.Add("NOTICE", "nick", "")
	f.Add("PRIVMSG", "#x", "line\r\ninjection")
	f.Fuzz(func(t *testing.T, command, target, body string) {
		if len(command)+len(target)+len(body) > 4096 { t.Skip() }
		line := Format(command, target, body)
		err := Validate(line)
		if err == nil && (strings.ContainsAny(line, "\r\n") || len(line) > 510) {
			t.Fatalf("Validate accepted unsafe line %q", line)
		}
	})
}
