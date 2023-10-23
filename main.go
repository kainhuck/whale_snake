package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math/rand"
)

var (
	screenWidth  = 480
	screenHeight = 390
	gridSize     = 30
	xGridCount   = screenWidth / gridSize
	yGridCount   = screenHeight / gridSize
)

const (
	dirNone = iota
	dirLeft
	dirRight
	dirUp
	dirDown
)

const (
	FastSpeed   = 2
	MiddleSpeed = 4
	SlowSpeed   = 6
)

type Position struct {
	X int
	Y int
}

type PositionWithColor struct {
	p Position
	c color.Color
}

type Game struct {
	moveDir   int
	snakeBody []PositionWithColor
	timer     int
	moveTime  int
	apple     PositionWithColor
	score     int
	best      int
	level     int
	stop      bool
	finished  bool
}

func (g *Game) randomApplePosition() Position {
	x := rand.Intn(xGridCount)
	y := rand.Intn(yGridCount)

OUTER:
	for {
		for _, each := range g.snakeBody {
			if each.p.X == x && each.p.Y == y {
				x = rand.Intn(xGridCount)
				y = rand.Intn(yGridCount)
				continue OUTER
			}
			return Position{
				X: x,
				Y: y,
			}
		}
	}
}

func (g *Game) canMove() bool {
	return g.timer%g.moveTime == 0
}

func (g *Game) collidesApple() bool {
	return g.snakeBody[0].p.X == g.apple.p.X && g.snakeBody[0].p.Y == g.apple.p.Y
}

func (g *Game) collidesSelf() bool {
	for i := 1; i < len(g.snakeBody); i++ {
		if g.snakeBody[0].p.X == g.snakeBody[i].p.X && g.snakeBody[0].p.Y == g.snakeBody[i].p.Y {
			return true
		}
	}
	return false
}

func (g *Game) reset() {
	g.snakeBody = []PositionWithColor{{
		p: Position{X: 6, Y: 8},
		c: color.RGBA{
			R: 255,
			G: 215,
			B: 0,
			A: 1,
		},
	}, {
		p: Position{X: 5, Y: 8},
		c: color.RGBA{
			R: uint8(rand.Intn(128) + 128),
			G: uint8(rand.Intn(128) + 128),
			B: uint8(rand.Intn(128) + 128),
			A: 1,
		},
	}, {
		p: Position{X: 4, Y: 8},
		c: color.RGBA{
			R: uint8(rand.Intn(128) + 128),
			G: uint8(rand.Intn(128) + 128),
			B: uint8(rand.Intn(128) + 128),
			A: 1,
		},
	}}
	g.timer = 0
	g.moveTime = SlowSpeed
	g.apple = PositionWithColor{
		p: g.randomApplePosition(),
		c: color.RGBA{
			R: uint8(rand.Intn(256)),
			G: uint8(rand.Intn(256)),
			B: uint8(rand.Intn(256)),
			A: 1,
		},
	}
	g.score = 0
	g.level = 1
	g.stop = false
	g.moveDir = dirRight
	g.finished = false
}

func (g *Game) collidesWall() bool {
	return g.snakeBody[0].p.X <= 0 || g.snakeBody[0].p.Y < 0 || g.snakeBody[0].p.X >= xGridCount || g.snakeBody[0].p.Y >= yGridCount
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.reset()
		return nil
	}

	if g.finished {
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		if g.moveDir != dirDown {
			g.moveDir = dirUp
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		if g.moveDir != dirUp {
			g.moveDir = dirDown
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		if g.moveDir != dirRight {
			g.moveDir = dirLeft
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		if g.moveDir != dirLeft {
			g.moveDir = dirRight
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.stop = !g.stop
	}

	if g.canMove() {

		if g.stop {
			return nil
		}

		if g.collidesWall() || g.collidesSelf() {
			g.finished = true
			return nil
		}

		if g.collidesApple() {
			g.snakeBody = append(g.snakeBody, PositionWithColor{c: g.apple.c})
			g.apple = PositionWithColor{
				p: g.randomApplePosition(),
				c: color.RGBA{
					R: uint8(rand.Intn(128) + 128),
					G: uint8(rand.Intn(128) + 128),
					B: uint8(rand.Intn(128) + 128),
					A: 1,
				},
			}
			g.score++
			if g.score > g.best {
				g.best = g.score
			}
			if g.score >= 10 && g.score < 20 {
				g.level = 2
				g.moveTime = MiddleSpeed
			} else if g.score >= 20 {
				g.level = 3
				g.moveTime = FastSpeed
			}
		}

		for i := len(g.snakeBody) - 1; i > 0; i-- {
			g.snakeBody[i].p.X = g.snakeBody[i-1].p.X
			g.snakeBody[i].p.Y = g.snakeBody[i-1].p.Y
		}

		switch g.moveDir {
		case dirRight:
			g.snakeBody[0].p.X++
		case dirLeft:
			g.snakeBody[0].p.X--
		case dirUp:
			g.snakeBody[0].p.Y--
		case dirDown:
			g.snakeBody[0].p.Y++
		}
	}

	g.timer++

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, v := range g.snakeBody {
		vector.DrawFilledRect(screen, float32(v.p.X*gridSize), float32(v.p.Y*gridSize), float32(gridSize), float32(gridSize), v.c, false)
	}

	vector.DrawFilledCircle(screen, float32(g.apple.p.X*gridSize)+float32(gridSize)/2, float32(g.apple.p.Y*gridSize)+float32(gridSize)/2, float32(gridSize)/2, g.apple.c, true)

	if !g.finished {
		if g.stop {
			ebitenutil.DebugPrint(screen, fmt.Sprintf("Score: %v, Best: %v, Level: %v, press space to continue", g.score, g.best, g.level))
		} else {
			ebitenutil.DebugPrint(screen, fmt.Sprintf("Score: %v, Best: %v, Level: %v, press space to stop", g.score, g.best, g.level))
		}
	} else {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Score: %v, Best: %v, Level: %v, press esc to reset", g.score, g.best, g.level))
	}

	return
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func newGame() *Game {
	g := &Game{}
	g.reset()
	return g
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("帷幄五彩斑斓贪吃蛇")
	_ = ebiten.RunGame(newGame())
}
