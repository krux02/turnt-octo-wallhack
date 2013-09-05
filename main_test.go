package main

import (
	"fmt"
	"testing"
)

func TestPalmGeneration(t *testing.T) {
	word := NewHeightMap(32, 32)
	pt := NewPalmTrees(word, 32)

	length := len(pt.positions)

	for i := 1; i < length; i++ {
		if pt.positions[pt.sortedX[i-1]][0] > pt.positions[pt.sortedX[i]][0] ||
			pt.positions[pt.sortedY[i-1]][1] > pt.positions[pt.sortedY[i]][1] ||
			pt.positions[pt.sortedXInv[i-1]][0] < pt.positions[pt.sortedXInv[i]][0] ||
			pt.positions[pt.sortedYInv[i-1]][1] < pt.positions[pt.sortedYInv[i]][1] {

			t.Log(fmt.Sprintln(pt.sortedX))
			t.Log(fmt.Sprintln(pt.sortedY))
			t.Log(fmt.Sprintln(pt.sortedXInv))
			t.Log(fmt.Sprintln(pt.sortedYInv))

			t.Fail()
		}
	}
}
