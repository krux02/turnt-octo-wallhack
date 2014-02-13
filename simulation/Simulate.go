package simulation

import (
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/world"
)

func EnterPortal(portal *world.Portal, camera *gamestate.Camera) {
	camera.SetFromMat4(portal.Target.ModelMat4().Mul4(portal.ModelMat4().Inv()).Mul4(camera.View()))
}

func PortalPassed(portal *world.Portal, pos1, pos2 mgl.Vec3f) bool {
	pos1w := mgl.Vec4f{pos1[0], pos1[1], pos1[2], 1}
	pos2w := mgl.Vec4f{pos2[0], pos2[1], pos2[2], 1}
	plane := portal.ClippingPlane(true)
	return pos1w.Dot(plane)*pos2w.Dot(plane) < 0
}

func Simulate(gs *gamestate.GameState) {
	UpdatePlayer(gs.Player, gs)
}

func UpdatePlayer(p *gamestate.Player, gs *gamestate.GameState) {
	rot := p.Input.Rotate
	p.Camera.Yaw(rot[0])
	p.Camera.Pitch(rot[1])
	p.Camera.Roll(rot[2])

	move := p.Input.Move
	if move.Len() > 0 {
		move = move.Normalize()
	}
	move = p.Camera.Direction.Rotate(move)

	if gs.Options.NoPlayerPhysics {
		move = move.Mul(0.1)
		p.Velocity = move
		p.Camera.MoveAbsolute(move)
	} else {
		move = move.Mul(0.01)
		p.Velocity = p.Velocity.Add(move)
		p.Camera.MoveAbsolute(p.Velocity)

		groundHeight := gs.World.HeightMap.Get2f(p.Position()[0], p.Position()[1])

		height := p.Camera.Position[2]
		minHeight := groundHeight + 1.5
		maxHeight := groundHeight + 20

		if height < minHeight {
			diff := minHeight - height
			p.Velocity[2] += diff
			p.Camera.Position[2] += diff
		}

		p.Velocity = p.Velocity.Mul(0.95)

		if height > maxHeight && p.Velocity[2] > -0.0 {
			p.Velocity[2] -= 0.001
		}

		groundFactor := (height - groundHeight) / 20
		if groundFactor > 1 {
			groundFactor = 1
		}
		if groundFactor < 0 {
			groundFactor = 0
		}
		groundFactor = 1 - groundFactor
		groundFactor = groundFactor * groundFactor
	}
}
