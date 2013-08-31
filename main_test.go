package main

import (
	"fmt"
	"testing"
)

func TestDiamondSquare(t *testing.T) {

	const w, h = 65, 65
	heights := NewHeightMap(w, h)

	DiamondSquare(heights, 65)

	if testing.Short() {

		t.Log("short")
	}

	if testing.Verbose() {
		t.Log("bla bla bla")
		fmt.Println("bla bla bla")
	}

	fmt.Println(heights)
}
