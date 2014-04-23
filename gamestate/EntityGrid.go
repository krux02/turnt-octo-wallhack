package gamestate

import (
	"github.com/krux02/turnt-octo-wallhack/math32"
)

type EntityGrid struct {
	W, H int
	data [][]*Entity
}

func CreateEntityGrid(w, h int) *EntityGrid {
	if (w & (w - 1)) != 0 {
		panic("no pow of 2 size")
	}
	if (h & (h - 1)) != 0 {
		panic("no pow of 2 size")
	}
	if w != h {
		panic("width and height needs to be equal")
	}
	return &EntityGrid{
		W:    w,
		H:    h,
		data: make([][]*Entity, w*h),
	}
}

func (this *EntityGrid) Get(x, y int) []*Entity {
	x = x & (this.W - 1)
	y = y & (this.H - 1)
	return this.data[this.W*y+x]
}

func (this *EntityGrid) Append(e *Entity) {
	x := int(math32.Floor(e.Position[0])) & (this.W - 1)
	y := int(math32.Floor(e.Position[1])) & (this.H - 1)
	index := this.W*y + x
	this.data[index] = append(this.data[index], e)
}

func remove(entities []*Entity, entity *Entity) []*Entity {
	for i, e := range entities {
		if e == entity {
			entities[i] = entities[len(entities)-1]
			return entities[:len(entities)-1]
		}
	}
	panic("arg entity not in the slice")
}

func (this *EntityGrid) Remove(e *Entity) {
	x := int(math32.Floor(e.Position[0])) & (this.W - 1)
	y := int(math32.Floor(e.Position[1])) & (this.H - 1)
	index := this.W*y + x
	this.data[index] = remove(this.data[index], e)
}
