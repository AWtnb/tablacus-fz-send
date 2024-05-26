package dir

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
)

func getChildItem(root string) (paths []string) {
	fs, err := os.ReadDir(root)
	if err != nil {
		return
	}
	for _, f := range fs {
		if strings.HasSuffix(f.Name(), ".ini") || strings.HasPrefix(f.Name(), "~$") {
			continue
		}
		p := filepath.Join(root, f.Name())
		paths = append(paths, p)
	}
	return
}

func Show(path string) {
	left := getChildItem(path)
	if len(left) < 1 {
		fmt.Printf("No items on '%s'.\n", path)
		return
	}
	if len(left) == 1 {
		fmt.Printf("Item on '%s':\n- '%s'\n", path, filepath.Base(left[0]))
		return
	}
	fmt.Printf("Items on '%s':\n", path)
	for i, l := range left {
		fmt.Printf("(%d/%d) - '%s'\n", i+1, len(left), filepath.Base(l))
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

func (d *Dir) Init(path string) {
	d.path = path
	d.member = getChildItem(d.path)
}

func (d *Dir) ExceptSelf() {
	paths := []string{}
	for _, p := range d.member {
		if p != d.path {
			paths = append(paths, p)
		}
	}
	d.member = paths
}

func (d *Dir) ExceptFiles() {
	paths := []string{}
	for _, p := range d.member {
		if fs, err := os.Stat(p); err == nil && fs.IsDir() {
			paths = append(paths, p)
		}
	}
	d.member = paths
}

func (d Dir) Member() []string {
	return d.member
}

func (d Dir) SelectItems(query string) (ps []string, err error) {
	if len(d.member) < 1 {
		return
	}
	idxs, err := fuzzyfinder.FindMulti(d.member, func(i int) string {
		return filepath.Base(d.member[i])
	}, fuzzyfinder.WithCursorPosition(fuzzyfinder.CursorPositionTop), fuzzyfinder.WithQuery(query))
	if err != nil {
		return
	}
	for _, i := range idxs {
		ps = append(ps, d.member[i])
	}
	return
}
