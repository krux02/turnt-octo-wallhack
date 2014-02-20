package gamestate

import (
	"math"
	"sort"
)

const Dims = 3

type KdElement interface {
	Position(dim int) float32
}

func Dist(this, that KdElement) float32 {
	var sum float32
	for i := 0; i < Dims; i++ {
		v := this.Position(i) - that.Position(i)
		sum += v * v
	}
	return float32(math.Sqrt(float64(sum)))
}

func abs(a float32) float32 {
	if a < 0 {
		return -a
	} else {
		return a
	}
}

type KdTree []KdElement

type nodeSorter struct {
	elements []KdElement
	dim      int
}

func (s *nodeSorter) Len() int {
	return len(s.elements)
}

func (s *nodeSorter) Swap(i, j int) {
	s.elements[i], s.elements[j] = s.elements[j], s.elements[i]
}

func (s *nodeSorter) Less(i, j int) bool {
	return s.elements[i].Position(s.dim) < s.elements[j].Position(s.dim)
}

func sortByDim(data []KdElement, dim int) {
	sort.Sort(&nodeSorter{data, dim})
}

func NewTree(data KdTree, dim int) KdTree {
	sortByDim(data, dim)
	dim = (dim + 1) % 3
	l := len(data)
	if l >= 2 {
		NewTree(data.Left(), dim)
		NewTree(data.Right(), dim)
	}
	return KdTree(data)
}

func (this KdTree) Left() KdTree {
	return this[0 : len(this)/2]
}

func (this KdTree) Right() KdTree {
	return this[len(this)/2+1 : len(this)]
}

func (this KdTree) Node() KdElement {
	return this[len(this)/2]
}

func (this KdTree) IsLeaf() bool {
	return len(this) == 1
}

func (this KdTree) RegionQuery(pos KdElement, radius float32, target []KdElement) []KdElement {
	return this.regionQuery(pos, radius, target, 0)
}

func (this KdTree) regionQuery(pos KdElement, radius float32, target []KdElement, dim int) []KdElement {
	dist := pos.Position(dim) - this.Node().Position(dim)

	// dist := pos - node
	// pos - radius <= node <=>  radius <= dist
	// pos + radius >= node <=>  radius >= -dist

	dim = (dim + 1) % Dims
	if radius < abs(dist) {
		target = append(target, this.Node())
	}
	if !this.IsLeaf() {
		if radius <= dist {
			target = this.Left().regionQuery(pos, radius, target, dim)
		}
		if radius >= -dist {
			target = this.Right().regionQuery(pos, radius, target, dim)
		}
	}
	return target
}

func (this KdTree) NearestQuery(pos KdElement) KdElement {
	return this.nearestQuery(pos, this.Node(), 0)
}

func (this KdTree) nearestQuery(pos, currentBest KdElement, dim int) KdElement {

	if Dist(pos, this.Node()) < Dist(pos, currentBest) {
		currentBest = this.Node()
	}
	if this.IsLeaf() {
		return currentBest
	}

	dist := pos.Position(dim) - this.Node().Position(dim)

	if dist < 0 {
		currentBest = this.Left().nearestQuery(pos, currentBest, (dim+1)%Dims)
		// test the other branch if there might be a candidate
		if Dist(pos, currentBest) > -dist {
			currentBest = this.Right().nearestQuery(pos, currentBest, (dim+1)%Dims)
		}
	} else {
		currentBest = this.Right().nearestQuery(pos, currentBest, (dim+1)%Dims)
		// test the other branch if there might be a candidate
		if Dist(pos, currentBest) > dist {
			currentBest = this.Left().nearestQuery(pos, currentBest, (dim+1)%Dims)
		}
	}

	return currentBest
}
