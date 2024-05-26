package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AWtnb/tablacus-fz-send/dir"
	"github.com/AWtnb/tablacus-fz-send/sender"
	"github.com/ktr0731/go-fuzzyfinder"
)

func main() {
	var (
		src   string
		dest  string
		focus string
	)
	flag.StringVar(&src, "src", "", "location of items to copy or move")
	flag.StringVar(&dest, "dest", "", "destination to copy or move")
	flag.StringVar(&focus, "focus", "", "path of currently focusing item")
	flag.Parse()
	if len(src) < 1 {
		src = os.ExpandEnv(`C:\Users\${USERNAME}\Desktop`)
	}
	os.Exit(run(src, dest, focus))
}

func report(err error) {
	if err == fuzzyfinder.ErrAbort {
		return
	}
	fmt.Printf("ERROR: %s\n", err.Error())
	fmt.Scanln()
}

func run(src string, dest string, focus string) int {
	if src == dest {
		report(fmt.Errorf("src and dest path should be different"))
		return 1
	}
	if src == ".." {
		src = filepath.Dir(dest)
	}

	sender := sender.Sender{Src: src, Dest: dest, Focus: focus}
	err := sender.Send()
	if err != nil {
		report(err)
		return 1
	}

	dir.Show(src)
	fmt.Scanln()
	return 0
}
