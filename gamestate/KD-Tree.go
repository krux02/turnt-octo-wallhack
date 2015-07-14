package gamestate

import (
	"math"
	"sort"
)

const Dims = 3

type KdElement interface {
	Dimension(dim int) float32
}

func Dist(this, that KdElement) float32 {
	var sum float32
	for i := 0; i < Dims; i++ {
		v := this.Dimension(i) - that.Dimension(i)
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
	return s.elements[i].Dimension(s.dim) < s.elements[j].Dimension(s.dim)
}

func sortByDim(data []KdElement, dim int) {
	sort.Sort(&nodeSorter{data, dim})
}

func NewTree(data KdTree) KdTree {
	return newTree(data, 0)
}

func newTree(data KdTree, dim int) KdTree {
	sortByDim(data, dim)
	dim = (dim + 1) % 3
	l := len(data)
	if l >= 2 {
		newTree(data.Left(), dim)
		newTree(data.Right(), dim)
	}
	return KdTree(data)
}

func (this KdTree) Left() KdTree {
	return this[0 : len(this)/2]
}

func (this KdTree) Right() KdTree {
	return this[len(this)/2+1:]
}

func (this KdTree) HasRight() bool {
	return len(this) > 2
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
	dist := pos.Dimension(dim) - this.Node().Dimension(dim)

	// dist := pos - node
	// pos - radius <= node <=>  radius <= dist
	// pos + radius >= node <=>  radius >= -dist

	dim = (dim + 1) % Dims
	if radius < abs(dist) {
		node := this.Node()
		target = append(target, node)
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

func (this KdTree) NearestQuery(pos KdElement, filter func(KdElement) bool) KdElement {
	return this.nearestQuery(pos, nil, filter, 0)
}

func (this KdTree) nearestQuery(pos, currentBest KdElement, filter func(KdElement) bool, dim int) KdElement {
	if currentBest == nil || Dist(pos, this.Node()) < Dist(pos, currentBest) {
		if filter(this.Node()) {
			currentBest = this.Node()
		}
	}
	if this.IsLeaf() {
		return currentBest
	}

	dist := pos.Dimension(dim) - this.Node().Dimension(dim)

	if dist < 0 {
		currentBest = this.Left().nearestQuery(pos, currentBest, filter, (dim+1)%Dims)
		// test the other branch if there might be a candidate
		if this.HasRight() {
			if currentBest == nil || Dist(pos, currentBest) > -dist {
				currentBest = this.Right().nearestQuery(pos, currentBest, filter, (dim+1)%Dims)
			}
		}
	} else {
		if this.HasRight() {
			currentBest = this.Right().nearestQuery(pos, currentBest, filter, (dim+1)%Dims)
		}
		// test the other branch if there might be a candidate
		if currentBest == nil || Dist(pos, currentBest) > dist {
			currentBest = this.Left().nearestQuery(pos, currentBest, filter, (dim+1)%Dims)
		}
	}

	return currentBest
}
