package app

import "testing"

func TestExpandAlias(t *testing.T) {
	tests := []struct { name, template, want string; args []string }{
		{"single argument", "/msg $1 hello", "/msg nick hello", []string{"nick", "there"}},
		{"remaining arguments", "/msg $1 $2-", "/msg nick hello there", []string{"nick", "hello", "there"}},
		{"all arguments", "/say $*", "/say one two", []string{"one", "two"}},
		{"escaped dollar", "/say $$1", "/say $1", nil},
		{"missing argument", "/msg $1 hello", "/msg  hello", nil},
	}
	for _, tt := range tests { t.Run(tt.name, func(t *testing.T) { got,err:=expandAlias(tt.template,tt.args); if err!=nil { t.Fatal(err) }; if got!=tt.want { t.Fatalf("got %q, want %q",got,tt.want) } }) }
}

func TestExpandAliasRequiresCommand(t *testing.T) {
	if _,err:=expandAlias("echo $*",nil); err==nil { t.Fatal("expected non-command template to fail") }
}

func TestValidAliasName(t *testing.T) {
	for _, name := range []string{"hello","greet2","my_alias"} { if !validAliasName(name) { t.Errorf("expected %q to be valid",name) } }
	for _, name := range []string{"","2fast","bad-name","has space"} { if validAliasName(name) { t.Errorf("expected %q to be invalid",name) } }
}
