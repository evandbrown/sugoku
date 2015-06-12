package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	b := NewBoard(2)
	Solve(b, 0, 0)
	b.Print()

	b = NewBoard(3)
	Solve(b, 0, 0)
	b.Print()

	b = NewBoard(4)
	Solve(b, 0, 0)
	b.Print()

	b = NewBoard(5)
	Solve(b, 0, 0)
	b.Print()
}

func NewBoard(slen int) *Board {
	b := new(Board)
	if slen == 0 {
		slen = 3
	}
	b.slen = slen
	b.padding = 5
	b.cells = make([][]int, b.slen*b.slen)
	for i := range b.cells {
		b.cells[i] = make([]int, b.slen*b.slen)
	}
	return b
}

func Solve(b *Board, x int, y int) (err error) {
	// Get the available values for this cell
	a := b.AvailableVals(x, y)
	for v := range a {
		b.cells[x][y] = a[v]
		if nx, ny := b.NextEmptyCell(x, y); nx != -1 {
			// If an error is returned, continue in the for loop
			if err := Solve(b, nx, ny); err != nil {
				//fmt.Println(err)
				continue
			} else {
				return nil
			}
		} else {
			// No more cells. We're done
			return nil
		}
	}
	b.cells[x][y] = 0
	return errors.New(fmt.Sprintf("No solution at %v, %v. Setting cell to 0", x, y))
}

type Board struct {
	cells   [][]int
	slen    int // length of the square
	padding int // padding when printing
}

func (b *Board) NextEmptyCell(x int, y int) (int, int) {
	// -1 if last cell
	if x+1 == b.slen*b.slen && y+1 == b.slen*b.slen {
		return -1, -1
	} else if x+1 < b.slen*b.slen {
		return x + 1, y
	} else {
		return 0, y + 1
	}
}

func (b *Board) Print() {
	fmt.Printf("\n\n\n")
	var c string
	for x := range b.cells {
		if x%b.slen == 0 {
			fmt.Println(strings.Repeat("-", b.slen*b.slen*3+(b.slen+1)*2-1))
		}
		for y := range b.cells[x] {
			if y%b.slen == 0 {
				fmt.Printf("| ")
			}
			if b.cells[x][y] == 0 {
				c = "*"
			} else {
				c = strconv.Itoa(b.cells[x][y])
			}
			fmt.Printf("%2s ", c)
		}
		fmt.Printf("|")
		fmt.Println()
	}
	fmt.Println(strings.Repeat("-", b.slen*b.slen*3+(b.slen+1)*2-1))
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
