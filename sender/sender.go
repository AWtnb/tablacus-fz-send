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
	"github.com/ktr0731/go-fuzzyfinder"
)

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

func (sdr Sender) targets() ([]string, error) {
	var d dir.Dir
	d.Init(sdr.Src, true)
	if fs, err := os.Stat(sdr.Dest); err == nil && fs.IsDir() {
		d.Except(sdr.Dest)
	}
	q := ""
	if 0 < len(sdr.Focus) {
		q = filepath.Base(sdr.Focus)
	}
	return d.SelectItems(q)
}

func (sdr Sender) destPath() (string, error) {
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
		dd.Init(sdr.Src, false)
		sds := dd.Member()
		if len(sds) < 1 {
			return "", ErrNoSubDir
		}
		idx, err := fuzzyfinder.Find(sds, func(i int) string {
			p := sds[i]
			rel, _ := filepath.Rel(sdr.Src, p)
			return filepath.Base(rel)
		})
		return sds[idx], err
	}
	return dir.Create(sdr.Src, sdr.Dest)
}

func (sdr Sender) sendItems(paths []string, dest string) error {
	var fes filesys.Entries
	fes.RegisterMulti(paths)
	dupls := fes.UnMovable(dest)
	for _, dp := range dupls {
		a := asker.Asker{Accept: "y", Reject: "n"}
		e := filesys.Entry{Path: dp}
		d := filesys.Entry{Path: dest}
		a.Ask(fmt.Sprintf("Name duplicated: %s in %s\nOverwrite?", e.DecoName(), d.DecoName()))
		if !a.Accepted() {
			fmt.Println("==> Skipped")
			fes.Exclude(dp)
		}
	}
	if fes.Size() < 1 {
		return nil
	}
	if err := fes.Copy(dest); err != nil {
		return err
	}

	if sdr.isDisposal() {
		return fes.Remove()
	}

	a := asker.Asker{Accept: "y", Reject: "n"}
	a.Ask("Delete original?")
	if a.Accepted() {
		if err := fes.Remove(); err != nil {
			return err
		}
	}
	return nil
}

func (sdr Sender) Send() error {
	ts, err := sdr.targets()
	if err != nil {
		return err
	}

	d, err := sdr.destPath()
	if err != nil {
		return err
	}

	if err := sdr.sendItems(ts, d); err != nil {
		return err
	}

	return nil
}
