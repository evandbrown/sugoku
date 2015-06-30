package main

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

type Square struct {
	available map[int]bool
	val       int
	row       int
	col       int
}

func (s Square) String() string {
	return fmt.Sprintf("%v", strconv.Itoa(s.val))
}

func (s Square) key() string {
	return strconv.Itoa(s.row) + strconv.Itoa(s.col)
}

func (s Square) numAvailable() int {
	i := 0
	for _, a := range s.available {
		if a {
			i++
		}
	}
	return i
}

func (s *Square) duplicate() *Square {
	s2 := *s
	s2.available = make(map[int]bool)
	for k, v := range s.available {
		s2.available[k] = v
	}
	return &s2
}

type Board struct {
	squares [][]*Square
	slen    int // length of the square
	padding int // padding when printing
}

func (b Board) eliminate(s *Square, val int) ([]*Square, error) {
	solved := make([]*Square, 0)
	peers := b.peers(s)
	for k, _ := range peers {
		if peers[k].numAvailable() == 1 && peers[k].val == 0 && peers[k].val == val {
			b.uneliminate(s, val)
			return nil, errors.New(fmt.Sprintf("Can't eliminate %v from %v because it's the last value. All peer possibilities were restored", val, peers[k]))
		}
		if peers[k].numAvailable() == 2 {
			solved = append(solved, peers[k])
		}
		peers[k].available[val] = false
	}
	return solved, nil
}

func (b Board) uneliminate(s *Square, val int) {
	peers := b.peers(s)
	peers = append(peers, s)
	s.available[val] = true
	for k, _ := range peers {
		peers[k].available = b.availableVals(peers[k])
	}
}

func (b Board) set(s *Square, val int) ([]*Square, error) {
	if val == 0 {
		s.val = 0
		b.uneliminate(s, val)
		return nil, nil
	} else {
		solved, err := b.eliminate(s, val)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Cannot set %v at %v as it causes a problem on propagation (%v)", val, s, err))
		} else {
			s.val = val
			return solved, nil
		}
	}
}

func (b Board) duplicate() *Board {
	b2 := b
	nc := make([][]*Square, len(b.squares))
	for r := range nc {
		nc[r] = make([]*Square, len(b.squares[r]))
		copy(nc[r], b.squares[r])
		for c := range nc[r] {
			nc[r][c] = nc[r][c].duplicate()
		}
	}
	b2.squares = nc
	return &b2
}

func (b Board) flatten() []*Square {
	squares := make([]*Square, len(b.squares)*len(b.squares))
	for r := range b.squares {
		for c := range b.squares[r] {
			squares[r*len(b.squares)+c] = b.squares[r][c]
		}
	}
	return squares
}
func (b Board) shuffle(a []int) {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}

func (b Board) nextEasiestSquare() *Square {
	var next, low *Square
	for next = b.squares[0][0]; next != nil; next = b.nextSquare(next) {
		if low == nil && next.val == 0 { // initialize low
			low = next
		}
		if next.val == 0 && next.numAvailable() < low.numAvailable() { //Mark a new low
			low = next
		}
	}
	return low
}

func (b Board) nextEmptySquare() *Square {
	var s *Square
	for s = b.squares[0][0]; s != nil && s.val != 0; {
		s = b.nextSquare(s)
	}
	return s
}

func (b Board) peers(s *Square) []*Square {
	peers := make(map[string]*Square)
	all := append(b.squaresInRow(s), b.squaresInCol(s)...)
	all = append(all, b.squaresInGroup(s)...)
	for i := range all {
		peers[all[i].key()] = all[i]
	}

	delete(peers, s.key())
	all = make([]*Square, len(peers))
	i := 0
	for _, p := range peers {
		all[i] = p
		i++
	}
	return all
}

func (b Board) initPossible() {
	for _, s := range b.flatten() {
		s.available = b.availableVals(s)
	}
}
func (b Board) nextSquare(s *Square) *Square {
	if s.row+1 == b.slen*b.slen && s.col+1 == b.slen*b.slen {
		return nil
	} else if s.col+1 < b.slen*b.slen {
		return b.squares[s.row][s.col+1]
	} else {
		return b.squares[s.row+1][0]
	}
}

func (b Board) String() string {
	var board string

	for r := range b.squares {
		if r%b.slen == 0 {
			board += fmt.Sprintf("%v\n", strings.Repeat("-", b.slen*b.slen*3+(b.slen+1)*2-1))
		}
		for c := range b.squares[r] {
			if c%b.slen == 0 {
				board += fmt.Sprintf("| ")
			}
			board += fmt.Sprintf("%2s ", b.squares[r][c].String())
		}
		board += fmt.Sprintf("|\n")
	}
	board += fmt.Sprintf("%v\n", strings.Repeat("-", b.slen*b.slen*3+(b.slen+1)*2-1))
	return board
}

func (b Board) availableVals(s *Square) map[int]bool {
	possible := b.possibleVals()
	peers := b.peers(s)
	for _, p := range peers {
		if p.val != 0 {
			delete(possible, p.val)
		}
	}
	return possible
}

func (b Board) possibleVals() map[int]bool {
	// Map of all possible
	p := make(map[int]bool)
	for i := 1; i <= b.slen*b.slen; i++ {
		p[i] = true
	}
	return p
}

func (b Board) squaresInRow(s *Square) []*Square {
	vals := make([]*Square, b.slen*b.slen)
	for i := 0; i < len(vals); i++ {
		vals[i] = b.squares[s.row][i]
	}
	return vals
}

func (b Board) squaresInCol(s *Square) []*Square {
	vals := make([]*Square, b.slen*b.slen)
	for i := 0; i < len(vals); i++ {
		vals[i] = b.squares[i][s.col]
	}
	return vals
}

func (b Board) squaresInGroup(s *Square) []*Square {
	squares := make([]*Square, b.slen*b.slen)

	firstSquare := b.firstSquareInGroup(s)
	i := 0
	for r := firstSquare.row; r < firstSquare.row+b.slen; r++ {
		for c := firstSquare.col; c < firstSquare.col+b.slen; c++ {
			squares[i] = b.squares[r][c]
			i++
		}
	}
	return squares
}

func (b Board) firstSquareInGroup(s *Square) *Square {
	var rr, cc int
	if s.row%b.slen == 0 {
		rr = s.row
	} else {
		rr = s.row - s.row%b.slen
	}
	if s.col%b.slen == 0 {
		cc = s.col
	} else {
		cc = s.col - s.col%b.slen
	}
	return b.squares[rr][cc]
}

func ParseBoard(s string) *Board {
	length := math.Sqrt(float64(len(s)))
	sq := math.Sqrt(length)
	board := newBoard(int(sq))
	for i, c := range s {
		if c == '.' {
			c = '0'
		}
		board.squares[i/9][i%9] = newSquare(int(c-'0'), i/9, i%9)
	}
	board.initPossible()
	return board
}

func newSquare(val int, row int, col int) *Square {
	s := &Square{val: val, row: row, col: col}
	s.available = make(map[int]bool)
	return s
}

func newBoard(slen int) *Board {
	b := new(Board)
	if slen == 0 {
		slen = 3
	}
	b.slen = slen
	b.padding = 5
	b.squares = make([][]*Square, b.slen*b.slen)
	for r := range b.squares {
		b.squares[r] = make([]*Square, b.slen*b.slen)
		for c := range b.squares[r] {
			b.squares[r][c] = newSquare(0, r, c)
		}
	}
	b.initPossible()
	return b
}

func mapToArr(a map[int]bool) []int {
	v := make([]int, 0)
	i := 0
	for val, avail := range a {
		if avail {
			v = append(v, val)
			i++
		}
	}
	return v
}
