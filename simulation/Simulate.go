package simulation

import (
	"fmt"
	mgl "github.com/Jragonmiris/mathgl"
	"github.com/krux02/turnt-octo-wallhack/gamestate"
)

var i int = 2

func EnterPortal(portal *gamestate.Portal, camera *gamestate.Camera) {

	//quat1 := portal.Orientation
	//quat2 := portal.Target.Orientation

	//camMat := camera.Orientation.Inverse().Mat4()
	//model1 := dir1.Mat4()
	//model2 := dir2.Mat4()

	//camMat = (model2.Mul4(model1.Inv()).Mul4(camMat))

	//camera.SetFromMat4(camMat)
	//camera.Orientation = camera.Orientation.Mul(dir2).Mul(dir1.Inverse())
	//camera.Orientation = dir1.Inverse().Mul(dir2).Mul(camera.Orientation)
	//camera.Orientation = quat2.Mul(quat1.Inverse()).Mul(camera.Orientation)
	//camera.Orientation = camera.Orientation.Inverse().Mul(dir2.Inverse()).Mul(dir1).Inverse()
	//camera.Orientation = camera.Orientation.Inverse().Mul(quat1.Inverse()).Mul(quat2).Inverse()

	pos1 := portal.Position
	pos2 := portal.Target.Position

	camera.Position = camera.Position.Sub(pos1).Add(pos2)

	camO := camera.Orientation
	targO := portal.Target.Orientation

	switch i {
	case 0:
		camera.Orientation = camO.Mul(targO)
	case 1:
		camera.Orientation = camO.Mul(targO.Inverse())
	case 2:
		camera.Orientation = targO.Mul(camO)
	case 3:
		camera.Orientation = targO.Inverse().Mul(camO)
	case 4:
		camera.Orientation = targO.Inverse().Mul(camO).Mul(targO)
	case 5:
		camera.Orientation = targO.Mul(camO).Mul(targO).Inverse()
	}

	fmt.Println("enter portal i:", i)

	//i = (i + 1) % 6
}

func PortalPassed(portal *gamestate.Portal, pos1, pos2 mgl.Vec3f) bool {
	pos1w := mgl.Vec4f{pos1[0], pos1[1], pos1[2], 1}
	pos2w := mgl.Vec4f{pos2[0], pos2[1], pos2[2], 1}
	plane := portal.ClippingPlane(true)
	return pos1w.Dot(plane)*pos2w.Dot(plane) < 0
}

func Simulate(gs *gamestate.GameState) {
	cam := &gs.Player.Camera
	oldPos := cam.Position
	UpdatePlayer(gs.Player, gs)
	newPos := cam.Position

	nearestPortal := gs.World.NearestPortal(oldPos)

	if PortalPassed(nearestPortal, oldPos, newPos) {
		//fmt.Println("before", oldPos, newPos, cam)
		EnterPortal(nearestPortal, cam)
		//fmt.Println("after", cam.Position, cam)
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
	move = p.Camera.Orientation.Rotate(move)

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
