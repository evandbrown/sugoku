package main

import (
	"errors"
	"fmt"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	v1 := HardBoard()
	PropagateAvailable(v1, v1.Flatten())
	Solve(v1)
	fmt.Println("Waiting for solution...")
	fmt.Println(v1)
}

// Propagate available values to provided squares
func PropagateAvailable(b *Board, squares []*Square) error {
	for i := range squares {
		s := squares[i]
		s.possible = b.availableVals(s)

		if len(s.possible) == 0 {
			return errors.New(fmt.Sprintf("Error: No solutions for %v", s))
		}

		if len(s.possible) == 1 && s.val == 0 {
			s.val = s.possible[0]
			fmt.Printf("\033[3;1H")
			fmt.Println(b)
			if err := PropagateAvailable(b, b.Peers(s)); err != nil {
				s.val = 0
				PropagateAvailable(b, b.Peers(s))
			}
		}
	}
	return nil
}

func Solve(b *Board) bool {
	s := b.NextEmptySquare()
	if s.row == -1 {
		fmt.Println("Finished!!")
		return false
	}
	fmt.Printf("\033[3;1H")
	for v := 0; v < len(s.possible); v++ {
		s.val = s.possible[v]
		fmt.Println(b)
		if err := PropagateAvailable(b, b.Peers(s)); err != nil {
			s.val = 0
			PropagateAvailable(b, b.Peers(s))
			return true
		}
		// Propagation may have solved the puzzle. Check
		if b.NextEmptySquare().row == -1 {
			return false
		}
		// If an error is returned, try the next value by continuing in the loop
		if err := Solve(b); err == true {
			continue
		} else {
			return false
		}
	}
	return false
}
