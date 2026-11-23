package irc

import (
	"bufio"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
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
	Proxy string
	Channels []string
}

type Event struct { Kind, Target, Nick, Text string }
type Client struct {
	cfg Config
	conn net.Conn
	write chan string
	events chan Event
	done chan struct{}
	once sync.Once
}
func NewClient(cfg Config) *Client {
	return &Client{cfg:cfg,write:make(chan string,64),events:make(chan Event,256),done:make(chan struct{})}
}
func (c *Client) Events() <-chan Event { return c.events }
func (c *Client) Done() <-chan struct{} { return c.done }

func (c *Client) Connect() error {
	d:=&net.Dialer{Timeout:12*time.Second}
	var conn net.Conn
	var err error
	if c.cfg.Proxy!="" {
		conn,err=dialSOCKS5(d,c.cfg.Proxy,c.cfg.Address)
	} else {
		conn,err=d.Dial("tcp",c.cfg.Address)
	}
	if err!=nil {return fmt.Errorf("network route failed: %w",err)}
	if c.cfg.TLS {
		host,_,splitErr:=net.SplitHostPort(c.cfg.Address)
		if splitErr!=nil {_=conn.Close();return fmt.Errorf("server must be host:port: %w",splitErr)}
		tlsConn:=tls.Client(conn,&tls.Config{MinVersion:tls.VersionTLS12,ServerName:host})
		_ = tlsConn.SetDeadline(time.Now().Add(12*time.Second))
		if err=tlsConn.Handshake(); err!=nil {_=conn.Close();return fmt.Errorf("TLS handshake failed: %w",err)}
		_ = tlsConn.SetDeadline(time.Time{})
		conn=tlsConn
	}
	c.conn=conn
	go c.writeLoop()
	go c.readLoop()
	c.Send("NICK",c.cfg.Nick)
	c.Send("USER",c.cfg.Nick,"0","*",c.cfg.RealName)
	return nil
}

// dialSOCKS5 sends the destination hostname to the proxy for remote DNS resolution.
func dialSOCKS5(d *net.Dialer, proxyAddr, targetAddr string) (net.Conn,error) {
	conn,err:=d.Dial("tcp",proxyAddr)
	if err!=nil {return nil,fmt.Errorf("connect to SOCKS5 proxy: %w",err)}
	fail:=func(e error)(net.Conn,error){_ = conn.Close();return nil,e}
	_ = conn.SetDeadline(time.Now().Add(12*time.Second))
	if _,err=conn.Write([]byte{5,1,0}); err!=nil{return fail(err)}
	reply:=make([]byte,2)
	if _,err=io.ReadFull(conn,reply); err!=nil{return fail(err)}
	if reply[0]!=5 || reply[1]!=0{return fail(fmt.Errorf("SOCKS5 proxy rejected unauthenticated method"))}
	host,portText,err:=net.SplitHostPort(targetAddr)
	if err!=nil{return fail(fmt.Errorf("target must be host:port: %w",err))}
	port,err:=net.LookupPort("tcp",portText)
	if err!=nil{return fail(fmt.Errorf("invalid target port: %w",err))}
	hostBytes:=[]byte(host)
	if len(hostBytes)==0 || len(hostBytes)>255{return fail(fmt.Errorf("invalid target hostname length"))}
	req:=[]byte{5,1,0,3,byte(len(hostBytes))}
	req=append(req,hostBytes...)
	p:=make([]byte,2);binary.BigEndian.PutUint16(p,uint16(port));req=append(req,p...)
	if _,err=conn.Write(req); err!=nil{return fail(err)}
	head:=make([]byte,4)
	if _,err=io.ReadFull(conn,head); err!=nil{return fail(err)}
	if head[0]!=5 || head[1]!=0{return fail(fmt.Errorf("SOCKS5 CONNECT rejected with code %d",head[1]))}
	var skip int
	switch head[3] {
	case 1: skip=4
	case 4: skip=16
	case 3:
		n:=[]byte{0};if _,err=io.ReadFull(conn,n);err!=nil{return fail(err)};skip=int(n[0])
	default:return fail(fmt.Errorf("invalid SOCKS5 address type"))
	}
	if _,err=io.CopyN(io.Discard,conn,int64(skip+2));err!=nil{return fail(err)}
	_ = conn.SetDeadline(time.Time{})
	return conn,nil
}
func (c *Client) Send(command string,params ...string) {
	line:=Format(command,params...)
	if Validate(line)!=nil{return}
	select{case c.write<-line:case <-c.done:}
}
func (c *Client) Close(){c.once.Do(func(){select{case c.write<-"QUIT :Leaving ZIRC":default:};if c.conn!=nil{_=c.conn.Close()};close(c.done)})}
func (c *Client) writeLoop(){for{select{case line:=<-c.write:
	_=c.conn.SetWriteDeadline(time.Now().Add(8*time.Second))
	if _,err:=fmt.Fprintf(c.conn,"%s\r\n",line);err!=nil{c.emit(Event{Kind:"error",Text:err.Error()});c.Close();return}
	case <-c.done:return}}}
func (c *Client) readLoop(){
	defer c.Close()
	sc:=bufio.NewScanner(c.conn);sc.Buffer(make([]byte,4096),8192)
	for sc.Scan(){
		raw:=sc.Text();m:=ParseLine(raw)
		if m.Command=="PING"{c.Send("PONG",m.Trailing);continue}
		nick:=m.Prefix;if i:=strings.IndexAny(nick,"!@");i>=0{nick=nick[:i]}
		switch m.Command{
		case "001":
			label:="Connected"
			if c.cfg.TLS{label+=" over verified TLS"}
			if c.cfg.Proxy!=""{label+=" through configured SOCKS5 proxy"}
			c.emit(Event{Kind:"status",Text:label+" as "+c.cfg.Nick})
			for _,ch:=range c.cfg.Channels{c.Send("JOIN",ch)}
		case "PRIVMSG":if len(m.Params)>0{c.emit(Event{Kind:"message",Target:m.Params[0],Nick:nick,Text:m.Trailing})}
		case "NOTICE","JOIN","PART","QUIT","ERROR":c.emit(Event{Kind:"notice",Text:raw})
		default:if strings.HasPrefix(m.Command,"4")||strings.HasPrefix(m.Command,"5"){c.emit(Event{Kind:"notice",Text:raw)}}
		}
	}
	if err:=sc.Err();err!=nil{c.emit(Event{Kind:"error",Text:err.Error()})}
	c.emit(Event{Kind:"status",Text:"Disconnected"})
}
func (c *Client) emit(e Event){select{case c.events<-e:default:}}
