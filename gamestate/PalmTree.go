package gamestate

import (
	mgl "github.com/Jragonmiris/mathgl"
)

// instance data for each tree
type PalmTree struct {
	Position_ws mgl.Vec4f
}

// forest
type PalmTreesInstanceData struct {
	Positions []PalmTree
}
