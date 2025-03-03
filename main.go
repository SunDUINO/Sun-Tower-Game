// -----------------------------------------------------------------------------
//
// -- Sun-Tower Game (2025)
// -- SunRiver https://forum.lothar-team.pl,
// -- versja 0.1.1
// -- Obracanie wieży   kursory lewo /prawo
// -- celowanie i styrzał kulą myszka + LMB

package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	state     string
	angle     float64
	tower     []Block
	ball      Ball
	ballsLeft int
	score     int
}

type Block struct {
	angle, y float64
	color    int
	active   bool
}

type Ball struct {
	x, y, vx, vy float64
	color        int
	active       bool
}

const (
	Width       = 800
	Height      = 600
	centerX     = Width / 2
	baseY       = Height - 200
	radius      = 100
	numFloors   = 11
	numColumns  = 12
	totalBlocks = numFloors * numColumns
	ballSpeed   = 5
	ballStartY  = Height - 50
)

// odswierzanie gry
func (g *Game) Update() error {
	if g.state == "title" {
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.startGame()
		}
	} else if g.state == "playing" {
		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			g.angle -= 0.05
		}
		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			g.angle += 0.05
		}
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && g.ballsLeft > 0 && !g.ball.active {
			g.throwBall()
		}
		g.updateBall()
	}
	return nil
}

// --- rysowanie gry
func (g *Game) Draw(screen *ebiten.Image) {
	if g.state == "title" {
		ebitenutil.DebugPrint(screen, "Tower Breaker\nPress SPACE to Start")
	} else if g.state == "playing" {
		for _, block := range g.tower {
			if block.active {
				x := centerX + radius*math.Cos(block.angle+g.angle)
				y := block.y
				blockImg := ebiten.NewImage(30, 30)
				blockImg.Fill(getColor(block.color))
				opt := &ebiten.DrawImageOptions{}
				opt.GeoM.Translate(x, y)
				screen.DrawImage(blockImg, opt)
			}
		}
		if g.ball.active {
			ballImg := ebiten.NewImage(10, 10)
			ballImg.Fill(getColor(g.ball.color))
			opt := &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(g.ball.x-5, g.ball.y-5)
			screen.DrawImage(ballImg, opt)
		}
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Score: %d\nBalls Left: %d", g.score, g.ballsLeft))
	}
}

// --- game start Spacja
func (g *Game) startGame() {
	g.state = "playing"
	g.ballsLeft = 16
	g.score = 0
	g.tower = generateTower()
	g.ball = Ball{}
}

func (g *Game) throwBall() {
	mouseX, mouseY := ebiten.CursorPosition()
	dx := float64(mouseX - centerX)
	dy := float64(mouseY - ballStartY)
	dist := math.Hypot(dx, dy)
	g.ball = Ball{
		x:      float64(centerX),
		y:      float64(ballStartY),
		vx:     (dx / dist) * ballSpeed,
		vy:     (dy / dist) * ballSpeed,
		color:  rand.Intn(4),
		active: true,
	}
	g.ballsLeft--
}

func (g *Game) updateBall() {
	if g.ball.active {
		g.ball.x += g.ball.vx
		g.ball.y += g.ball.vy
		for j := range g.tower {
			if g.tower[j].active && checkCollision(g.ball, g.tower[j], g.angle) {
				if g.ball.color == g.tower[j].color {
					g.destroyConnectedBlocks(j)
				} else {
					g.ball.vx = -g.ball.vx
					g.ball.vy = -g.ball.vy
				}
				g.ball.active = false
				break
			}
		}
	}
}

func checkCollision(ball Ball, block Block, angle float64) bool {
	x := centerX + radius*math.Cos(block.angle+angle)
	y := block.y
	return math.Hypot(ball.x-x, ball.y-y) < 20
}

func (g *Game) destroyConnectedBlocks(index int) {
	color := g.tower[index].color
	g.destroyBlock(index, color)
}

func (g *Game) destroyBlock(index, color int) {
	if index < 0 || index >= len(g.tower) || !g.tower[index].active || g.tower[index].color != color {
		return
	}
	g.tower[index].active = false
	g.score += 10
	directions := []int{-1, 1, -numColumns, numColumns}
	for _, d := range directions {
		g.destroyBlock(index+d, color)
	}
}

func generateTower() []Block {
	rand.Seed(time.Now().UnixNano())
	tower := make([]Block, totalBlocks)
	for floor := 0; floor < numFloors; floor++ {
		for col := 0; col < numColumns; col++ {
			angle := (2 * math.Pi / float64(numColumns)) * float64(col)
			tower[floor*numColumns+col] = Block{
				angle:  angle,
				y:      float64(baseY - floor*30),
				color:  rand.Intn(4),
				active: true,
			}
		}
	}
	return tower
}

func getColor(index int) color.Color {
	switch index {
	case 0:
		return color.RGBA{255, 0, 0, 255}
	case 1:
		return color.RGBA{0, 255, 0, 255}
	case 2:
		return color.RGBA{0, 0, 255, 255}
	case 3:
		return color.RGBA{255, 255, 0, 255}
	default:
		return color.Black
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return Width, Height
}

func main() {
	game := &Game{state: "title"}
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
