package helpers

import mgl "github.com/Jragonmiris/mathgl"

func HomogenXYZ(v mgl.Vec4f) mgl.Vec3f {
	v = Homogen(v)
	return mgl.Vec3f{v[0], v[1], v[2]}
}

func Homogen(v mgl.Vec4f) mgl.Vec4f {
	return v.Mul(1 / v[3])
}

func HomogenDist(p1, p2 mgl.Vec4f) float32 {
	return HomogenDiff(p1, p2).Len()
}

func HomogenDiff(p1, p2 mgl.Vec4f) mgl.Vec4f {
	return Homogen(p1).Sub(Homogen(p2))
}

func HomogenDiffXYZ(p1, p2 mgl.Vec4f) mgl.Vec3f {
	v := HomogenDiff(p1, p2)
	return mgl.Vec3f{v[0], v[1], v[2]}
}
