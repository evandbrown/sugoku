package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func main() {
	b := NewBoard(3)
	b.cells[0] = []int{3, 4, 0, 0, 0, 0, 9, 0, 0}
	b.cells[1] = []int{9, 0, 0, 8, 4, 0, 0, 0, 0}
	b.cells[2] = []int{0, 0, 8, 0, 0, 2, 0, 0, 5}
	b.cells[3] = []int{0, 2, 4, 0, 0, 0, 0, 1, 0}
	b.cells[4] = []int{0, 0, 6, 4, 0, 7, 8, 0, 0}
	b.cells[5] = []int{0, 3, 0, 0, 0, 0, 7, 5, 0}
	b.cells[6] = []int{2, 0, 0, 5, 0, 0, 4, 0, 0}
	b.cells[7] = []int{0, 0, 0, 0, 2, 6, 0, 0, 7}
	b.cells[8] = []int{0, 0, 5, 0, 0, 0, 0, 2, 3}
	b = NewBoard(4)
	r, c := b.NextEmptyCell(0, 0)
	Solve(b, r, c)
	fmt.Println()
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

func Solve(b *Board, r int, c int) error {
	// Get the available values for this cell
	fmt.Printf("\033[3;1H")
	time.Sleep(1 * time.Millisecond)
	a := b.AvailableVals(r, c)
	for v := range a {
		b.cells[r][c] = a[v]
		fmt.Printf(b.Print())
		if nr, nc := b.NextEmptyCell(r, c); nr != -1 {
			// If an error is returned, continue in the for loop
			if err := Solve(b, nr, nc); err != nil {
				continue
			} else {
				return nil
			}
		} else {
			// No more cells. We're done
			return nil
		}
	}
	b.cells[r][c] = 0
	return errors.New(fmt.Sprintf("No solution at %v, %v. Setting cell to 0", r, c))
}

type Board struct {
	cells   [][]int
	slen    int // length of the square
	padding int // padding when printing
}

func (b *Board) NextEmptyCell(r int, c int) (int, int) {
	if r == 0 && c == 0 && b.cells[r][c] == 0 {
		return r, c
	}
	for r, c = b.NextCell(r, c); r != -1 && b.cells[r][c] != 0; {
		r, c = b.NextCell(r, c)
	}
	return r, c
}
func (b *Board) NextCell(r int, c int) (int, int) {
	// -1 if last cell
	if r+1 == b.slen*b.slen && c+1 == b.slen*b.slen {
		return -1, -1
	} else if c+1 < b.slen*b.slen {
		return r, c + 1
	} else {
		return r + 1, 0
	}
}

func (b *Board) Print() string {
	var board string

	var v string
	for r := range b.cells {
		if r%b.slen == 0 {
			board += fmt.Sprintf("%v\n", strings.Repeat("-", b.slen*b.slen*3+(b.slen+1)*2-1))
		}
		for c := range b.cells[r] {
			if c%b.slen == 0 {
				board += fmt.Sprintf("| ")
			}
			if b.cells[r][c] == 0 {
				v = "*"
			} else {
				v = strconv.Itoa(b.cells[r][c])
			}
			board += fmt.Sprintf("%2s ", v)
		}
		board += fmt.Sprintf("|\n")
	}
	board += fmt.Sprintf("%v", strings.Repeat("-", b.slen*b.slen*3+(b.slen+1)*2-1))
	return board
}

func (b *Board) Shuffle(a []int) {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}

func (b *Board) AvailableVals(r int, c int) []int {
	used := append(b.valsInRow(r), b.valsInCol(c)...)
	used = append(used, b.valsInSquare(r, c)...)
	p := b.possibleVals()

	for i := range used {
		delete(p, used[i])
	}
	available := make([]int, 0, len(p))
	for k := range p {
		available = append(available, k)
	}
	//sort.Ints(available)
	b.Shuffle(available)
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

func (b *Board) valsInRow(r int) []int {
	vals := make([]int, b.slen*b.slen)
	for i := 0; i < len(vals); i++ {
		vals[i] = b.cells[r][i]
	}
	return vals
}

func (b *Board) valsInCol(c int) []int {
	vals := make([]int, b.slen*b.slen)
	for i := 0; i < len(vals); i++ {
		vals[i] = b.cells[i][c]
	}
	return vals
}

func (b *Board) valsInSquare(r int, c int) []int {
	vals := make([]int, b.slen*b.slen*b.slen)

	rr, cc := b.TopLeftCoord(r, c)
	i := 0
	for r = rr; r < rr+b.slen; r++ {
		for c = cc; c < cc+b.slen; c++ {
			vals[i] = b.cells[r][c]
			i++
		}
	}
	return vals
}

func (b *Board) TopLeftCoord(r int, c int) (int, int) {
	var rr, cc int
	if r%b.slen == 0 {
		rr = r
	} else {
		rr = r - r%b.slen
	}
	if c%b.slen == 0 {
		cc = c
	} else {
		cc = c - c%b.slen
	}
	return rr, cc
}
