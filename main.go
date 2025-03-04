// -----------------------------------------------------------------------------
//
// -- Sun-Tower Game (2025)
// -- SunRiver https://forum.lothar-team.pl,
// -- wersja 0.1.2
// -- Obracanie wieży   kursory lewo /prawo
// -- celowanie i strzał kulą myszka + LMB

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
	nextBall  Ball // Nowa kula na górze ekranu
	ballsLeft int
	score     int
}

type Block struct {
	angle, y float64
	color    int
	active   bool
}

type Ball struct {
	x, y, vx, vy, radius float64
	color                int
	active               bool
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
		g.checkCollisions() // Dodanie sprawdzania kolizji
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.state == "title" {
		ebitenutil.DebugPrint(screen, "Sun Tower\nPress SPACE to Start")
	} else if g.state == "playing" {
		// Rysowanie bloków z cieniowaniem
		for _, block := range g.tower {
			if block.active {
				x := centerX + radius*math.Cos(block.angle+g.angle)
				y := block.y
				blockImg := ebiten.NewImage(40, 40)
				blockImg.Fill(getColor(block.color))
				opt := &ebiten.DrawImageOptions{}
				opt.GeoM.Translate(x, y)
				// Dodanie efektu cienia (proste 3D)
				shadow := ebiten.NewImage(45, 45)
				shadow.Fill(color.RGBA{0, 0, 0, 100})
				shadowOpt := &ebiten.DrawImageOptions{}
				shadowOpt.GeoM.Translate(x+2, y+2)
				screen.DrawImage(shadow, shadowOpt)
				screen.DrawImage(blockImg, opt)
			}
		}

		// Rysowanie aktywnej kuli
		if g.ball.active {
			g.drawBall(screen, g.ball)
		}
		// Rysowanie następnej kuli na górze ekranu
		g.drawBall(screen, g.nextBall)
		// Rysowanie kuli gotowej do strzału na dole ekranu
		g.drawBall(screen, Ball{x: float64(centerX), y: ballStartY, radius: g.ball.radius, color: g.ball.color, active: true})
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Score: %d\nBalls Left: %d", g.score, g.ballsLeft))
	}
}

func (g *Game) throwBall() {
	g.ball = g.nextBall
	g.ball.x = float64(centerX)
	g.ball.y = ballStartY
	mouseX, mouseY := ebiten.CursorPosition()
	dx := float64(mouseX - centerX)
	dy := float64(mouseY - ballStartY)
	dist := math.Hypot(dx, dy)
	g.ball.vx = (dx / dist) * ballSpeed
	g.ball.vy = (dy / dist) * ballSpeed
	g.ball.active = true
	g.ballsLeft--
	g.nextBall = generateRandomBall()
}

func (g *Game) drawBall(screen *ebiten.Image, ball Ball) {
	if ball.active {
		ballImg := ebiten.NewImage(int(ball.radius*2), int(ball.radius*2))
		ballImg.Fill(getColor(ball.color))
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(ball.x-ball.radius, ball.y-ball.radius)
		screen.DrawImage(ballImg, opt)
	}
}

func (g *Game) updateBall() {
	if g.ball.active {
		g.ball.x += g.ball.vx
		g.ball.y += g.ball.vy
		if g.ball.y < 0 {
			g.ball.active = false
		}
	}
}

func (g *Game) startGame() {
	g.state = "playing"
	g.ballsLeft = 16
	g.score = 0
	g.tower = generateTower()
	g.ball = Ball{radius: randFloat(8, 15)}
	g.nextBall = generateRandomBall()
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return Width, Height
}

func randFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func generateRandomBall() Ball {
	return Ball{
		radius: randFloat(8, 15),
		color:  rand.Intn(4),
		active: true,
	}
}

// -- stara funkcja -------------------------------
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

// Sprawdzanie kolizji kuli z blokami
/*
func (g *Game) checkCollisions() {
	for i := range g.tower {
		block := &g.tower[i]
		if block.active && g.ball.active {
			// Sprawdzamy, czy kula zderza się z blokiem
			x := centerX + radius*math.Cos(block.angle+g.angle)
			y := block.y

			// Sprawdzamy kolizję na prostokątnym bloku
			if g.ball.x+g.ball.radius > x && g.ball.x-g.ball.radius < x+30 &&
				g.ball.y+g.ball.radius > y && g.ball.y-g.ball.radius < y+30 {
				block.active = false  // Zniszczenie bloku
				g.score++             // Zwiększenie wyniku
				g.ball.active = false // Kula przestaje istnieć
			}
		}
	}
}
*/
// ---- Poprawiona kolizja ....
func (g *Game) checkCollisions() {
	for i := range g.tower {
		block := &g.tower[i]
		if block.active && g.ball.active {
			// Sprawdzamy, czy kula zderza się z blokiem
			x := centerX + radius*math.Cos(block.angle+g.angle)
			y := block.y

			// Sprawdzamy kolizję na prostokątnym bloku
			if g.ball.x+g.ball.radius > x && g.ball.x-g.ball.radius < x+30 &&
				g.ball.y+g.ball.radius > y && g.ball.y-g.ball.radius < y+30 {
				block.active = false  // Zniszczenie bloku
				g.score++             // Zwiększenie wyniku
				g.ball.active = false // Kula przestaje istnieć

				// Sprawdzanie innych bloków w tym samym kolorze
				for j := range g.tower {
					otherBlock := &g.tower[j]
					if otherBlock.active && block.color == otherBlock.color {
						// Obliczanie odległości między blokami
						otherX := centerX + radius*math.Cos(otherBlock.angle+g.angle)
						otherY := otherBlock.y

						// Sprawdzamy, czy blok znajduje się blisko
						if math.Hypot(otherX-x, otherY-y) < 50 { // Możesz dostosować 50 do odpowiedniej odległości
							otherBlock.active = false // Zniszczenie tego bloku
							g.score++                 // Zwiększenie wyniku
						}
					}
				}
			}
		}
	}
}

func getColor(index int) color.Color {
	switch index {
	case 0:
		return color.RGBA{255, 0, 0, 255} // Czerwony
	case 1:
		return color.RGBA{0, 255, 0, 255} // Zielony
	case 2:
		return color.RGBA{0, 0, 255, 255} // Niebieski
	case 3:
		return color.RGBA{255, 255, 0, 255} // Żółty
	default:
		return color.Black
	}
}

func main() {
	game := &Game{state: "title"}
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
