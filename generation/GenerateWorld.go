package generation

import (
	mgl "github.com/Jragonmiris/mathgl"
	gs "github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	//"image"
	"math/rand"
)

func GenerateWorld(W, H, N int) (world *gs.World) {
	heights := gs.NewHeightMap(W, H)
	DiamondSquare(heights, float32(W))
	Portals := RandomPortals(heights, N)

	kdTree := make([]gs.KdElement, len(Portals))
	for i, portal := range Portals {
		kdTree[i] = portal
	}
	kdTree = gs.NewTree(kdTree)

	world = &gs.World{heights, kdTree, Portals}
	return
}

func RandomPortals(hm *gs.HeightMap, N int) []*gs.Portal {
	if N&1 != 0 {
		panic("not an even number of portals")
	}

	PortalMesh := gs.QuadMesh()
	normal := mgl.Vec4f{0, 0, 1, 0}

	Portals := make([]*gs.Portal, N)
	for i := 0; i < N; i++ {
		x := rand.Float32() * float32(hm.W)
		y := rand.Float32() * float32(hm.H)
		z := hm.Get2f(x, y) + 5
		q := mgl.Quatf{rand.Float32(), mgl.Vec3f{rand.Float32(), rand.Float32(), rand.Float32()}}.Normalize()
		Portals[i] = &gs.Portal{
			gs.Entity{mgl.Vec4f{x, y, z, 1}, q},
			PortalMesh,
			normal,
			nil,
		}
	}
	for i := range Portals {
		j := i ^ 1
		Portals[i].Target = Portals[j]
	}
	return Portals
}

func DiamondSquare(m *gs.HeightMap, factor float32) {
	w, h := m.W, m.H

	stepSize := w

	squares := func() {
		for i := 0; i < w; i += stepSize {
			for j := 0; j < h; j += stepSize {
				sum := m.Get(i, j)
				sum += m.Get(i+stepSize, j)
				sum += m.Get(i, j+stepSize)
				sum += m.Get(i+stepSize, j+stepSize)

				h := (sum / 4.0) + (rand.Float32()-0.5)*factor
				x := i + stepSize/2
				y := j + stepSize/2
				m.Set(x, y, h)
			}
		}
	}

	diamonds := func() {
		for i := 0; i <= w; i += stepSize {
			for j := 0; j <= h; j += stepSize {
				if ((i+j)/stepSize)%2 == 1 {
					sum := float32(0)
					count := float32(0.0)
					if i != 0 {
						sum += m.Get(i-stepSize, j)
						count += 1
					}
					if i < w-1 {
						sum += m.Get(i+stepSize, j)
						count += 1
					}
					if j != 0 {
						sum += m.Get(i, j-stepSize)
						count += 1
					}
					if j < h-1 {
						sum += m.Get(i, j+stepSize)
						count += 1
					}

					h := (sum / count) + (rand.Float32()-0.5)*factor
					m.Set(i, j, h)
				}
			}
		}
	}

	for stepSize > 0 {
		squares()

		stepSize /= 2
		if stepSize == 0 {
			break
		}
		diamonds()

		factor *= 0.5
	}
}

func GenerateWaterTextures() {
	hm1 := gs.NewHeightMap(64, 64)
	hm2 := gs.NewHeightMap(64, 64)
	DiamondSquare(hm1, 1)
	DiamondSquare(hm2, 1)

	helpers.SaveImage("hm1.png", hm1.ExportImage())
	helpers.SaveImage("hm2.png", hm2.ExportImage())
}
