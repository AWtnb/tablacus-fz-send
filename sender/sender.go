package sender

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AWtnb/go-asker"
	"github.com/AWtnb/tablacus-fz-send/dir"
	"github.com/AWtnb/tablacus-fz-send/filesys"
	"github.com/fatih/color"
	"github.com/ktr0731/go-fuzzyfinder"
)

func finishLabel(message string) {
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Printf("%s %s\n", green("[DONE]"), message)
}

var (
	ErrNoSubDir    = errors.New("no subdir to move")
	ErrInvalidDest = errors.New("invalid dest path")
)

type Sender struct {
	Src   string
	Dest  string
	Focus string
}

func (sdr Sender) isDisposal() bool {
	return sdr.Dest == "_obsolete"
}

func (sdr Sender) Targets() ([]string, error) {
	var d dir.Dir
	d.Init(sdr.Src)
	if fs, err := os.Stat(sdr.Dest); err == nil && fs.IsDir() {
		d.Except(sdr.Dest)
	}
	q := ""
	if 0 < len(sdr.Focus) {
		q = filepath.Base(sdr.Focus)
	}
	return d.SelectItems(q)
}

func (sdr Sender) DestPath() (string, error) {
	if fs, err := os.Stat(sdr.Dest); err == nil && fs.IsDir() {
		return sdr.Dest, nil
	}
	if strings.Contains(sdr.Dest, string(os.PathSeparator)) {
		return "", ErrInvalidDest
	}
	if sdr.isDisposal() {
		return dir.Create(sdr.Src, sdr.Dest)
	}
	if len(sdr.Dest) < 1 {
		var dd dir.Dir
		dd.Init(sdr.Src)
		dd.ExceptFiles()
		dd.ExceptSelf()
		sds := dd.Member()
		if len(sds) < 1 {
			return "", ErrNoSubDir
		}
		idx, err := fuzzyfinder.Find(sds, func(i int) string {
			return filepath.Base(sds[i])
		})
		return sds[idx], err
	}
	p := filepath.Join(sdr.Src, sdr.Dest)
	if fs, err := os.Stat(p); err == nil && fs.IsDir() {
		return p, nil
	}
	return dir.Create(sdr.Src, sdr.Dest)
}

func (sdr Sender) sendItems(paths []string, dest string) error {
	var fes filesys.Entries
	fes.RegisterMulti(paths)
	dupls := fes.UnMovable(dest)
	if 0 < len(dupls) {
		for _, dp := range dupls {
			a := asker.Asker{Accept: "y", Reject: "n"}
			a.Ask(fmt.Sprintf("Name duplicated: '%s'\noverwrite?", filepath.Base(dp)))
			if !a.Accepted() {
				fmt.Println("==> skipped")
				fes.Exclude(dp)
			}
		}
	}
	if fes.Size() < 1 {
		return nil
	}
	if err := fes.CopyTo(dest); err != nil {
		return err
	}
	finishLabel("successfully copied everything")
	fes.Show()

	if sdr.isDisposal() {
		finishLabel("removed original items")
		return fes.Remove()
	}

	a := asker.Asker{Accept: "y", Reject: "n"}
	a.Ask("\n==> Delete original?")
	if a.Accepted() {
		if err := fes.Remove(); err != nil {
			return err
		}
	}
	return nil
}

func (sdr Sender) Send() error {
	t, err := sdr.Targets()
	if err != nil {
		return err
	}

	d, err := sdr.DestPath()
	if err != nil {
		return err
	}

	if err := sdr.sendItems(t, d); err != nil {
		return err
	}

	cyan := color.New(color.BgCyan)
	cyan.Println("\n[FINISHED]")
	return nil
}
