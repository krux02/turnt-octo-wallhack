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

func Log2(uint64) uint64

func XYZ(v mgl.Vec4f) mgl.Vec3f {
	return mgl.Vec3f{v[0], v[1], v[2]}
}

func Vector(n mgl.Vec3f) mgl.Vec4f {
	return mgl.Vec4f{n[0], n[1], n[2], 0}
}

func Point(p mgl.Vec3f) mgl.Vec4f {
	return mgl.Vec4f{p[0], p[1], p[2], 1}
}

func TriangleIntersection(V1, V2, V3, O, D mgl.Vec3f) (out float32, hit bool) {
	const EPSILON = 0.000001
	var e1, e2 mgl.Vec3f //Edge1, Edge2
	var P, Q, T mgl.Vec3f
	var det, inv_det, u, v, t float32

	//Find vectors for two edges sharing V1
	e1 = V2.Sub(V1)
	e2 = V3.Sub(V1)
	//Begin calculating determinant - also used to calculate u parameter
	P = D.Cross(e2)
	//if determinant is near zero, ray lies in plane of triangle
	det = e1.Dot(P)
	//NOT CULLING
	if det > -EPSILON && det < EPSILON {
		return 0, false
	}
	inv_det = 1 / det

	//calculate distance from V1 to ray origin
	T = O.Sub(V1)

	//Calculate u parameter and test bound
	u = T.Dot(P) * inv_det
	//The intersection lies outside of the triangle
	if u < 0 || u > 1 {
		return 0, false
	}

	//Prepare to test v parameter
	Q = T.Cross(e1)

	//Calculate V parameter and test bound
	v = D.Dot(Q) * inv_det
	//The intersection lies outside of the triangle
	if v < 0 || u+v > 1 {
		return 0, false
	}

	t = e2.Dot(Q) * inv_det

	if t > EPSILON { //ray intersection
		return t, true
	}

	// No hit, no win
	return 0, false
}
