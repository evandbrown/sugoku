package main

import (
	"errors"
	"fmt"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	done := make(chan *Board)
	defer close(done)
	boards := make([]*Board, 2)
	for i := range boards {
		boards[i] = ParseBoard(BOARD_HARD_1)
		go solve(boards[i], done)
	}
	printProgress(boards, done)
}

func printProgress(boards []*Board, stop chan *Board) {
	for {
		select {
		case fin := <-stop:
			if SHOW_STATUS {
				if STATUS_OVERWRITE {
					fmt.Printf("\033[3;1H")
				}
				for _, b := range boards {
					fmt.Println(".")
					fmt.Println(b)
				}
			} else {
				fmt.Println(fin)
			}

			return
		default:
			if SHOW_STATUS {
				if STATUS_OVERWRITE {
					fmt.Printf("\033[3;1H")
				}
				for _, b := range boards {
					fmt.Println(".")
					fmt.Println(b)
				}
			}
		}
	}
}

func solve(b *Board, done chan *Board) error {
	time.Sleep(10 * time.Millisecond)
	s := b.nextEasiestSquare()
	if s == nil {
		done <- b
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
	b.shuffle(try)
	for _, val := range try {
		_, err := b.set(s, val)
		if err != nil {
			return errors.New(fmt.Sprintf("There was an error trying to solve %v with %v: %v", s, val, err))
		} else {
			if err = solve(b, done); err != nil {
				b.set(s, 0)
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
	SHOW_STATUS      = true
	STATUS_OVERWRITE = true
)
