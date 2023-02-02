// The MIT License (MIT)
//
// Copyright (c) 2015-2016 Martin Lindhe
// Copyright (c) 2016      Hajime Hoshi
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
// THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

//go:build example
// +build example

package main

import (
	"log"
  _ "image/jpeg"
  _ "image/png"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
  "github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
  img1 *ebiten.Image
  img2 *ebiten.Image
  img3 *ebiten.Image
  words *ebiten.Image
  danWidth int
  danHeight int
  nImg int = 3
)

func init() {
	rand.Seed(time.Now().UnixNano())

	var err error
	img1, _, err = ebitenutil.NewImageFromFile("dan.jpg")
	if err != nil {
		log.Fatal(err)
	}
  danWidth, danHeight = img1.Size()
  log.Print(danWidth)
  log.Print(danHeight)

	img2, _, err = ebitenutil.NewImageFromFile("dan2.jpg")
	if err != nil {
		log.Fatal(err)
	}

	img3, _, err = ebitenutil.NewImageFromFile("dan3.png")
	if err != nil {
		log.Fatal(err)
	}

	words, _, err = ebitenutil.NewImageFromFile("wordart.png")
	if err != nil {
		log.Fatal(err)
	}
}

// World represents the game state.
type World struct {
	area   []int
	width  int
	height int
}

// NewWorld creates a new world.
func NewWorld(width, height int, maxInitLiveCells int) *World {
	w := &World{
		area:   make([]int, width*height),
		width:  width,
		height: height,
	}
	w.init(maxInitLiveCells)
	return w
}

// init inits world with a random state.
func (w *World) init(maxLiveCells int) {
	for i := 0; i < maxLiveCells; i++ {
		x := rand.Intn(w.width)
		y := rand.Intn(w.height)
    val := rand.Intn(nImg) + 1
		w.area[y*w.width+x] = val
	}
}

// Update game state by one tick.
func (w *World) Update() {
	width := w.width
	height := w.height
	next := make([]int, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pop := neighbourCount(w.area, width, height, x, y)
			switch {
			case pop < 2:
				// rule 1. Any live cell with fewer than two live neighbours
				// dies, as if caused by under-population.
				next[y*width+x] = 0

			case (pop == 2 || pop == 3) && w.area[y*width+x] > 0:
				// rule 2. Any live cell with two or three live neighbours
				// lives on to the next generation.
				next[y*width+x] = rand.Intn(nImg) + 1

			case pop > 3:
				// rule 3. Any live cell with more than three live neighbours
				// dies, as if by over-population.
				next[y*width+x] = 0

			case pop == 3:
				// rule 4. Any dead cell with exactly three live neighbours
				// becomes a live cell, as if by reproduction.
				next[y*width+x] = rand.Intn(nImg) + 1
			}
		}
	}
	w.area = next
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// neighbourCount calculates the Moore neighborhood of (x, y).
func neighbourCount(a []int, width, height, x, y int) int {
	c := 0
	for j := -1; j <= 1; j++ {
		for i := -1; i <= 1; i++ {
			if i == 0 && j == 0 {
				continue
			}
			x2 := x + i
			y2 := y + j
			if x2 < 0 || y2 < 0 || width <= x2 || height <= y2 {
				continue
			}
			if a[y2*width+x2] > 0 {
				c++
			}
		}
	}
	return c
}

const (
	screenWidth  = 16
	screenHeight = 16
)

type Game struct {
	world  *World
	//pixels []byte
}

func (g *Game) Update() error {
	g.world.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
  for i, p := range g.world.area {
    if p > 0 {
      op := &ebiten.DrawImageOptions{}
      y := int(i / screenWidth)
      x := int(i % screenHeight)
      //op.GeoM.Scale(0.5,0.5)
      op.GeoM.Translate(float64(x*400),float64(y*400))
      if p == 1 {
        screen.DrawImage(img1, op)
      } else if p == 2{
        screen.DrawImage(img2, op)
      } else {
        screen.DrawImage(img3, op)
      }
    }
  }
  op := &ebiten.DrawImageOptions{}
  op.GeoM.Translate(0.0, float64(screenHeight/4*400))
  screen.DrawImage(words,op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth*400, screenHeight*400
}

func main() {
	g := &Game{
		world: NewWorld(screenWidth, screenHeight, int((screenWidth*screenHeight)/4)),
	}

  ebiten.SetMaxTPS(1)
	ebiten.SetWindowSize(screenWidth*20, screenHeight*20)
	ebiten.SetWindowTitle("Game of Life (Ebiten Demo)")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
