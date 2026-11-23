package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/zakshub/zirc-client/internal/irc"
)

func Run(cfg irc.Config) error {
	client:=irc.NewClient(cfg)
	if err:=client.Connect(); err!=nil {return err}
	defer client.Close()

	fmt.Print("\033[2J\033[H")
	fmt.Println("\033[1;36mZIRC CLIENT\033[0m  |  Privacy-first terminal IRC")
	fmt.Printf("Server: %s | Transport: ",cfg.Address)
	if cfg.TLS {fmt.Println("TLS encrypted")} else {fmt.Println("\033[1;31mPLAINTEXT: your traffic is not encrypted\033[0m")}
	if cfg.Proxy!="" { fmt.Println("Routing: SOCKS5 proxy configured; origin IP is presented to the proxy, not directly to IRC server") } else { fmt.Println("\\033[1;31mRouting: DIRECT; IRC server can see your public IP unless it provides a cloak\\033[0m") }
	fmt.Println("Commands: /join #channel  /nick name  /msg nick message  /quit")
	fmt.Println(strings.Repeat("─",72))

	input:=make(chan string,8)
	go func(){s:=bufio.NewScanner(os.Stdin); for s.Scan(){input<-s.Text()}; close(input)}()
	tick:=time.NewTicker(100*time.Millisecond)
	defer tick.Stop()
	active:="status"
	for {
		select {
		case ev:=<-client.Events():
			switch ev.Kind {
			case "status": fmt.Printf("\n\033[1;32m* %s\033[0m\n",ev.Text)
			case "error": fmt.Printf("\n\033[1;31m! %s\033[0m\n",ev.Text)
			case "message": fmt.Printf("\n\033[36m<%s>\033[0m %s\n",ev.Nick,ev.Text)
			case "notice": fmt.Printf("\n\033[90m%s\033[0m\n",ev.Text)
			}
			fmt.Printf("\033[1;33m[%s] > \033[0m",active)
		case line,ok:=<-input:
			if !ok {return nil}
			line=strings.TrimSpace(line)
			if line=="" {fmt.Printf("\033[1;33m[%s] > \033[0m",active);continue}
			if strings.HasPrefix(line,"/") {
				parts:=strings.Fields(line)
				switch strings.ToLower(parts[0]) {
				case "/quit": return nil
				case "/join":
					if len(parts)>1 {client.Send("JOIN",parts[1]);active=parts[1]}
				case "/nick":
					if len(parts)>1 {client.Send("NICK",parts[1])}
				case "/msg":
					if len(parts)>2 {target:=parts[1]; body:=strings.TrimSpace(strings.TrimPrefix(line,parts[0]+" "+target));client.Send("PRIVMSG",target,body);fmt.Printf("-> %s: %s\n",target,body)}
				case "/help":
					fmt.Println("Commands: /join #channel, /nick name, /msg nick message, /quit")
				case "/part":
					client.Send("PART",active);active="status"
				default:fmt.Println("Unknown command. Use /help")
				}
			} else if active!="status" {
				client.Send("PRIVMSG",active,line)
				fmt.Printf("\033[32m<you>\033[0m %s\n",line)
			} else {fmt.Println("Join a channel first with /join #channel")}
			fmt.Printf("\033[1;33m[%s] > \033[0m",active)
		case <-client.Done():
			return fmt.Errorf("IRC connection closed")
		case <-tick.C:
		}
	}
}
