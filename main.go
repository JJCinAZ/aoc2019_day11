package main

import (
	"cloud.google.com/aoc2019/day11/intcode"
	"fmt"
	"io/ioutil"
)

type vector struct {
	x, y int
}

type panel struct {
	paintCount int
	color      int
}

var (
	panels                   map[vector]panel
	curX, curY, curDirection int
	lastColor                int
)

func main() {
	part2(67, 20)
}

func part1() {
	var (
		vm *intcode.VM
	)
	data, err := ioutil.ReadFile("pgm.dat")
	if err != nil {
		panic(err)
	}
	panels = make(map[vector]panel, 1024)
	lastColor = -1
	pgm := intcode.Compile(string(data))
	vm = intcode.NewVM(1, pgm, handleInput, handleOutput)
	vm.Pgm.Debug(false)
	if err := vm.ExecPgm(); err != nil {
		panic(err)
	}
	w, h := getDimensions()
	fmt.Println(len(panels), w, h)
}

func part2(w, h int) {
	var (
		vm *intcode.VM
	)
	data, err := ioutil.ReadFile("pgm.dat")
	if err != nil {
		panic(err)
	}
	curX, curY = 0, h>>1
	panels = make(map[vector]panel, 1024)
	// Start on white panel in center
	panels[vector{curX, curY}] = panel{0, 1}
	lastColor = -1
	pgm := intcode.Compile(string(data))
	vm = intcode.NewVM(1, pgm, handleInput, handleOutput)
	vm.Pgm.Debug(false)
	if err := vm.ExecPgm(); err != nil {
		panic(err)
	}
	pixels := make([]byte, w*h)
	for i := range pixels {
		pixels[i] = '.'
	}
	for v, p := range panels {
		if p.color == 1 {
			pixels[v.y*w+v.x] = '#'
		}
	}
	for i := 0; i < w*h; i += w {
		fmt.Println(string(pixels[i : i+w]))
	}
}

func getDimensions() (int, int) {
	var minX, minY, maxX, maxY int
	for v, _ := range panels {
		if v.x < minX {
			minX = v.x
		} else if v.x > maxX {
			maxX = v.x
		}
		if v.y < minY {
			minY = v.y
		} else if v.y > maxY {
			maxY = v.y
		}
	}
	return maxX - minX + 1, maxY - minY + 1
}

func handleInput() int {
	v := vector{curX, curY}
	if p, exists := panels[v]; exists {
		return p.color
	}
	return 0
}

func handleOutput(b int) {
	if lastColor == -1 {
		lastColor = b
		return
	}
	turn := b
	v := vector{curX, curY}
	if p, exists := panels[v]; exists {
		p.paintCount++
		p.color = lastColor
		panels[v] = p
	} else {
		panels[v] = panel{1, lastColor}
	}
	lastColor = -1
	switch turn {
	case 0:
		if curDirection == 0 {
			curDirection = 270
		} else {
			curDirection -= 90
		}
	case 1:
		if curDirection == 270 {
			curDirection = 0
		} else {
			curDirection += 90
		}
	}
	switch curDirection {
	case 0:
		curY--
	case 90:
		curX++
	case 180:
		curY++
	case 270:
		curX--
	}
}
