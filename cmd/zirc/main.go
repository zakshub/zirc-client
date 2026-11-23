package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/zakshub/zirc-client/internal/app"
	"github.com/zakshub/zirc-client/internal/irc"
)

type listFlag []string
func (l *listFlag) String() string { return strings.Join(*l,",") }
func (l *listFlag) Set(v string) error { *l=append(*l,v);return nil }

func main(){
	var channels listFlag
	server:=flag.String("server","","IRC server host:port")
	nick:=flag.String("nick","zirc-user","IRC nickname")
	tlsOn:=flag.Bool("tls",true,"Use TLS with certificate verification")
	proxy:=flag.String("proxy","","SOCKS5 proxy host:port, e.g. 127.0.0.1:9050 for local Tor")
	realname:=flag.String("realname","ZIRC User","Non-identifying IRC real name")
	flag.Var(&channels,"join","Channel to join; repeatable")
	flag.Parse()
	if *server==""{
		fmt.Fprintln(os.Stderr,"Usage: zirc --server irc.example.net:6697 --tls --proxy 127.0.0.1:9050 --nick nickname")
		os.Exit(2)
	}
	if *proxy!="" && !strings.Contains(*proxy,":"){
		fmt.Fprintln(os.Stderr,"Invalid --proxy. Expected SOCKS5 host:port.")
		os.Exit(2)
	}
	cfg:=irc.Config{Address:*server,Nick:*nick,RealName:*realname,TLS:*tlsOn,Proxy:*proxy,Channels:channels}
	if err:=app.Run(cfg);err!=nil{fmt.Fprintln(os.Stderr,"zirc:",err);os.Exit(1)}
}
