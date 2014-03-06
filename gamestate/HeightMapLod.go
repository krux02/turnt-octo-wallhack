package gamestate

import "github.com/krux02/turnt-octo-wallhack/helpers"

func (this *HeightMap) MaxLod() int {
	a := helpers.Log2(uint64(this.W))
	b := helpers.Log2(uint64(this.H))
	if a < b {
		return int(a)
	} else {
		return int(b)
	}
}

/*
func (this *HeightMap) BoundsLod(x, y, lod int) (min, max float32) {
	x = x & (this.W - 1)
	y = y & (this.H - 1)
	w := this.W
	h := this.H

	if lod > this.MaxLod() {
		panic("lod too high")
	}

	if lod == 0 {
		return this.Get(x, y), this.Get(x, y)
	} else {
		x, y = x>>1, y>>1
	}

	offset := 0
	for ; lod > 0; lod -= 1 {
		offset += x * y
		x, y = x>>1, y>>1
		w, h = w>>1, h>>1
		lod -= 1
	}

	min = this.MinTree[offset+w*y+x]
	max = this.MaxTree[offset+w*y+x]
	return
}
*/

/*
func (m *HeightMap) MinHm() (minHM *HeightMap, maxHM *HeightMap) {
	out = NewHeightMap(m.W/2, m.H/2)
	for x := 0; x < out.W; x++ {
		for y := 0; y < out.H; y++ {
			v := [4]float32{
				m.Get(2*x, 2*y),
				m.Get(2*x, 2*y+1),
				m.Get(2*x+1, 2*y),
				m.Get(2*x+1, 2*y+1),
			}
			min, max := v[0], v[0]
			for i := 1; i < 4; i++ {
				if v[i] < min {
					min = v[i]
				}
				if v[i] > max {
					max = v[i]
				}
			}
			minHm.Set(x, y, min)
			maxHm.Set(x, y, max)
		}
	}
}
*/
