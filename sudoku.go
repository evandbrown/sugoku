package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Running!")
	b := NewBoard()
	b.Print()

	fmt.Printf("Possible vals at 0, 0: %v\n", b.AvailableVals(8, 8))
	x, y := b.TopLeftCoord(1, 1)
	fmt.Printf("Top-left for 1,1: %v,%v\n", x, y)
	x, y = b.TopLeftCoord(8, 8)
	fmt.Printf("Top-left for 8,8: %v,%v\n", x, y)
	fmt.Printf("Available for 1,1: %v\n", b.AvailableVals(0, 0))
}

type Board struct {
	cells   [][]int
	slen    int // length of the square
	padding int // padding when printing
}

func (b *Board) Print() {
	var c string
	for x := range b.cells {
		if x%b.slen == 0 {
			fmt.Println(strings.Repeat(" -", b.slen*b.slen+b.slen))
		}
		for y := range b.cells[x] {
			if y%b.slen == 0 {
				fmt.Printf("| ")
			}
			if b.cells[x][y] == 0 {
				c = "*"
			} else {
				c = string(b.cells[x][y])
			}
			fmt.Printf("%v ", c)
		}
		fmt.Printf("|")
		fmt.Println()
	}
	fmt.Println(strings.Repeat(" -", b.slen*b.slen+b.slen))
}

func (b *Board) AvailableVals(x int, y int) []int {
	used := append(b.valsInRow(y), b.valsInCol(x)...)
	used = append(used, b.valsInSquare(x, y)...)
	p := b.possibleVals()

	for i := range used {
		delete(p, used[i])
	}
	available := make([]int, 0, len(p))
	for k := range p {
		available = append(available, k)
	}
	return available
}

func (b *Board) possibleVals() map[int]bool {
	// Map of all possible
	p := make(map[int]bool)
	for i := 1; i <= b.slen*b.slen; i++ {
		p[i] = true
	}
	return p
}
func (b *Board) valsInRow(y int) []int {
	vals := make([]int, b.slen*b.slen)
	for i := 0; i < len(vals); i++ {
		vals[i] = b.cells[i][y]
	}
	return vals
}

func (b *Board) valsInCol(x int) []int {
	vals := make([]int, b.slen*b.slen)
	for i := 0; i < len(vals); i++ {
		vals[i] = b.cells[x][i]
	}
	return vals
}

func (b *Board) valsInSquare(x int, y int) []int {
	vals := make([]int, b.slen*b.slen*b.slen)

	xx, yy := b.TopLeftCoord(x, y)
	i := 0
	for x = xx; x < xx+b.slen; x++ {
		for y = yy; y < yy+b.slen; y++ {
			vals[i] = b.cells[x][y]
			i++
		}
	}
	return vals
}

func (b *Board) TopLeftCoord(x int, y int) (xx int, yy int) {
	if x%b.slen == 0 {
		xx = x
	} else {
		xx = x - x%b.slen
	}
	if y%b.slen == 0 {
		yy = y
	} else {
		yy = y - y%b.slen
	}
	return
}

func NewBoard() *Board {
	b := new(Board)
	b.slen = 3
	b.padding = 5
	b.cells = make([][]int, b.slen*b.slen)
	for i := range b.cells {
		b.cells[i] = make([]int, b.slen*b.slen)
	}
	return b
}
