package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
	"image"
	"image/color"
	"log"
	"math/rand"
	"snake/input"
	"sync"
)

const (
	BOARDWIDTH  = 10
	BOARDHEIGHT = 10
	TILESIZE    = 64
)

type Board struct {
	Grid           [BOARDHEIGHT][BOARDWIDTH]int
	Fruit          [2]int
	SnakeLength    int
	SnakeDirection int
	SnakePlace     [2]int
}

type Game struct {
	Input *input.Manager
	Board *Board
}

func NewBoard() *Board {
	var b Board
	b.Fruit = [2]int{2, 2}
	b.SnakePlace = [2]int{5, 5}
	b.SnakeLength = 3
	b.SnakeDirection = 1
	return &b
}

func (g *Game) MoveFruit() {
	log.Printf("fruit at %+v", g.Board.Fruit)
	for g.Board.Grid[g.Board.Fruit[0]][g.Board.Fruit[1]] > 0 {
		g.Board.Fruit[1] = rand.Intn(BOARDWIDTH)
		g.Board.Fruit[0] = rand.Intn(BOARDHEIGHT)
		log.Printf("moved fruit to %+v", g.Board.Fruit)

	}
}

func (g *Game) AteFruit() {
	g.Board.SnakeLength++
	for i := 0; i < BOARDHEIGHT; i++ {
		for j := 0; j < BOARDWIDTH; j++ {
			if g.Board.Grid[i][j] > 0 {
				g.Board.Grid[i][j] = g.Board.Grid[i][j] + 1
			}
		}
	}
	g.MoveFruit()
}

var frame int

func (g *Game) Update() (err error) {
	g.Input.Update()
	frame++
	if frame%30 == 0 {

		for i := 0; i < BOARDHEIGHT; i++ {
			for j := 0; j < BOARDWIDTH; j++ {
				if g.Board.Grid[i][j] > 0 {
					g.Board.Grid[i][j] = g.Board.Grid[i][j] - 1
				}
			}
		}
		if g.Board.SnakePlace[0] < BOARDHEIGHT && g.Board.SnakePlace[1] < BOARDWIDTH {
			g.Board.Grid[g.Board.SnakePlace[0]][g.Board.SnakePlace[1]] = g.Board.SnakeLength
		}

		if g.Board.SnakePlace[0] == g.Board.Fruit[0] && g.Board.SnakePlace[1] == g.Board.Fruit[1] {
			g.AteFruit()
		}
		switch g.Board.SnakeDirection {
		case 1:
			g.Board.SnakePlace[1] += 1
		case 2:
			g.Board.SnakePlace[0] += 1
		case 3:
			g.Board.SnakePlace[1] -= 1
		case 4:
			g.Board.SnakePlace[0] -= 1

		}
		if g.Board.SnakePlace[0] >= BOARDHEIGHT || g.Board.SnakePlace[1] >= BOARDWIDTH || g.Board.SnakePlace[0] < 0 || g.Board.SnakePlace[1] < 0 ||
			g.Board.Grid[g.Board.SnakePlace[0]][g.Board.SnakePlace[1]] > 0 {
			panic(nil)
		}
	}

inputloop:
	for {
		select {
		case KeyEvent := <-g.Input.Stream:
			log.Println(KeyEvent)
			switch KeyEvent.Key {
			case ebiten.KeyArrowRight:
				if g.Board.SnakeDirection == 2 || g.Board.SnakeDirection == 4 {
					g.Board.SnakeDirection = 1
				}
			case ebiten.KeyArrowDown:
				if g.Board.SnakeDirection == 1 || g.Board.SnakeDirection == 3 {
					g.Board.SnakeDirection = 2
				}
			case ebiten.KeyArrowLeft:
				if g.Board.SnakeDirection == 2 || g.Board.SnakeDirection == 4 {
					g.Board.SnakeDirection = 3
				}
			case ebiten.KeyArrowUp:
				if g.Board.SnakeDirection == 1 || g.Board.SnakeDirection == 3 {
					g.Board.SnakeDirection = 4
				}
			}
		default:
			break inputloop
		}
	}
	return
}

func (g *Game) Draw(screen *ebiten.Image) {
	fruitX := g.Board.Fruit[1] * TILESIZE
	fruitY := g.Board.Fruit[0] * TILESIZE

	fruitCanvas := screen.SubImage(image.Rect(fruitX, fruitY, fruitX+TILESIZE, fruitY+TILESIZE))
	fruitCanvas.(*ebiten.Image).Fill(colornames.Red)

	for i := 0; i < BOARDHEIGHT; i++ {
		for j := 0; j < BOARDWIDTH; j++ {
			if g.Board.Grid[i][j] > 0 {
				colStart := j * TILESIZE
				rowStart := i * TILESIZE
				objectCanvas := screen.SubImage(image.Rect(colStart, rowStart, colStart+TILESIZE, rowStart+TILESIZE))
				objectCanvas.(*ebiten.Image).Fill(color.White)
			}
		}
	}
}

func (g *Game) Layout(w, h int) (gw, gh int) {
	return 64 * BOARDWIDTH, 64 * BOARDHEIGHT
}

var startupOnce sync.Once

func (g *Game) Initialize() {
	startupOnce.Do(func() {
		log.SetFlags(log.Ltime | log.Llongfile)
		log.Printf("logs initialized")
		g.Input = input.New()

		// register keys that we want to listen for in combos
		g.Input.RegisterKey(ebiten.KeyArrowUp)
		g.Input.RegisterKey(ebiten.KeyArrowDown)
		g.Input.RegisterKey(ebiten.KeyArrowLeft)
		g.Input.RegisterKey(ebiten.KeyArrowRight)
		g.Board = NewBoard()

	})
}

func main() {
	var snake Game
	snake.Initialize()
	ebiten.SetWindowTitle("snake")
	if err := ebiten.RunGame(&snake); err != nil {
		log.Fatal(err)
	}
	fmt.Println()
}
