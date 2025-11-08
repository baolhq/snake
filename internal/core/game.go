package core

import (
	"math/rand"
	"slices"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"baolhq/snake/internal/models"
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

	// remove chosen food cell from free
	g.freeCells = append(g.freeCells[:idx], g.freeCells[idx+1:]...)
}

func (g *Game) handleInput() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	if g.gameOver {
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			*g = *NewGame()
		}
		return nil
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) && g.dir.Y != 1 {
		g.dir = models.Point{X: 0, Y: -1}
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) && g.dir.Y != -1 {
		g.dir = models.Point{X: 0, Y: 1}
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) && g.dir.X != 1 {
		g.dir = models.Point{X: -1, Y: 0}
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) && g.dir.X != -1 {
		g.dir = models.Point{X: 1, Y: 0}
	}

	return nil
}

func (g *Game) Update() error {
	if err := g.handleInput(); err != nil {
		return err
	}

	// update timer
	g.timer += time.Millisecond * 16
	if g.timer < 120*time.Millisecond {
		return nil
	}
	g.timer = 0

	// movement
	head := g.snake[0]
	newHead := models.Point{X: head.X + g.dir.X, Y: head.Y + g.dir.Y}

	// kill on out-of-bounds
	if newHead.X < 0 || newHead.X >= ScreenW/cell || newHead.Y < 0 || newHead.Y >= ScreenH/cell {
		g.gameOver = true
		return nil
	}

	// will we grow this tick? (eating the food)
	willGrow := newHead == g.food

	// collision check:
	// allow moving into the current tail when the tail will move away (i.e. when not growing
	// and pendingGrowth does not block tail movement).
	if slices.Contains(g.snake, newHead) {
		tail := g.snake[len(g.snake)-1]
		// if newHead == tail and tail will be removed this tick, it's allowed
		tailWillBeRemoved := !willGrow && !(len(g.pendingGrowth) > 0 && tail == g.pendingGrowth[0])
		if !(newHead == tail && tailWillBeRemoved) {
			g.gameOver = true
			return nil
		}
	}

	g.freeCells = removePoints(g.freeCells, []models.Point{newHead})

	// prepend new head
	g.snake = append([]models.Point{newHead}, g.snake...)

	// eating: queue a pending growth spot and spawn new food
	if willGrow {
		g.pendingGrowth = append(g.pendingGrowth, g.food)
		g.spawnFood()
	}

	// decide whether to remove tail this tick:
	// if tail equals the first pendingGrowth spot, then consume the pendingGrowth and DON'T pop tail (grow)
	tail := g.snake[len(g.snake)-1]
	if len(g.pendingGrowth) > 0 && tail == g.pendingGrowth[0] {
		// consume pending growth: do not remove tail, just advance the queue
		g.pendingGrowth = g.pendingGrowth[1:]
	} else {
		// normal move: drop the tail
		g.snake = g.snake[:len(g.snake)-1]
		g.freeCells = append(g.freeCells, tail)
	}

	return nil
}
