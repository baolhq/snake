package core

import (
	"math/rand"
	"slices"
	"time"

	mng "baolhq/snake/internal/managers"
	"baolhq/snake/internal/models"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenW = 480
	ScreenH = 480
	cell    = 20
)

type Game struct {
	snake         []models.Point
	pendingGrowth []models.Point
	freeCells     []models.Point
	dir           models.Point
	food          models.Point
	timer         time.Duration
	gameOver      bool
}

func removePoints(list []models.Point, points []models.Point) []models.Point {
	out := list[:0]
	for _, p := range list {
		if !slices.Contains(points, p) {
			out = append(out, p)
		}
	}
	return out
}

func NewGame() *Game {
	totalX := ScreenW / cell
	totalY := ScreenH / cell

	free := make([]models.Point, 0, totalX*totalY)
	for y := range totalY {
		for x := range totalX {
			free = append(free, models.Point{X: x, Y: y})
		}
	}

	s := models.NewSnake(ScreenW, cell)
	free = removePoints(free, s)

	g := &Game{
		snake:     s,
		dir:       models.Point{X: 1, Y: 0},
		freeCells: free,
	}
	g.spawnFood()
	return g
}

func (g *Game) spawnFood() {
	idx := rand.Intn(len(g.freeCells))
	g.food = g.freeCells[idx]
	g.freeCells = append(g.freeCells[:idx], g.freeCells[idx+1:]...)
}

func (g *Game) HandleInput() error {
	mng.Input.Update()

	if mng.Input.WasPressed(mng.ActionPause) {
		return ebiten.Termination
	}
	if g.gameOver && mng.Input.WasPressed(mng.ActionEnter) {
		*g = *NewGame()
		return nil
	}

	switch {
	case mng.Input.WasPressed(mng.ActionUp) && g.dir.Y != 1:
		g.dir = models.Point{X: 0, Y: -1}
	case mng.Input.WasPressed(mng.ActionDown) && g.dir.Y != -1:
		g.dir = models.Point{X: 0, Y: 1}
	case mng.Input.WasPressed(mng.ActionLeft) && g.dir.X != 1:
		g.dir = models.Point{X: -1, Y: 0}
	case mng.Input.WasPressed(mng.ActionRight) && g.dir.X != -1:
		g.dir = models.Point{X: 1, Y: 0}
	}

	return nil
}

func (g *Game) Update() error {
	if err := g.HandleInput(); err != nil {
		return err
	}

	g.timer += time.Millisecond * 16
	if g.timer < 120*time.Millisecond {
		return nil
	}
	g.timer = 0

	head := g.snake[0]
	newHead := models.Point{X: head.X + g.dir.X, Y: head.Y + g.dir.Y}

	if newHead.X < 0 || newHead.X >= ScreenW/cell || newHead.Y < 0 || newHead.Y >= ScreenH/cell {
		g.gameOver = true
		return nil
	}

	willGrow := newHead == g.food

	if slices.Contains(g.snake, newHead) {
		tail := g.snake[len(g.snake)-1]
		tailWillBeRemoved := !willGrow && !(len(g.pendingGrowth) > 0 && tail == g.pendingGrowth[0])
		if !(newHead == tail && tailWillBeRemoved) {
			g.gameOver = true
			return nil
		}
	}

	g.freeCells = removePoints(g.freeCells, []models.Point{newHead})
	g.snake = append([]models.Point{newHead}, g.snake...)

	if willGrow {
		g.pendingGrowth = append(g.pendingGrowth, g.food)
		g.spawnFood()
	}

	tail := g.snake[len(g.snake)-1]
	if len(g.pendingGrowth) > 0 && tail == g.pendingGrowth[0] {
		g.pendingGrowth = g.pendingGrowth[1:]
	} else {
		g.snake = g.snake[:len(g.snake)-1]
		g.freeCells = append(g.freeCells, tail)
	}

	return nil
}
