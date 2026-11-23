package app

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/zakshub/zirc-client/internal/irc"
)

type model struct {
	client *irc.Client
	active string
	channels []string
	lines []string
	input string
	status string
	offset int
}
func Run(cfg irc.Config) error {
	if !isTerminal(){return fmt.Errorf("ZIRC requires an interactive terminal")}
	state,err:=makeRaw(int(os.Stdin.Fd()));if err!=nil{return fmt.Errorf("cannot enable terminal mode: %w",err)}
	defer state.restore(int(os.Stdin.Fd()))
	client:=irc.NewClient(cfg)
	if err:=client.Connect();err!=nil{return err}
	defer client.Close()
	m:=&model{client:client,active:"status",channels:[]string{"status"},status:"Connecting"}
	if cfg.TLS {m.status="TLS enabled"}
	if cfg.Proxy!="" {m.status+=" | SOCKS5 proxy configured; Tor not independently verified"} else {m.status+=" | DIRECT ROUTE: IRC server sees your public IP"}
	fmt.Print("\033[?1049h\033[?25l")
	defer fmt.Print("\033[?25h\033[?1049l")
	keys:=make(chan byte,128)
	go func(){buf:=make([]byte,1);for{n,e:=os.Stdin.Read(buf);if e!=nil{close(keys);return};if n>0{keys<-buf[0]}}}()
	sig:=make(chan os.Signal,1);signal.Notify(sig,os.Interrupt,syscall.SIGTERM);defer signal.Stop(sig)
	tick:=time.NewTicker(100*time.Millisecond);defer tick.Stop()
	for {
		select {
		case ev:=<-client.Events():m.handle(ev)
		case b,ok:=<-keys:if !ok{return nil};if m.key(b){return nil}
		case <-sig:return nil
		case <-client.Done():m.status="Disconnected";m.draw();return nil
		case <-tick.C:m.draw()
		}
	}
}
func isTerminal()bool{return terminalAvailable(int(os.Stdin.Fd()))}
func (m *model)add(s string){m.lines=append(m.lines,s);if len(m.lines)>3000{m.lines=m.lines[len(m.lines)-3000:]}}
func (m *model)handle(e irc.Event){
	switch e.Kind {
	case "status":m.status=e.Text;m.add("* "+e.Text)
	case "error":m.status=e.Text;m.add("! "+e.Text)
	case "message":
		if !contains(m.channels,e.Target){m.channels=append(m.channels,e.Target)}
		m.add("<"+e.Nick+"> "+e.Text)
	case "notice":m.add(e.Text)
	}
}
func contains(xs []string,s string)bool{for _,x:=range xs{if x==s{return true}};return false}
func (m *model)key(b byte)bool{
	switch b{
	case 3:return true
	case 13:
		s:=strings.TrimSpace(m.input);m.input=""
		if s!=""{m.command(s)}
	case 8,127:if len(m.input)>0{m.input=m.input[:len(m.input)-1]}
	case 9:if len(m.channels)>0{for i,ch:=range m.channels{if ch==m.active{m.active=m.channels[(i+1)%len(m.channels)];break}}}
	case 27:
	default:if b>=32&&b<127{m.input+=string(b)}
	}
	return false
}
func (m *model)command(s string){
	if strings.HasPrefix(s,"/"){
		p:=strings.Fields(s);if len(p)==0{return}
		switch strings.ToLower(p[0]){
		case "/quit":m.client.Close()
		case "/join":if len(p)>1{m.client.Send("JOIN",p[1]);if !contains(m.channels,p[1]){m.channels=append(m.channels,p[1])};m.active=p[1]}
		case "/nick":if len(p)>1{m.client.Send("NICK",p[1])}
		case "/msg":if len(p)>2{target:=p[1];body:=strings.TrimSpace(strings.TrimPrefix(s,p[0]+" "+target));m.client.Send("PRIVMSG",target,body);m.add("-> "+target+": "+body)}
		case "/part":if m.active!="status"{m.client.Send("PART",m.active);m.add("* Leaving "+m.active);m.active="status"}
		case "/help":m.add("Commands: /join #channel | /nick name | /msg nick message | /part | /quit")
		default:m.add("Unknown command. Type /help")
		}
		return
	}
	if m.active=="status"{m.add("Join a channel first with /join #channel");return}
	m.client.Send("PRIVMSG",m.active,s);m.add("<you> "+s)
}
func (m *model)draw(){
	w,h:=terminalSize();if w<50{w=50};if h<12{h=12}
	var b strings.Builder
	b.WriteString("\033[H\033[2J")
	b.WriteString("\033[48;5;235m\033[38;5;81m ZIRC CLIENT \033[0m")
	b.WriteString("\033[38;5;245m  "+m.status+"\033[0m\n")
	b.WriteString("\033[38;5;240m"+strings.Repeat("─",w)+"\033[0m\n")
	side:=20;bodyW:=w-side-3;if bodyW<12{side=0;bodyW=w-2}
	if side>0{b.WriteString("\033[38;5;81m NETWORK / CHANNELS\033[0m");b.WriteString(strings.Repeat(" ",max(1,side-18)));b.WriteString("\033[38;5;81m│ "+m.active+"\033[0m\n")}else{b.WriteString("\033[38;5;81m "+m.active+"\033[0m\n")}
	rows:=h-6
	start:=len(m.lines)-rows-m.offset;if start<0{start=0}
	end:=len(m.lines)-m.offset;if end>len(m.lines){end=len(m.lines)};if end<start{end=start}
	for i:=0;i<rows;i++{
		if side>0{
			if i<len(m.channels){ch:=m.channels[i];if ch==m.active{b.WriteString("\033[48;5;238m\033[38;5;220m "+pad(ch,side-2)+"\033[0m")}else{b.WriteString("\033[38;5;250m "+pad(ch,side-2)+"\033[0m")}}else{b.WriteString(strings.Repeat(" ",side))}
			b.WriteString("\033[38;5;240m│\033[0m")
		}
		idx:=start+i;if idx<end{line:=m.lines[idx];if len(line)>bodyW{line=line[:bodyW]};b.WriteString(line)}
		b.WriteString("\033[K\n")
	}
	b.WriteString("\033[38;5;240m"+strings.Repeat("─",w)+"\033[0m\n")
	b.WriteString("\033[38;5;220m> \033[0m"+m.input+"\033[K\n")
	b.WriteString("\033[38;5;245m Tab: switch channel  /help: commands  Ctrl+C: quit\033[0m")
	fmt.Print(b.String())
}
func pad(s string,n int)string{if len(s)>n{return s[:n]};return s+strings.Repeat(" ",n-len(s))}
func max(a,b int)int{if a>b{return a};return b}
