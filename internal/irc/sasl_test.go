package irc

import (
	"encoding/base64"
	"strings"
	"testing"
)

func TestSASLPlainNegotiation(t *testing.T) {
	c:=NewClient(Config{Nick:"tester",SASLUser:"account",SASLPass:"private"})
	c.handleCAP(Message{Params:[]string{"*","LS"},Trailing:"multi-prefix sasl=PLAIN"})
	if got:=<-c.write;got!="CAP REQ :sasl"{t.Fatalf("unexpected capability request: %s",got)}
	c.handleCAP(Message{Params:[]string{"*","ACK"},Trailing:"sasl"})
	if got:=<-c.write;got!="AUTHENTICATE PLAIN"{t.Fatalf("unexpected SASL mechanism: %s",got)}
	c.saslNegotiating=true
	c.handleAuthenticateChallenge()
	got:=<-c.write
	if !strings.HasPrefix(got,"AUTHENTICATE "){t.Fatalf("unexpected SASL payload command: %s",got)}
	decoded,err:=base64.StdEncoding.DecodeString(strings.TrimPrefix(got,"AUTHENTICATE "))
	if err!=nil{t.Fatal(err)}
	if string(decoded)!="\x00account\x00private"{t.Fatalf("unexpected SASL payload contents")}
	if !c.saslPayloadSent{t.Fatal("SASL payload was not marked sent")}
}

func (c *Client) handleAuthenticateChallenge() { c.sendSASLPlain() }
