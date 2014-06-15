package gamestate

import (
	mgl "github.com/krux02/mathgl/mgl32"
	"math"
	"math/rand"
	"testing"
)

func TestCameraMatrix(t *testing.T) {
	var camera1, camera2 Camera
	camera1.Orientation.W = 1
	camera1.Yaw(rand.Float32() * 2 * math.Pi)
	camera1.Pitch(rand.Float32() * 2 * math.Pi)
	camera1.Roll(rand.Float32() * 2 * math.Pi)
	camera1.Position = mgl.Vec4{rand.Float32() * 10, rand.Float32() * 10, rand.Float32() * 10, 1}
	t.Log(camera1)

	view1 := camera1.View()
	camera2.SetView(view1)
	view2 := camera2.View()

	err := SquareErrorMat4(view1.Sub(view2))

	if err > 1e-10 {
		t.Error(err, view1, view2)
	}
}

func SquareErrorMat4(m mgl.Mat4) (err float32) {
	for i := 0; i < 16; i++ {
		err += m[i] * m[i]
	}
	return
}
