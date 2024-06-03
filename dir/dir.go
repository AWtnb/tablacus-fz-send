package dir

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/AWtnb/go-walk"
	"github.com/AWtnb/tablacus-fz-send/filesys"
	"github.com/ktr0731/go-fuzzyfinder"
)

var ErrNoItem = errors.New("no item to send")

func getChildItem(root string, depth int, all bool) (paths []string) {
	var d walk.Dir
	d.Init(root, all, depth, "")
	found, err := d.GetChildItem()
	if err != nil {
		return
	}
	for _, f := range found {
		n := filepath.Base(f)
		if f == root || strings.HasSuffix(n, ".ini") || strings.HasPrefix(n, "~$") {
			continue
		}
		paths = append(paths, f)
	}
	return
}

func Show(path string) {
	pe := filesys.Entry{Path: path}
	dp := pe.DecoName()
	left := getChildItem(path, 1, true)
	if len(left) < 1 {
		fmt.Printf("(now %s is empty)\n", dp)
		return
	}
	if len(left) == 1 {
		fmt.Printf("Left item on %s:\n", dp)
		e := filesys.Entry{Path: left[0]}
		fmt.Printf(" - %s\n", e.DecoName())
		return
	}
	fmt.Printf("Left items on %s:\n", dp)
	for i, p := range left {
		e := filesys.Entry{Path: p}
		fmt.Printf("- %s %s\n", filesys.PadCount(i+1, len(left)), e.DecoName())
	}
}

func getPerm(path string) fs.FileMode {
	s := string(os.PathSeparator)
	elems := strings.Split(path, s)
	for i := 0; i < len(elems); i++ {
		ln := len(elems) - i
		p := strings.Join(elems[0:ln], s)
		if fs, err := os.Stat(p); err == nil && fs.IsDir() {
			return fs.Mode() & os.ModePerm
		}
	}
	return 0700
}

func Create(dir string, name string) (string, error) {
	p := filepath.Join(dir, name)
	if fs, err := os.Stat(p); err == nil && fs.IsDir() {
		return p, nil
	}
	return p, os.MkdirAll(p, getPerm(p))
}

type Dir struct {
	path   string
	member []string
}

func (d *Dir) Init(path string, all bool) {
	d.path = path
	d.member = getChildItem(d.path, -1, all)
}

func (d *Dir) Except(path string) {
	paths := []string{}
	for _, p := range d.member {
		if strings.HasPrefix(p, path) {
			continue
		}
		paths = append(paths, p)
	}
	d.member = paths
}

func (d *Dir) ExceptDir() {
	paths := []string{}
	for _, p := range d.member {
		if fs, err := os.Stat(p); err == nil && !fs.IsDir() {
			paths = append(paths, p)
		}
	}
	d.member = paths
}

func (d Dir) Member() []string {
	return d.member
}

func (d Dir) SelectItems() (ps []string, err error) {
	if len(d.member) < 1 {
		err = ErrNoItem
		return
	}
	idxs, err := fuzzyfinder.FindMulti(d.member, func(i int) string {
		p := d.member[i]
		rel, _ := filepath.Rel(d.path, p)
		if fs, err := os.Stat(p); err == nil && fs.IsDir() {
			return fmt.Sprintf("%s \U0001F4C1", rel)
		}
		return rel
	}, fuzzyfinder.WithCursorPosition(fuzzyfinder.CursorPositionTop))
	if err != nil {
		return
	}
	for _, i := range idxs {
		ps = append(ps, d.member[i])
	}
	return
}
