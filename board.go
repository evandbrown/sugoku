package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

type Board struct {
	cells   [][]int
	slen    int // length of the square
	padding int // padding when printing
}

func (b *Board) Duplicate() *Board {
	b2 := *b
	nc := make([][]int, len(b.cells))
	for i := range nc {
		nc[i] = make([]int, len(b.cells[i]))
		copy(nc[i], b.cells[i])
	}
	b2.cells = nc
	return &b2
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

func BoardFromS(s string) *Board {
	length := math.Sqrt(float64(len(s)))
	sq := math.Sqrt(length)
	board := NewBoard(int(sq))
	for i, c := range s {
		if c == '.' {
			c = '0'
		}
		board.cells[i/9][i%9] = int(c - '0')
	}
	return board
}

func HardBoard() *Board {
	b := BoardFromS(".15.....................8....6....1..3.2.....2.............8..2........6.........")
	return b
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
