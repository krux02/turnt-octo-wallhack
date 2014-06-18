package simulation

import (
	//	"fmt"
	//mgl "github.com/krux02/mathgl/mgl32"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
	"github.com/krux02/turnt-octo-wallhack/helpers"
	"github.com/krux02/turnt-octo-wallhack/particles"
)

func Simulate(gs *gamestate.GameState, ps *particles.ParticleSystem) {
	player := gs.World.Player
	cam := &gs.World.Player.Camera
	oldPos := cam.Position
	UpdatePlayer(gs.World.Player, gs)
	newPos := cam.Position

	transform := gs.World.PortalTransform(oldPos, newPos)
	cam.SetModel(transform.Mul4(cam.Model()))
	player.Velocity = transform.Mul4x1(player.Velocity)

	if gs.Options.ParticlePhysics {
		ps.DoStep(gs)
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

	move_xyz := p.Camera.Orientation.Rotate(move.Vec3())
	move = move_xyz.Vec4(0)

	if !gs.Options.PlayerPhysics {
		move = move.Mul(0.1)
		p.Velocity = move
		p.Camera.MoveAbsolute(move)
	} else {
		move = move.Mul(0.01)
		p.Velocity = p.Velocity.Add(move)
		p.Camera.MoveAbsolute(p.Velocity)

		pos := p.Position()

		groundHeight := gs.World.HeightMap.Get2f(pos[0], pos[1])
		groundNormal := gs.World.HeightMap.Normal2f(pos[0], pos[1])

		height := p.Camera.Position[2]
		minHeight := groundHeight + 1.5
		maxHeight := groundHeight + 20

		if height < minHeight {
			diff := minHeight - height
			p.Velocity = p.Velocity.Add(helpers.Vector(groundNormal.Mul(diff)))
			//p.Camera.Position[2] += diff
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
