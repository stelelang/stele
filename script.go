package stele

import "deedles.dev/mk"

type Script struct {
	Declarations Declarations
}

type Declarations struct {
	m map[string]any
}

func (d *Declarations) init() {
	if d.m == nil {
		mk.Map(&d.m, 0)
	}
}

func (d *Declarations) Add(id string, val any) bool {
	d.init()

	_, ok := d.m[id]
	if ok {
		return false
	}

	d.m[id] = val
	return true
}
