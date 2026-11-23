package app

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"unicode"

	"github.com/zakshub/zirc-client/internal/irc"
)

type model struct {
	client *irc.Client
	identity irc.Config
	active string
	channels []string
	lines map[string][]string
	users map[string][]string
	topics map[string]string
	input string
	status string
	offset int
	quitting bool
}

func Run(cfg irc.Config) error {
	if !isTerminal() { return fmt.Errorf("ZIRC requires an interactive terminal") }
	state,err:=makeRaw(int(os.Stdin.Fd()))
	if err!=nil { return fmt.Errorf("cannot enable terminal mode: %w",err) }
	defer state.restore(int(os.Stdin.Fd()))
	client:=irc.NewClient(cfg)
	if err:=client.Connect();err!=nil { return err }
	defer func(){ if client!=nil { client.Close() } }()
	m:=&model{client:client,identity:cfg,active:"status",channels:[]string{"status"},lines:map[string][]string{"status":{}},users:map[string][]string{},topics:map[string]string{},status:"Connecting"}
	if cfg.TLS { m.status="TLS enabled" }
	if cfg.Proxy!="" { m.status+=" | SOCKS5 proxy configured; Tor not independently verified" } else { m.status+=" | DIRECT ROUTE: IRC server sees your public IP" }
	fmt.Print("\033[?1049h\033[?25l")
	defer fmt.Print("\033[?25h\033[?1049l")
	keys:=make(chan rune,128)
	go func(){reader:=bufio.NewReader(os.Stdin);for{r,_,e:=reader.ReadRune();if e!=nil{close(keys);return};keys<-r}}()
	sig:=make(chan os.Signal,1);signal.Notify(sig,os.Interrupt,syscall.SIGTERM);defer signal.Stop(sig)
	tick:=time.NewTicker(100*time.Millisecond);defer tick.Stop()
	clientEvents,clientDone:=client.Events(),client.Done()
	var reconnectTimer *time.Timer
	var reconnectC <-chan time.Time
	attempt:=0
	scheduleReconnect:=func() bool {
		attempt++
		if attempt>5 { m.status="Reconnect attempts exhausted. Use /quit and start ZIRC again."; m.add("status","! Automatic reconnect stopped after five attempts"); return false }
		delay:=reconnectDelay(attempt)
		m.status=fmt.Sprintf("Disconnected. Retrying through the same configured route in %s (%d/5)",delay,attempt)
		m.add("status","* "+m.status)
		reconnectTimer=time.NewTimer(delay)
		reconnectC=reconnectTimer.C
		return true
	}
	for {
		select {
		case ev:=<-clientEvents:m.handle(ev)
		case r,ok:=<-keys:if !ok{return nil};if m.key(r){return nil}
		case <-sig:return nil
		case <-clientDone:
			clientEvents=nil;clientDone=nil
			m.add("status","* Connection lost")
			if m.quitting { return nil }
			if !scheduleReconnect(){m.draw();return nil}
		case <-reconnectC:
			reconnectC=nil
			if reconnectTimer!=nil { reconnectTimer=nil }
			cfg.Nick=m.identity.Nick
			cfg.Channels=channelTargets(m.channels)
			next:=irc.NewClient(cfg)
			if err:=next.Connect();err!=nil {
				m.add("status","! Reconnect failed: "+err.Error())
				if !scheduleReconnect(){m.draw();return nil}
			} else {
				client=next;m.client=next
				clientEvents,clientDone=next.Events(),next.Done()
				attempt=0
				m.status="Reconnected through the configured route"
				m.add("status","* "+m.status)
			}
		case <-tick.C:if m.quitting{return nil};m.draw()
		}
	}
}

func isTerminal()bool{return terminalAvailable(int(os.Stdin.Fd()))}
func (m *model)add(target,line string){
	if target==""{target="status"}
	if m.lines[target]==nil{m.lines[target]=[]string{}}
	m.lines[target]=append(m.lines[target],line)
	if len(m.lines[target])>3000{m.lines[target]=m.lines[target][len(m.lines[target])-3000:]}
}
func (m *model)ensureChannel(ch string){
	if ch==""{return}
	if !contains(m.channels,ch){m.channels=append(m.channels,ch)}
	if m.lines[ch]==nil{m.lines[ch]=[]string{}}
	if m.users[ch]==nil{m.users[ch]=[]string{}}
}
func (m *model)handle(e irc.Event){
	switch e.Kind {
	case "status":m.status=e.Text;m.add("status","* "+e.Text)
	case "error":m.status=e.Text;m.add(m.active,"! "+e.Text)
	case "message":
		target:=e.Target
		if !strings.HasPrefix(target,"#"){target=e.Nick}
		m.ensureChannel(target)
		m.add(target,"<"+e.Nick+"> "+e.Text)
	case "join":
		m.ensureChannel(e.Target)
		m.users[e.Target]=appendUnique(m.users[e.Target],e.Nick)
		m.add(e.Target,"* "+e.Nick+" joined")
	case "part":
		m.removeUser(e.Target,e.Nick)
		m.add(e.Target,"* "+e.Nick+" left")
	case "quit":
		for _,ch:=range m.channels{m.removeUser(ch,e.Nick);if ch!="status"{m.add(ch,"* "+e.Nick+" quit")}}
	case "names":
		m.ensureChannel(e.Target)
		for _,name:=range strings.Fields(e.Text){name=strings.TrimLeft(name,"~&@%+");if name!=""{m.users[e.Target]=appendUnique(m.users[e.Target],name)}}
	case "topic":
		m.ensureChannel(e.Target);m.topics[e.Target]=e.Text;m.add(e.Target,"* Topic: "+e.Text)
	case "notice":m.add(m.active,e.Text)
	}
}
func appendUnique(xs []string,s string)[]string{if s==""||contains(xs,s){return xs};return append(xs,s)}
func contains(xs []string,s string)bool{for _,x:=range xs{if x==s{return true}};return false}
func (m *model)removeUser(ch,nick string){xs:=m.users[ch];out:=xs[:0];for _,x:=range xs{if x!=nick{out=append(out,x)}};m.users[ch]=out}
func (m *model)key(r rune)bool{
	switch r{
	case 3:return true
	case 13:
		s:=strings.TrimSpace(m.input);m.input=""
		if s!=""{m.command(s)}
	case 8,127:if len(m.input)>0{runes:=[]rune(m.input);m.input=string(runes[:len(runes)-1])}
	case 9:if len(m.channels)>0{for i,ch:=range m.channels{if ch==m.active{m.active=m.channels[(i+1)%len(m.channels)];m.offset=0;break}}}
	case 27:
	default:if r>=32&&r!=127&&r!=27{m.input+=string(r)}
	}
	return false
}
func (m *model)command(s string){
	if strings.HasPrefix(s,"/"){
		p:=strings.Fields(s);if len(p)==0{return}
		switch strings.ToLower(p[0]){
		case "/quit":m.quitting=true;m.client.Close()
		case "/join":if len(p)>1{m.client.Send("JOIN",p[1]);m.ensureChannel(p[1]);m.active=p[1]}
		case "/nick":if len(p)>1{m.client.Send("NICK",p[1]);m.identity.Nick=p[1]}
		case "/randomnick":
			b:=make([]byte,5);if _,err:=rand.Read(b);err!=nil{m.add(m.active,"! Cannot generate random nickname");return}
			nick:="zirc-"+hex.EncodeToString(b);m.client.Send("NICK",nick);m.identity.Nick=nick;m.add("status","* Requested random nickname "+nick)
		case "/msg","/query":
			if len(p)>1{target:=p[1];m.ensureChannel(target);m.active=target
				if len(p)>2{body:=strings.TrimSpace(strings.TrimPrefix(s,p[0]+" "+target));m.client.Send("PRIVMSG",target,body);m.add(target,"<you> "+body)}
			}
		case "/part":if m.active!="status"{ch:=m.active;m.client.Send("PART",ch);m.add(ch,"* Leaving "+ch);m.removeChannel(ch);m.active="status"}
		case "/names":if m.active!="status"{m.client.Send("NAMES",m.active)}
		case "/topic":if m.active!="status"{if len(p)>1{body:=strings.TrimSpace(strings.TrimPrefix(s,p[0]+" "));m.client.Send("TOPIC",m.active,body)}else{m.client.Send("TOPIC",m.active)}}
		case "/invite":if len(p)>1&&m.active!="status"{m.client.Send("INVITE",p[1],m.active)}
		case "/mode":if m.active!="status"{args:=append([]string{m.active},p[1:]...);m.client.Send("MODE",args...)}
		case "/search":
			if len(p)>1 {
				query:=strings.TrimSpace(strings.TrimPrefix(s,p[0]))
				matches:=searchLines(m.lines[m.active],query,20)
				m.add("status",fmt.Sprintf("Search in %s for %q: %d match(es)",m.active,query,len(matches)))
				for _,match:=range matches { m.add("status",match) }
			} else { m.add("status","Usage: /search words") }
		case "/privacy":
			route:="DIRECT. IRC server sees your public IP."
			if m.identity.Proxy!=""{route="SOCKS5 proxy configured. Tor is not independently verified."}
			m.add("status","Identity: nick="+m.identity.Nick+" username="+m.identity.Username+" realname="+m.identity.RealName)
			m.add("status","Transport: TLS="+fmt.Sprint(m.identity.TLS)+" | "+route)
		case "/help":m.add("status","Commands: /join #channel | /nick name | /randomnick | /msg nick message | /query nick | /part | /names | /topic [text] | /invite nick | /mode [flags] | /search words | /privacy | /quit")
		default:m.add(m.active,"Unknown command. Type /help")
		}
		return
	}
	if m.active=="status"{m.add("status","Join a channel first with /join #channel");return}
	m.client.Send("PRIVMSG",m.active,s);m.add(m.active,"<you> "+s)
}
func (m *model)removeChannel(ch string){
	out:=m.channels[:0];for _,x:=range m.channels{if x!=ch{out=append(out,x)}};m.channels=out
	delete(m.lines,ch);delete(m.users,ch);delete(m.topics,ch)
}
func (m *model)draw(){
	w,h:=terminalSize();if w<50{w=50};if h<12{h=12}
	var b strings.Builder
	b.WriteString("\033[H\033[2J")
	b.WriteString("\033[48;5;235m\033[38;5;81m ZIRC CLIENT \033[0m")
	b.WriteString("\033[38;5;245m  "+m.status+"\033[0m\n")
	if topic:=m.topics[m.active];topic!=""{b.WriteString("\033[38;5;220m Topic: "+topic+"\033[0m\n")}else{b.WriteString("\033[38;5;240m"+strings.Repeat("─",w)+"\033[0m\n")}
	side:=24;bodyW:=w-side-3;if bodyW<12{side=0;bodyW=w-2}
	if side>0{b.WriteString("\033[38;5;81m NETWORK / CHANNELS\033[0m");b.WriteString(strings.Repeat(" ",max(1,side-18)));b.WriteString("\033[38;5;81m│ "+m.active+"\033[0m\n")}else{b.WriteString("\033[38;5;81m "+m.active+"\033[0m\n")}
	rows:=h-6
	lines:=m.lines[m.active]
	start:=len(lines)-rows-m.offset;if start<0{start=0}
	end:=len(lines)-m.offset;if end>len(lines){end=len(lines)};if end<start{end=start}
	for i:=0;i<rows;i++{
		if side>0{
			switch {
			case i<len(m.channels):
				ch:=m.channels[i];if ch==m.active{b.WriteString("\033[48;5;238m\033[38;5;220m "+pad(ch,side-2)+"\033[0m")}else{b.WriteString("\033[38;5;250m "+pad(ch,side-2)+"\033[0m")}
			case i==len(m.channels):b.WriteString("\033[38;5;81m "+pad("USERS",side-2)+"\033[0m")
			case i>len(m.channels)&&i-len(m.channels)-1<len(m.users[m.active]):b.WriteString("\033[38;5;250m "+pad(m.users[m.active][i-len(m.channels)-1],side-2)+"\033[0m")
			default:b.WriteString(strings.Repeat(" ",side))
			}
			b.WriteString("\033[38;5;240m│\033[0m")
		}
		idx:=start+i;if idx<end{line:=lines[idx];if displayWidth(line)>bodyW{line=truncateCells(line,bodyW)};b.WriteString(line)}
		b.WriteString("\033[K\n")
	}
	b.WriteString("\033[38;5;240m"+strings.Repeat("─",w)+"\033[0m\n")
	b.WriteString("\033[38;5;220m> \033[0m"+m.input+"\033[K\n")
	b.WriteString("\033[38;5;245m Tab: switch view  /help: commands  Ctrl+C: quit\033[0m")
	fmt.Print(b.String())
}
func pad(s string,n int)string {
	if n<=0 { return "" }
	s=truncateCells(s,n)
	return s+strings.Repeat(" ",n-displayWidth(s))
}

type cellWidthState struct {
	width int
	joined bool
	regional int
}

func (s *cellWidthState) add(r rune) {
	if r==0x200d { s.joined=true; return }
	if r>=0x1f1e6 && r<=0x1f1ff {
		if s.regional%2==0 { s.width+=2 }
		s.regional++
		s.joined=false
		return
	}
	s.regional=0
	w:=runeCellWidth(r)
	if s.joined && w==2 { s.joined=false; return }
	s.joined=false
	s.width+=w
}

func displayWidth(s string) int {
	var state cellWidthState
	for _,r:=range s { state.add(r) }
	return state.width
}

func runeCellWidth(r rune) int {
	if r==0 || unicode.IsControl(r) || unicode.Is(unicode.Mn,r) || unicode.Is(unicode.Me,r) || unicode.Is(unicode.Cf,r) { return 0 }
	if r>=0xfe00 && r<=0xfe0f || r>=0xe0100 && r<=0xe01ef { return 0 }
	if isWideRune(r) { return 2 }
	return 1
}

func isWideRune(r rune) bool {
	return r>=0x1100 && r<=0x115f ||
		r==0x2329 || r==0x232a ||
		r>=0x2e80 && r<=0xa4cf && r!=0x303f ||
		r>=0xac00 && r<=0xd7a3 ||
		r>=0xf900 && r<=0xfaff ||
		r>=0xfe10 && r<=0xfe19 ||
		r>=0xfe30 && r<=0xfe6f ||
		r>=0xff00 && r<=0xff60 ||
		r>=0xffe0 && r<=0xffe6 ||
		r>=0x2600 && r<=0x27bf ||
		r>=0x1f300 && r<=0x1faff ||
		r>=0x20000 && r<=0x3fffd
}

func truncateCells(s string,limit int) string {
	if limit<=0 { return "" }
	var b strings.Builder
	var state cellWidthState
	for _,r:=range s {
		next:=state
		next.add(r)
		if next.width>limit { break }
		state=next
		b.WriteRune(r)
	}
	return b.String()
}
func max(a,b int)int{if a>b{return a};return b}

func searchLines(lines []string, query string, limit int) []string {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" || limit <= 0 { return nil }
	var matches []string
	for i, line := range lines {
		if strings.Contains(strings.ToLower(line), query) {
			matches = append(matches, fmt.Sprintf("[%d] %s", i+1, line))
			if len(matches) == limit { break }
		}
	}
	return matches
}

func reconnectDelay(attempt int) time.Duration {
	if attempt < 1 { return 0 }
	if attempt > 5 { attempt = 5 }
	return time.Duration(1<<uint(attempt-1)) * time.Second
}

func channelTargets(views []string) []string {
	var channels []string
	for _, view := range views {
		if len(view)>0 && strings.ContainsRune("#&+!",rune(view[0])) {
			channels=append(channels,view)
		}
	}
	return channels
}
