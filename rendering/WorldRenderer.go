package rendering

import (
	mgl "github.com/krux02/mathgl"
	"github.com/krux02/turnt-octo-wallhack/particles"
	"github.com/krux02/turnt-octo-wallhack/world"
)

type WorldRenderer struct {
	HeightMapRenderer *HeightMapRenderer
	MeshRenderer      *MeshRenderer
	Portal            MeshRenderData
	PalmTrees         *PalmTrees
	ParticleSystem    *particles.ParticleSystem
}

func NewWorldRenderer(w *world.World) *WorldRenderer {
	portalData := w.Portals[0].Mesh
	mr := NewMeshRenderer()
	return &WorldRenderer{
		HeightMapRenderer: NewHeightMapRenderer(w.HeightMap),
		MeshRenderer:      mr,
		Portal:            mr.CreateMeshRenderData(portalData),
		PalmTrees:         NewPalmTrees(w.HeightMap, 10000),
		ParticleSystem:    particles.NewParticleSystem(w, 75000, mgl.Vec3f{32, 32, 32}, 1, 250),
	}
}

func (this *WorldRenderer) Render(w *world.World) {
}
