package simulation

import (
	//	"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
)

func PortalPassed(portal *gamestate.Portal, pos1, pos2 mgl.Vec4f) bool {
	m := portal.View()
	pos1 = m.Mul4x1(pos1)
	pos2 = m.Mul4x1(pos2)
	n := portal.Normal

	a := n.Dot(pos1)
	b := n.Dot(pos2)

	if a*b < 0 {
		c := a / (a - b)
		pos3 := pos1.Mul(c).Add(pos2.Mul(1 - c))
		return -1 < pos3[0] && pos3[0] < 1 && -1 < pos3[1] && pos3[1] < 1
	} else {
		return false
	}
}

func Simulate(gs *gamestate.GameState) {
	player := gs.Player
	cam := &gs.Player.Camera
	oldPos := cam.Position
	UpdatePlayer(gs.Player, gs)
	newPos := cam.Position

	nearestPortal := gs.World.NearestPortal(oldPos)

	if PortalPassed(nearestPortal, oldPos, newPos) {
		// Enter Portal
		transform := nearestPortal.Transform()
		cam.SetModel(transform.Mul4(cam.Model()))
		player.Velocity = transform.Mul4x1(player.Velocity)
	}
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
	move_xyz := mgl.Vec3f{move[0], move[1], move[2]}
	move_xyz = p.Camera.Orientation.Rotate(move_xyz)
	move = mgl.Vec4f{move_xyz[0], move_xyz[1], move_xyz[2], 0}

	if !gs.Options.PlayerPhysics {
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
