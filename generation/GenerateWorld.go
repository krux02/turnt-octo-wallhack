package generation

import (
	"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	gs "github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	"github.com/krux02/turnt-octo-wallhack/math32"
	//"image"
	"math/rand"
)

func sigmuid(x float32) float32 {
	return ((x / math32.Sqrt(1+x*x)) + 1) / 2
}

func r() float32 {
	return rand.Float32()
}

func GenerateWorld(W, H, N int) (world *gs.World) {
	heights := gs.NewHeightMap(W, H)
	DiamondSquare(heights, float32(W))
	minH, maxH := heights.Bounds()

	fmt.Println(minH, maxH)
	water := &gs.Water{
		W:          W,
		H:          H,
		LowerBound: mgl.Vec3f{0, 0, minH},
		UpperBound: mgl.Vec3f{float32(W), float32(H), maxH},
	}
	Portals := RandomPortals(heights, N)
	for i, x := range heights.Data {
		heights.Data[i] = math32.Mix(0, x-5, sigmuid((x-5)/10)) + 5
	}
	kdTree := make([]gs.KdElement, len(Portals))
	for i, portal := range Portals {
		kdTree[i] = portal
	}
	kdTree = gs.NewTree(kdTree)

	forest := GeneratePalmTrees(heights, 5000)
	npcs := make([]gs.IRenderEntity, 255)
	for i := 0; i < 255; i++ {
		x := r() * float32(W)
		y := r() * float32(H)
		h := heights.Get2f(x, y) + 1
		npcs[i] = &gs.Npc{mgl.Vec4f{x, y, h, 1}, mgl.QuatIdentf()}
	}
	startPos := mgl.Vec4f{5, 5, 10, 1}
	player := &gs.Player{Camera: *gs.NewCameraFromPos4f(startPos)}
	world = &gs.World{heights, water, kdTree, Portals, forest, npcs, player}
	return
}

func RandomPortals(hm *gs.HeightMap, N int) []*gs.Portal {
	if N&1 != 0 {
		panic("not an even number of portals")
	}

	normal := mgl.Vec4f{0, 0, 1, 0}
	mesh := gs.QuadMesh()

	Portals := make([]*gs.Portal, N)
	for i := 0; i < N; i++ {
		x := rand.Float32() * float32(hm.W)
		y := rand.Float32() * float32(hm.H)
		z := hm.Get2f(x, y) + 5
		q := mgl.Quatf{r(), mgl.Vec3f{r(), r(), r()}}.Normalize()
		Portals[i] = &gs.Portal{
			gs.Entity{mgl.Vec4f{x, y, z, 1}, q},
			normal,
			nil,
			mesh,
		}
	}
	for i := range Portals {
		j := i ^ 1
		Portals[i].Target = Portals[j]
	}
	return Portals
}

func GeneratePalmTrees(hm *gs.HeightMap, count int) *gs.Forest {
	trees := make([]gs.PalmTree, count)

	for i := 0; i < count; i++ {

		var x, y float32
		for true {
			x = rand.Float32() * float32(hm.W)
			y = rand.Float32() * float32(hm.H)
			if hm.Normal2f(x, y)[2] > 0.65 && hm.Get2f(x, y) > 10 {
				break
			}
		}

		z := hm.Get2f(x, y)
		trees[i] = gs.PalmTree{mgl.Vec4f{x, y, z, 1}}
	}

	forest := new(gs.Forest)
	forest.Positions = trees
	forest.Model = mgl.Ident4f()
	return forest
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
