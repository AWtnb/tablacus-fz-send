package filesys

import (
	"fmt"

	"github.com/fatih/color"
)

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

func (es Entries) CopyTo(dest string) error {
	for _, ent := range es.entries {
		d := ent.DecoName(false)
		if err := ent.CopyTo(dest); err != nil {
			return err
		}
		fmt.Printf("- %s ==> %s\n", d, color.GreenString("Copied"))
	}
	return nil
}

func (es Entries) Remove() error {
	for _, ent := range es.entries {
		d := ent.DecoName(true)
		if err := ent.Remove(); err != nil {
			return err
		}
		fmt.Printf("- %s ==> %s\n", d, color.MagentaString("Deleted"))
	}
	return nil
}
