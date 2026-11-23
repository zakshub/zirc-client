package main

import (
	"flag"
	"fmt"
	"os"

	fyneapp "fyne.io/fyne/v2/app"
	"github.com/zakshub/zirc-client/internal/desktop"
)

var version = "dev"

func main() {
	showVersion := flag.Bool("version", false, "Print ZIRC desktop version")
	flag.Parse()
	if *showVersion {
		fmt.Println("ZIRC Desktop", version)
		return
	}

	a := fyneapp.NewWithID("com.zakshub.zirc")
	w := a.NewWindow("ZIRC")
	w.Resize(fyneSize(1280, 820))
	desktop.New(w).Show()
	a.Run()
}

func fyneSize(width, height float32) fyne.Size {
	return fyne.NewSize(width, height)
}
