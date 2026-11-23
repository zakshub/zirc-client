package irc

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type Config struct {
	Address string
	Nick string
	RealName string
	TLS bool
	Channels []string
}

type Event struct {
	Kind string
	Target string
	Nick string
	Text string
}

type Client struct {
	cfg Config
	conn net.Conn
	write chan string
	events chan Event
	done chan struct{}
	once sync.Once
}

func NewClient(cfg Config) *Client {
	return &Client{cfg:cfg, write:make(chan string,64), events:make(chan Event,256), done:make(chan struct{})}
}
func (c *Client) Events() <-chan Event { return c.events }
func (c *Client) Done() <-chan struct{} { return c.done }

func (c *Client) Connect() error {
	d := &net.Dialer{Timeout:12*time.Second}
	var conn net.Conn
	var err error
	if c.cfg.TLS {
		host,_,splitErr := net.SplitHostPort(c.cfg.Address)
		if splitErr != nil { return fmt.Errorf("server must be host:port: %w",splitErr) }
		conn,err = tls.DialWithDialer(d,"tcp",c.cfg.Address,&tls.Config{MinVersion:tls.VersionTLS12,ServerName:host})
	} else {
		conn,err = d.Dial("tcp",c.cfg.Address)
	}
	if err != nil { return fmt.Errorf("connect failed: %w",err) }
	c.conn=conn
	go c.writeLoop()
	go c.readLoop()
	c.Send("NICK",c.cfg.Nick)
	c.Send("USER",c.cfg.Nick,"0","*",c.cfg.RealName)
	return nil
}
func (c *Client) Send(command string, params ...string) {
	line:=Format(command,params...)
	if Validate(line)!=nil { return }
	select { case c.write<-line: case <-c.done: }
}
func (c *Client) Close() {
	c.once.Do(func(){
		select {case c.write<-"QUIT :Leaving ZIRC":default:}
		if c.conn!=nil { _=c.conn.Close() }
		close(c.done)
	})
}
func (c *Client) writeLoop() {
	for {
		select {
		case line:=<-c.write:
			_ = c.conn.SetWriteDeadline(time.Now().Add(8*time.Second))
			if _,err:=fmt.Fprintf(c.conn,"%s\r\n",line); err!=nil { c.emit(Event{Kind:"error",Text:err.Error()}); c.Close(); return }
		case <-c.done:return
		}
	}
}
func (c *Client) readLoop() {
	defer c.Close()
	sc:=bufio.NewScanner(c.conn)
	sc.Buffer(make([]byte,4096),8192)
	for sc.Scan() {
		raw:=sc.Text()
		m:=ParseLine(raw)
		if m.Command=="PING" { c.Send("PONG",m.Trailing); continue }
		nick:=m.Prefix
		if i:=strings.IndexAny(nick,"!@"); i>=0 {nick=nick[:i]}
		switch m.Command {
		case "001":
			c.emit(Event{Kind:"status",Text:"Connected securely as "+c.cfg.Nick})
			for _,ch:=range c.cfg.Channels {c.Send("JOIN",ch)}
		case "PRIVMSG":
			if len(m.Params)>0 {c.emit(Event{Kind:"message",Target:m.Params[0],Nick:nick,Text:m.Trailing})}
		case "NOTICE","JOIN","PART","QUIT","ERROR":
			c.emit(Event{Kind:"notice",Text:raw})
		default:
			if strings.HasPrefix(m.Command,"4") || strings.HasPrefix(m.Command,"5") {c.emit(Event{Kind:"notice",Text:raw})}
		}
	}
	if err:=sc.Err(); err!=nil {c.emit(Event{Kind:"error",Text:err.Error()})}
	c.emit(Event{Kind:"status",Text:"Disconnected"})
}
func (c *Client) emit(e Event) {select {case c.events<-e:default:}}
