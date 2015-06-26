package main

import (
	"errors"
	"fmt"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(1)
	//runtime.GOMAXPROCS(runtime.NumCPU())
	v1 := ParseBoard(BOARD_HARD_1)
	status := make(chan *Board)
	go printProgress(status)
	if err := solve(v1, status); err != nil {
		fmt.Println(err)
		return
	}
	<-status
	if !SHOW_STATUS {
		fmt.Println(v1)
	}
}

func printProgress(status chan *Board) {
	for b := range status {
		if SHOW_STATUS {
			if STATUS_OVERWRITE {
				fmt.Printf("\033[3;1H")
			}
			fmt.Println(b)
		}
	}
}

func solve(b *Board, status chan *Board) error {
	s := b.NextEasiestSquare()
	status <- b
	if s == nil {
		close(status)
		return nil
	}
	try := make([]int, 0)
	i := 0
	for val, ok := range s.available {
		if ok {
			try = append(try, val)
			i++
		}
	}
	for _, val := range try {
		_, err := b.Set(s, val)
		if err != nil {
			return errors.New(fmt.Sprintf("There was an error trying to solve %v with %v: %v", s, val, err))
		} else {
			if err = solve(b, status); err != nil {
				b.Set(s, 0)
				continue
			} else {
				return nil
			}
		}
	}
	return errors.New(fmt.Sprintf("Failed to solve %v with all values %v", s, try))
}

const (
	BOARD_EMPTY      = "................................................................................."
	BOARD_HARD_1     = ".15.....................8....6....1..3.2.....2.............8..2........6........."
	BOARD_HARD_2     = "....7..2.8.......6.1.2.5...9.54....8.........3....85.1...3.2.8.4.......9.7..6...."
	SHOW_STATUS      = false
	STATUS_OVERWRITE = true
)
