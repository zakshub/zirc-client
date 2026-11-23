package irc

import "testing"

func TestParsePrivmsg(t *testing.T) {
	m := ParseLine(":nick!u@host PRIVMSG #linux :hello world")
	if m.Command != "PRIVMSG" || m.Prefix != "nick!u@host" || m.Trailing != "hello world" || len(m.Params) != 1 || m.Params[0] != "#linux" {
		t.Fatalf("unexpected parse: %#v", m)
	}
}
func TestPing(t *testing.T) {
	m := ParseLine("PING :server")
	if m.Command != "PING" || m.Trailing != "server" { t.Fatalf("unexpected: %#v", m) }
}
func TestFormat(t *testing.T) {
	if got := Format("PRIVMSG", "#x", "hello world"); got != "PRIVMSG #x :hello world" { t.Fatalf("got %q", got) }
}
func TestRejectNewlines(t *testing.T) {
	if Validate("QUIT\r\nPRIVMSG #x :oops") == nil { t.Fatal("expected rejection") }
}
