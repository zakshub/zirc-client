//go:build !linux

package app

import "fmt"

type terminalState struct{}
func terminalAvailable(fd int)bool{return false}
func makeRaw(fd int)(*terminalState,error){return nil,fmt.Errorf("ZIRC initial TUI currently supports Linux only")}
func(s *terminalState)restore(fd int){}
func terminalSize()(int,int){return 80,24}
