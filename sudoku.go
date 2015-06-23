package main

import (
	"fmt"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	v1 := HardBoard()
	done := make(chan *Board)
	defer close(done)
	r, c := v1.NextEmptyCell(0, 0)
	go Solve(v1, done, r, c)
	go Solve(v1.Duplicate(), done, r, c)
	go Solve(v1.Duplicate(), done, r, c)
	go Solve(v1.Duplicate(), done, r, c)
	fmt.Println("Waiting for solution...")
	solution := <-done
	fmt.Println(solution.Print())
}

func Solve(b *Board, done chan *Board, r int, c int) bool {
	// Get the available values for this cell
	//fmt.Printf("\033[3;1H")
	//time.Sleep(0 * time.Millisecond)
	a := b.AvailableVals(r, c)
	for v := range a {
		b.cells[r][c] = a[v]
		//fmt.Printf(b.Print())
		// Get next cell
		if nr, nc := b.NextEmptyCell(r, c); nr != -1 {
			// If an error is returned, try the next value by continuing in the loop
			if err := Solve(b, done, nr, nc); err == true {
				continue
			} else {
				return false
			}
		} else {
			// No more cells. We're done
			done <- b
			return false
		}
	}
	b.cells[r][c] = 0
	return true
}
