package app

import (
	"testing"

	"github.com/zakshub/zirc-client/internal/irc"
)

func TestChannelBuffersStaySeparate(t *testing.T) {
	m:=&model{active:"status",channels:[]string{"status"},lines:map[string][]string{"status":{}},users:map[string][]string{},topics:map[string]string{}}
	m.handle(irc.Event{Kind:"message",Target:"#one",Nick:"alice",Text:"first"})
	m.handle(irc.Event{Kind:"message",Target:"#two",Nick:"bob",Text:"second"})
	if got:=m.lines["#one"];len(got)!=1||got[0]!="<alice> first"{t.Fatalf("unexpected first channel buffer: %#v",got)}
	if got:=m.lines["#two"];len(got)!=1||got[0]!="<bob> second"{t.Fatalf("unexpected second channel buffer: %#v",got)}
}

func TestChannelMembersAndTopic(t *testing.T) {
	m:=&model{active:"#one",channels:[]string{"status"},lines:map[string][]string{"status":{}},users:map[string][]string{},topics:map[string]string{}}
	m.handle(irc.Event{Kind:"join",Target:"#one",Nick:"alice"})
	m.handle(irc.Event{Kind:"names",Target:"#one",Text:"@alice +bob charlie"})
	m.handle(irc.Event{Kind:"topic",Target:"#one",Text:"project discussion"})
	if len(m.users["#one"])!=3{t.Fatalf("expected three unique users, got %#v",m.users["#one"])}
	if m.topics["#one"]!="project discussion"{t.Fatalf("unexpected topic: %q",m.topics["#one"])}
	m.handle(irc.Event{Kind:"part",Target:"#one",Nick:"bob"})
	if contains(m.users["#one"],"bob"){t.Fatal("parted user remains in member list")}
}

func TestSearchLinesIsCaseInsensitiveAndBounded(t *testing.T) {
	lines := []string{"hello world", "another line", "HELLO again", "hello third"}
	got := searchLines(lines, "Hello", 2)
	if len(got) != 2 { t.Fatalf("expected two results, got %#v", got) }
	if got[0] != "[1] hello world" || got[1] != "[3] HELLO again" {
		t.Fatalf("unexpected search results: %#v", got)
	}
	if got := searchLines(lines, "missing", 10); len(got) != 0 {
		t.Fatalf("expected no results, got %#v", got)
	}
}
