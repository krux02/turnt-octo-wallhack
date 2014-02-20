package gamestate

import (
	mgl "github.com/Jragonmiris/mathgl"
	"math/rand"
	"testing"
)

type MyPos struct {
	pos mgl.Vec3f
}

func (mp MyPos) Dimension(dim int) float32 {
	return mp.pos[dim]
}

func randPos() MyPos {
	return MyPos{mgl.Vec3f{rand.Float32(), rand.Float32(), rand.Float32()}}
}

const length = 200

func TestKdTree(t *testing.T) {
	positions := make([]KdElement, length)
	for i := 0; i < length; i++ {
		positions[i] = randPos()
	}

	center := mgl.Vec3f{0.5, 0.5, 0.5}
	nearestPos := positions[0]
	dist := nearestPos.(MyPos).pos.Sub(center).Len()
	for i := 1; i < length; i++ {
		dist2 := positions[i].(MyPos).pos.Sub(center).Len()
		if dist2 < dist {
			dist = dist2
			nearestPos = positions[i]
		}
	}

	tree := NewTree(positions)

	nearestPosTree := tree.NearestQuery(MyPos{center}, func(kd KdElement) bool { return true })

	if nearestPos != nearestPosTree {
		t.Log(nearestPos)
		t.Log(nearestPosTree)
		t.Fail()
	}
}
