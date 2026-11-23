//go:build linux

package app

import (
	"os"
	"syscall"
	"unsafe"
)

type terminalState struct{ old syscall.Termios }
func terminalAvailable(fd int)bool{var t syscall.Termios;_,_,e:=syscall.Syscall(syscall.SYS_IOCTL,uintptr(fd),uintptr(syscall.TCGETS),uintptr(unsafe.Pointer(&t)));return e==0}
func makeRaw(fd int)(*terminalState,error){
	var old syscall.Termios
	_,_,e:=syscall.Syscall(syscall.SYS_IOCTL,uintptr(fd),uintptr(syscall.TCGETS),uintptr(unsafe.Pointer(&old)))
	if e!=0{return nil,e}
	raw:=old
	raw.Iflag &^= syscall.BRKINT|syscall.ICRNL|syscall.INPCK|syscall.ISTRIP|syscall.IXON
	raw.Oflag &^= syscall.OPOST
	raw.Cflag |= syscall.CS8
	raw.Lflag &^= syscall.ECHO|syscall.ICANON|syscall.IEXTEN|syscall.ISIG
	raw.Cc[syscall.VMIN]=1;raw.Cc[syscall.VTIME]=0
	_,_,e=syscall.Syscall(syscall.SYS_IOCTL,uintptr(fd),uintptr(syscall.TCSETS),uintptr(unsafe.Pointer(&raw)))
	if e!=0{return nil,e}
	return &terminalState{old:old},nil
}
func(s *terminalState)restore(fd int){_,_,_=syscall.Syscall(syscall.SYS_IOCTL,uintptr(fd),uintptr(syscall.TCSETS),uintptr(unsafe.Pointer(&s.old)))}
func terminalSize()(int,int){
	var ws struct{Row,Col,Xpixel,Ypixel uint16}
	_,_,e:=syscall.Syscall(syscall.SYS_IOCTL,uintptr(os.Stdout.Fd()),uintptr(syscall.TIOCGWINSZ),uintptr(unsafe.Pointer(&ws)))
	if e!=0||ws.Col==0{return 80,24};return int(ws.Col),int(ws.Row)
}
