package filesys

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"
)

func getDepth(path string) int {
	return len(strings.Split(path, string(os.PathSeparator)))
}

type Entries struct {
	entries []Entry
}

func (es *Entries) Register(path string) {
	ent := Entry{Path: path}
	es.entries = append(es.entries, ent)
}

func (es *Entries) RegisterMulti(paths []string) {
	for _, p := range paths {
		es.Register(p)
	}
}

func (es Entries) UnMovable(dest string) (paths []string) {
	for _, ent := range es.entries {
		if ent.ExistsOn(dest) {
			paths = append(paths, ent.Path)
		}
	}
	return
}

func (es Entries) Size() int {
	return len(es.entries)
}

func (es *Entries) Exclude(path string) {
	var ents []Entry
	for _, ent := range es.entries {
		if ent.Path != path {
			ents = append(ents, ent)
		}
	}
	es.entries = ents
}

func (es Entries) Sorted() []Entry {
	ss := es.entries
	sort.Slice(ss, func(i, j int) bool {
		return getDepth(ss[i].Path) > getDepth(ss[j].Path)
	})
	return ss
}

func (es Entries) Copy(src string, dest string) error {
	for _, ent := range es.entries {
		d := ent.DecoName()
		if err := ent.CopyTo(dest); err != nil {
			return err
		}
		fmt.Printf("- %s%s%s ==> %s to '%s'\n", filepath.Dir(ent.Path), string(os.PathSeparator), d, color.GreenString("Copied"), dest)
	}
	return nil
}

func (es Entries) Remove(from string) error {
	for _, ent := range es.Sorted() {
		d := ent.DecoName()
		if err := ent.Remove(); err != nil {
			return err
		}
		fmt.Printf("- %s ==> %s from '%s'\n", d, color.HiMagentaString("Deleted"), filepath.Dir(ent.Path))
	}
	return nil
}
