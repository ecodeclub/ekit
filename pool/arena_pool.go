package pool

import (
	"arena"
)

type ArenaPool struct {
	New func() any

	Chain  []*Box
	Cursor int
	memMap map[string]*arena.Arena
}

type Box struct {
	Mem    *arena.Arena
	Object any
}

func NewArenaPool(new func() any) *ArenaPool {
	return &ArenaPool{
		New:    new,
		Cursor: -1,
	}
}

func (a *ArenaPool) NewX() (*Box, error) {
	mem := arena.NewArena()
	box := arena.New[Box](mem)
	box.Mem = mem
	box.Object = a.New()
	return box, nil
}

func (a *ArenaPool) Get() (*Box, error) {
	if a.Cursor == -1 {
		X, err := a.NewX()
		if err != nil {
			return nil, err
		}
		return X, nil
	}
	ret := a.Chain[a.Cursor]
	a.Cursor--
	return ret, nil
}

func (a *ArenaPool) Put(X *Box) error {
	a.Chain = append(a.Chain, X)
	a.Cursor++
	return nil
}

func (b *Box) Free() error {
	b.Mem.Free()
	return nil
}
