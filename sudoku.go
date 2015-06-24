package main

import (
	"errors"
	"fmt"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	v1 := HardBoard()
	//v1 = EmptyBoard()
	if err := InitialPropagation(v1); err != nil {
		fmt.Println("Invalid board")
		return
	}
	if err := Solve(v1); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(v1)
}

func InitialPropagation(b *Board) error {
	for r := range b.squares {
		for c := range b.squares[r] {
			if b.squares[r][c].val != 0 {
				if err := b.Propagate(b.squares[r][c], true); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func Solve(b *Board) error {
	s := b.NextEmptySquare()
	if s.row == -1 {
		fmt.Println("Finished!!")
		return nil
	}
	for _, p := range s.possible {
		err := b.Set(s, p)
		if err != nil {
			return errors.New(fmt.Sprintf("There was an error trying to solve %v with %v: %v", s, p, err))
		} else {
			if err = Solve(b); err != nil {
				b.ClearForward(s)
				continue
			} else {
				return nil
			}
		}
	}
	return errors.New(fmt.Sprintf("Failed to solve %v with all values %v", s, s.possible))
}
