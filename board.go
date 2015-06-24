package main

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"
)

type Board struct {
	squares [][]*Square
	slen    int // length of the square
	padding int // padding when printing
}

type Square struct {
	possible []int
	locked   bool
	val      int
	row      int
	col      int
}

func (s *Square) String() string {
	return fmt.Sprintf("(%v,%v) = %v %v", s.row+1, s.col+1, s.val, s.possible)
}

func (s *Square) Key() string {
	return strconv.Itoa(s.row) + strconv.Itoa(s.col)
}

func (s *Square) AvailString() string {
	return fmt.Sprintf("%v %v (%v,%v)", s.val, s.possible, s.row+1, s.col+1)
}
func (s *Square) ValString() string {
	return fmt.Sprintf("%v", strconv.Itoa(s.val))
}
func (b *Board) Peers(s *Square) []*Square {
	peers := make(map[string]*Square)
	all := append(b.squaresInRow(s), b.squaresInCol(s)...)
	all = append(all, b.squaresInGroup(s)...)
	for i := range all {
		peers[all[i].Key()] = all[i]
	}

	delete(peers, s.Key())
	all = make([]*Square, len(peers))
	i := 0
	for _, p := range peers {
		all[i] = p
		i++
	}
	return all
}

func (b *Board) Propagate(s *Square, setPeers bool) error {
	peers := b.Peers(s)
	for i := range peers {
		peers[i].possible = b.availableVals(peers[i])
		if len(peers[i].possible) == 0 && peers[i].val == 0 {
			return errors.New(fmt.Sprintf("Cannot propagate %v as it leaves no options for %v\n", s, peers[i]))
		}
		if len(peers[i].possible) == 1 && peers[i].val == 0 && setPeers {
			return b.Set(peers[i], peers[i].possible[0])
		}
	}
	return nil
}

func (b *Board) Set(s *Square, val int) error {
	s.val = val
	if err := b.Propagate(s, true); err != nil {
		s.val = 0
		b.Propagate(s, false)
		return errors.New(fmt.Sprintf("Cannot set %v at %v as it causes a problem on propagation (%v)", val, s, err))
	}
	//fmt.Printf("\033[3;1H")
	return nil
}

func (b *Board) Duplicate() *Board {
	b2 := *b
	nc := make([][]*Square, len(b.squares))
	for i := range nc {
		nc[i] = make([]*Square, len(b.squares[i]))
		copy(nc[i], b.squares[i])
	}
	b2.squares = nc
	return &b2
}

func (b *Board) Flatten() []*Square {
	squares := make([]*Square, len(b.squares)*len(b.squares))
	for r := range b.squares {
		for c := range b.squares[r] {
			squares[r*len(b.squares)+c] = b.squares[r][c]
		}
	}
	return squares
}

func (b *Board) NextEmptySquare() *Square {
	var s *Square
	for s = b.squares[0][0]; s.row != -1 && s.val != 0; {
		s = b.NextSquare(s)
	}
	return s
}

func (b *Board) NextSquare(s *Square) *Square {
	// -1 if last cell
	if s.row+1 == b.slen*b.slen && s.col+1 == b.slen*b.slen {
		return NewSquare(0, -1, -1)
	} else if s.col+1 < b.slen*b.slen {
		return b.squares[s.row][s.col+1]
	} else {
		return b.squares[s.row+1][0]
	}
}

func (b *Board) ClearForward(s *Square) {
	sq := b.NextSquare(s)
	if sq.row == -1 {
		return
	}
	if !sq.locked {
		sq.val = 0
	}
	b.ClearForward(sq)
}

func (b *Board) String() string {
	var board string

	for r := range b.squares {
		if r%b.slen == 0 {
			board += fmt.Sprintf("%v\n", strings.Repeat("-", b.slen*b.slen*3+(b.slen+1)*2-1))
		}
		for c := range b.squares[r] {
			if c%b.slen == 0 {
				board += fmt.Sprintf("| ")
			}
			board += fmt.Sprintf("%2s ", b.squares[r][c].ValString())
		}
		board += fmt.Sprintf("|\n")
	}
	board += fmt.Sprintf("%v\n", strings.Repeat("-", b.slen*b.slen*3+(b.slen+1)*2-1))
	return board
}

func (b *Board) Shuffle(a []int) {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}

func (b *Board) availableVals(s *Square) []int {
	p := b.possibleVals()
	peers := b.Peers(s)
	for i := range peers {
		delete(p, peers[i].val)
	}
	available := make([]int, 0, len(p))
	for k := range p {
		available = append(available, k)
	}
	sort.Ints(available)
	//b.Shuffle(available)
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

func (b *Board) squaresInRow(s *Square) []*Square {
	vals := make([]*Square, b.slen*b.slen)
	for i := 0; i < len(vals); i++ {
		vals[i] = b.squares[s.row][i]
	}
	return vals
}

func (b *Board) squaresInCol(s *Square) []*Square {
	vals := make([]*Square, b.slen*b.slen)
	for i := 0; i < len(vals); i++ {
		vals[i] = b.squares[i][s.col]
	}
	return vals
}

func (b *Board) squaresInGroup(s *Square) []*Square {
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

func (b *Board) firstSquareInGroup(s *Square) *Square {
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

func BoardFromS(s string) *Board {
	length := math.Sqrt(float64(len(s)))
	sq := math.Sqrt(length)
	board := NewBoard(int(sq))
	for i, c := range s {
		if c == '.' {
			c = '0'
		}
		board.squares[i/9][i%9] = NewSquare(int(c-'0'), i/9, i%9)
	}
	return board
}

func HardBoard() *Board {
	b := BoardFromS(".15.....................8....6....1..3.2.....2.............8..2........6.........")
	return b
}

func EmptyBoard() *Board {
	b := BoardFromS(".................................................................................")
	return b
}

func NewBoard(slen int) *Board {
	b := new(Board)
	if slen == 0 {
		slen = 3
	}
	b.slen = slen
	b.padding = 5
	b.squares = make([][]*Square, b.slen*b.slen)
	for r := range b.squares {
		b.squares[r] = make([]*Square, b.slen*b.slen)
	}
	return b
}

func NewSquare(val int, row int, col int) *Square {
	s := &Square{val: val, row: row, col: col}
	if s.val != 0 {
		s.locked = true
	}
	return s
}
