package irc

import (
	"encoding/binary"
	"io"
	"net"
	"testing"
)

func TestSOCKS5UsesProxySideDNS(t *testing.T) {
	listener,err:=net.Listen("tcp","127.0.0.1:0")
	if err!=nil {t.Fatal(err)}
	defer listener.Close()
	serverErr:=make(chan error,1)
	go func(){
		conn,e:=listener.Accept()
		if e!=nil {serverErr<-e;return}
		defer conn.Close()
		greeting:=make([]byte,3)
		if _,e=io.ReadFull(conn,greeting);e!=nil{serverErr<-e;return}
		if greeting[0]!=5||greeting[1]!=1||greeting[2]!=0{serverErr<-io.ErrUnexpectedEOF;return}
		_,e=conn.Write([]byte{5,0});if e!=nil{serverErr<-e;return}
		head:=make([]byte,5)
		if _,e=io.ReadFull(conn,head);e!=nil{serverErr<-e;return}
		if head[0]!=5||head[1]!=1||head[3]!=3{serverErr<-io.ErrUnexpectedEOF;return}
		name:=make([]byte,int(head[4]))
		if _,e=io.ReadFull(conn,name);e!=nil{serverErr<-e;return}
		port:=make([]byte,2)
		if _,e=io.ReadFull(conn,port);e!=nil{serverErr<-e;return}
		if string(name)!="irc.example.invalid"{serverErr<-io.ErrUnexpectedEOF;return}
		if binary.BigEndian.Uint16(port)!=6697{serverErr<-io.ErrUnexpectedEOF;return}
		_,e=conn.Write([]byte{5,0,0,1,127,0,0,1,0,1})
		serverErr<-e
	}()
	conn,err:=dialSOCKS5(&net.Dialer{},listener.Addr().String(),"irc.example.invalid:6697")
	if err!=nil{t.Fatal(err)}
	_ = conn.Close()
	if err=<-serverErr;err!=nil{t.Fatal(err)}
}
