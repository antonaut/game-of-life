package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/gif"
	"os"
	//"time"
)

const (
	WIDTH          = 100
	HEIGHT         = 100
	NSTEPS         = 1000
	GIF_DELAY_TIME = 2
	CLI_SLEEP_TIME = 100 // Milliseconds
)

var (
	current  [][]bool
	prev     [][]bool
	filename string
)

func updateSquareAt(x, y int, grid [][]bool) {
	nNeighbours := 0

	if x != WIDTH {
		// Right
		if grid[y][x+1] {
			nNeighbours++
		}

		// Top Right
		if y != 0 {
			if grid[y-1][x+1] {
				nNeighbours++
			}
		}

		// Bottom Right
		if y != HEIGHT {
			if grid[y+1][x+1] {
				nNeighbours++
			}
		}
	}

	if y != HEIGHT {

		// Bottom
		if grid[y+1][x] {
			nNeighbours++
		}

		// Bottom left
		if x != 0 {
			if grid[y+1][x-1] {
				nNeighbours++
			}
		}

	}

	if y != 0 {

		// Top
		if grid[y-1][x] {
			nNeighbours++
		}

		// Top Left
		if x != 0 {
			if grid[y-1][x-1] {
				nNeighbours++
			}
		}
	}

	if x != 0 {
		// Left
		if grid[y][x-1] {
			nNeighbours++
		}
	}

	current[y][x] = (grid[y][x] && nNeighbours == 2) || nNeighbours == 3
}

func printGrid() {
	for y := 0; y < HEIGHT; y++ {
		for x := 0; x < WIDTH; x++ {
			if current[y][x] {
				fmt.Print("o")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}

}

// Glider gun as per
// http://en.wikipedia.org/wiki/Gun_(cellular_automaton)#mediaviewer/File:Game_of_life_glider_gun.svg
// Implemented column wise
//
// and 'acorn pattern' as per
// http://pmav.eu/stuff/javascript-game-of-life-v3.1.1/
func initGrid() {

	// left square
	current[6][2] = true
	current[7][2] = true
	current[6][3] = true
	current[7][3] = true

	// "C"
	current[6][12] = true
	current[7][12] = true
	current[8][12] = true

	current[5][13] = true
	current[9][13] = true

	current[4][14] = true
	current[10][14] = true

	current[4][15] = true
	current[10][15] = true

	// "->"
	current[7][16] = true

	current[5][17] = true
	current[9][17] = true

	current[6][18] = true
	current[7][18] = true
	current[8][18] = true

	current[7][19] = true

	// "<="
	current[4][22] = true
	current[5][22] = true
	current[6][22] = true

	current[4][23] = true
	current[5][23] = true
	current[6][23] = true

	current[3][24] = true
	current[7][24] = true

	current[2][26] = true
	current[3][26] = true
	current[7][26] = true
	current[8][26] = true

	// right square
	current[4][36] = true
	current[5][36] = true
	current[4][37] = true
	current[5][37] = true

	// Acorn
	current[50][50] = true

	current[48][51] = true
	current[50][51] = true

	current[49][53] = true

	current[50][54] = true
	current[50][55] = true
	current[50][56] = true
}

func main() {

	flag.StringVar(&filename, "f", "tmp.gif", "The name of the animated gif to be written.")

	imgs := make([]*image.Paletted, NSTEPS)
	delays := make([]int, NSTEPS)

	for i := 0; i < NSTEPS; i++ {
		delays[i] = GIF_DELAY_TIME
	}

	current = make([][]bool, HEIGHT+1)
	prev = make([][]bool, HEIGHT+1)

	for s := 0; s < HEIGHT+1; s++ {
		current[s] = make([]bool, WIDTH+1)
		prev[s] = make([]bool, WIDTH+1)
	}

	initGrid()

	for i := 0; i < NSTEPS; i++ {
		fmt.Printf("**** Animating gif (%d%%)- Step %d of %d ****\n", int(float64(i)/float64(NSTEPS)*100.0), i, NSTEPS)
		r := image.Rect(0, 0, WIDTH, HEIGHT)
		imgs[i] = image.NewPaletted(r, palette.Plan9)
		//printGrid()

		for y := 0; y < HEIGHT; y++ {
			for x := 0; x < WIDTH; x++ {
				prev[y][x] = current[y][x]
			}
		}

		for y := 0; y < HEIGHT; y++ {
			for x := 0; x < WIDTH; x++ {

				var c color.Gray

				if current[y][x] {
					c.Y = 0
				} else {
					c.Y = 255
				}

				imgs[i].Set(x, y, c)
				updateSquareAt(x, y, prev)
			}
		}

		//time.Sleep(CLI_SLEEP_TIME * time.Millisecond)
	}

	var gImg = gif.GIF{imgs, delays, -1}

	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	err = gif.EncodeAll(file, &gImg)
	if err != nil {
		panic(err)
	}

}
