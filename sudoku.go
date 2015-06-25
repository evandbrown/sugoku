package main

import (
	"errors"
	"fmt"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	v1 := HardBoard()
	if err := InitialPropagation(v1); err != nil {
		fmt.Println("Invalid board")
		return
	}
	for _, s := range v1.Flatten() {
		if s.row < 3 && s.col < 3 {
			fmt.Println(s)
		}
	}
	fmt.Println(v1)
	if err := Solve(v1); err != nil {
		fmt.Println(err)
		return
	}
	//	fmt.Println(v1)
}

func InitialPropagation(b *Board) error {
	b.InitPossible()
	return nil
}

func Solve(b *Board) error {
	s := b.NextEmptySquare()
	if s.row == -1 {
		fmt.Println("Finished!!")
		return nil
	}
	for val, available := range s.available {
		if available {
			_, err := b.Set(s, val)
			if err != nil {
				return errors.New(fmt.Sprintf("There was an error trying to solve %v with %v: %v", s, val, err))
			} else {
				if err = Solve(b); err != nil {
					b.Set(s, 0)
					continue
				} else {
					return nil
				}
			}
		}
	}
	return errors.New(fmt.Sprintf("Failed to solve %v with all values %v", s, s.available))
}
