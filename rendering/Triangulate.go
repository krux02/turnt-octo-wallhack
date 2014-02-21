package rendering

func TriangulationIndices(w, h int) []int32 {
	indexCount := 6 * w * h
	indices := make([]int32, indexCount, indexCount)

	i := 0
	put := func(v int) {
		indices[i] = int32(v)
		i += 1
	}

	flat := func(x, y int) int {
		return (w+1)*y + x
	}

	quad := func(x, y int) {
		v1 := flat(x, y)
		v2 := flat(x+1, y)
		v3 := flat(x, y+1)
		v4 := flat(x+1, y+1)

		put(v1)
		put(v2)
		put(v3)

		put(v3)
		put(v2)
		put(v4)
	}

	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			quad(i, j)
		}
	}

	return indices
}
