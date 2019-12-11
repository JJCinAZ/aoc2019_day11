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
)

func main() {
	data, err := ioutil.ReadFile("pgm.dat")
	if err != nil {
		panic(err)
	}
	part1(string(data))
	part2(string(data))
}

func part1(sourceCode string) {
	panels = make(map[vector]panel, 1024)
	ioChan := make(chan int)
	vm := intcode.NewVM(1, intcode.Compile(sourceCode), ioChan)
	go handleIO(ioChan)
	vm.Pgm.Debug(false)
	if err := vm.ExecPgm(); err != nil {
		panic(err)
	}
	<-ioChan
	close(ioChan)
	fmt.Println(len(panels))
}

func part2(sourceCode string) {
	curX, curY = 0, 0
	curDirection = 0
	panels = make(map[vector]panel, 1024)
	panels[vector{curX, curY}] = panel{0, 1}  // Start on white panel
	ioChan := make(chan int)
	vm := intcode.NewVM(1, intcode.Compile(sourceCode), ioChan)
	go handleIO(ioChan)
	vm.Pgm.Debug(false)
	if err := vm.ExecPgm(); err != nil {
		panic(err)
	}
	<-ioChan
	close(ioChan)
	printBitmap()
}

func printBitmap() {
	minX, maxX, minY, maxY := getExtents()
	w := maxX - minX + 1
	h := maxY - minY + 1
	pixels := make([]rune, w*h)
	for i := range pixels {
		pixels[i] = ' '
	}
	for v, p := range panels {
		if p.color == 1 {
			pixels[(v.y - minY) * w + (v.x - minX)] = '🀫'
		}
	}
	for i := 0; i < w*h; i += w {
		fmt.Println(string(pixels[i : i+w]))
	}
}

func getExtents() (minX, maxX, minY, maxY int) {
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
	return
}

func takePicture() int {
	v := vector{curX, curY}
	if p, exists := panels[v]; exists {
		return p.color
	}
	return 0
}

func handleIO(ioChan chan int) {
	for {
		ioChan <- takePicture()
		lastColor, ok := <-ioChan
		if !ok {
			return
		}
		turn := <-ioChan
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
}
