package filesys

import (
	"fmt"
	"os"
	"sort"
	"strings"
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

func (es Entries) Copy(dest string) error {
	for _, ent := range es.entries {
		de := Entry{Path: dest}
		fmt.Printf("- Coping to %s: %s%s\n", de.DecoName(), ent.DecoRelPath(), ent.DecoName())
		if err := ent.CopyTo(dest); err != nil {
			return err
		}
	}
	return nil
}

func (es Entries) Remove() error {
	for _, ent := range es.Sorted() {
		fmt.Printf("- Deleting: %s%s\n", ent.DecoRelPath(), ent.DecoName())
		if err := ent.Remove(); err != nil {
			return err
		}
	}
	return nil
}
