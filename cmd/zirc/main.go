package main

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/zakshub/zirc-client/internal/app"
	"github.com/zakshub/zirc-client/internal/irc"
)

var version = "dev"

type listFlag []string
func (l *listFlag) String() string { return strings.Join(*l,",") }
func (l *listFlag) Set(v string) error { *l=append(*l,v);return nil }

func randomNick() (string,error) {
	b:=make([]byte,5)
	if _,err:=rand.Read(b);err!=nil{return "",err}
	return "zirc-"+hex.EncodeToString(b),nil
}

func main(){
	var channels listFlag
	server:=flag.String("server","","IRC server host:port")
	nick:=flag.String("nick","zirc-user","IRC nickname")
	random:=flag.Bool("random-nick",false,"Generate a random local nickname")
	username:=flag.String("username","","IRC username; defaults to nickname")
	tlsOn:=flag.Bool("tls",true,"Use TLS with certificate verification")
	proxy:=flag.String("proxy","","SOCKS5 proxy host:port; no direct fallback")
	proxyUser:=flag.String("proxy-user","","SOCKS5 username")
	passwordEnv:=flag.String("proxy-password-env","","Environment variable containing SOCKS5 password; never pass password as an argument")
	saslUser:=flag.String("sasl-user","","IRC SASL PLAIN account username")
	saslPasswordEnv:=flag.String("sasl-password-env","","Environment variable containing SASL password; never pass password as an argument")
	realname:=flag.String("realname","ZIRC User","IRC real name")
	preview:=flag.Bool("identity-preview",false,"Print the connection identity and route before connecting")
	showVersion:=flag.Bool("version",false,"Print ZIRC version")
	flag.Var(&channels,"join","Channel to join; repeatable")
	flag.Parse()
	if *showVersion { fmt.Println("ZIRC",version); return }

	if *random {
		generated,err:=randomNick()
		if err!=nil{fmt.Fprintln(os.Stderr,"Cannot generate nickname:",err);os.Exit(1)}
		*nick=generated
	}
	if *username=="" {*username=*nick}
	password:=""
	if *passwordEnv!="" {
		password=os.Getenv(*passwordEnv)
		if password==""{fmt.Fprintln(os.Stderr,"The configured proxy password environment variable is empty or unavailable.");os.Exit(2)}
	}
	saslPassword:=""
	if *saslPasswordEnv!="" { saslPassword=os.Getenv(*saslPasswordEnv);if saslPassword==""{fmt.Fprintln(os.Stderr,"The configured SASL password environment variable is empty or unavailable.");os.Exit(2)} }
	if (*saslUser=="")!=(saslPassword=="") { fmt.Fprintln(os.Stderr,"SASL requires both --sasl-user and --sasl-password-env.");os.Exit(2) }
	if saslPassword!="" && !*tlsOn { fmt.Fprintln(os.Stderr,"SASL PLAIN requires TLS to protect credentials in transit.");os.Exit(2) }
	if (*proxyUser=="")!=(password=="") {
		fmt.Fprintln(os.Stderr,"SOCKS5 authentication requires both --proxy-user and --proxy-password-env.")
		os.Exit(2)
	}
	if *proxy!="" && !strings.Contains(*proxy,":") {
		fmt.Fprintln(os.Stderr,"Invalid --proxy. Expected SOCKS5 host:port.")
		os.Exit(2)
	}
	if *proxy=="" && (*proxyUser!=""||password!="") {
		fmt.Fprintln(os.Stderr,"SOCKS5 credentials require --proxy.")
		os.Exit(2)
	}
	if *preview {
		route:="DIRECT: IRC server can see your public IP"
		if *proxy!=""{route="SOCKS5 proxy configured; proxy operator can see your source IP; Tor is not independently verified"}
		fmt.Printf("ZIRC identity preview\nServer: %s\nNickname: %s\nUsername: %s\nReal name: %s\nTLS: %t\nSASL account: %s\nRoute: %s\n",*server,*nick,*username,*realname,*tlsOn,*saslUser,route)
		if *server==""{fmt.Fprintln(os.Stderr,"Set --server to connect after reviewing this identity.") ;os.Exit(2)}
		return
	}
	if *server=="" {
		fmt.Fprintln(os.Stderr,"Usage: zirc --server irc.example.net:6697 --tls --proxy 127.0.0.1:9050 --random-nick")
		os.Exit(2)
	}
	cfg:=irc.Config{Address:*server,Nick:*nick,Username:*username,RealName:*realname,TLS:*tlsOn,Proxy:*proxy,ProxyUser:*proxyUser,ProxyPass:password,SASLUser:*saslUser,SASLPass:saslPassword,Channels:channels}
	if err:=app.Run(cfg);err!=nil{fmt.Fprintln(os.Stderr,"zirc:",err);os.Exit(1)}
}
