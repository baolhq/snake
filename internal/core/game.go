package core

import (
	"math"
	"math/rand"
	"slices"
	"time"

	"baolhq/snake/internal/assets"
	"baolhq/snake/internal/consts"
	mng "baolhq/snake/internal/managers"
	"baolhq/snake/internal/models"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Snake      *models.Snake
	dir        models.Point
	pendingDir []models.Point
	food       models.Point
	freeCells  []models.Point
	timer      time.Duration
	accel      bool
	accelTimer time.Duration
	accelDelay time.Duration
	gameOver   bool
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
	totalX, totalY := consts.ScreenWidth/consts.CellSize, consts.ScreenHeight/consts.CellSize

	free := make([]models.Point, 0, totalX*totalY)
	for y := range totalY {
		for x := range totalX {
			free = append(free, models.Point{X: x, Y: y})
		}
	}

	initial := []models.Point{
		{X: totalX/2 + 1, Y: totalY / 2},
		{X: totalX / 2, Y: totalY / 2},
		{X: totalX/2 - 1, Y: totalY / 2},
	}
	s := models.NewSnake(initial, consts.CellSize)
	free = removePoints(free, s.Segments)

	g := &Game{
		Snake:      s,
		dir:        models.Point{X: 1, Y: 0},
		freeCells:  free,
		accelDelay: 200 * time.Millisecond,
		pendingDir: make([]models.Point, 0, 2),
	}
	g.spawnFood()
	return g
}

func (g *Game) spawnFood() {
	if len(g.freeCells) == 0 {
		return
	}
	idx := rand.Intn(len(g.freeCells))
	g.food = g.freeCells[idx]
	g.freeCells = append(g.freeCells[:idx], g.freeCells[idx+1:]...)
}

func (g *Game) queueDirection(action mng.Action, dx, dy int) {
	if !mng.Input.WasPressed(action) {
		return
	}

	last := g.dir
	if len(g.pendingDir) > 0 {
		last = g.pendingDir[len(g.pendingDir)-1]
	}

	if last.X == -dx && last.Y == -dy {
		return // cannot reverse
	}

	if len(g.pendingDir) < 2 {
		g.pendingDir = append(g.pendingDir, models.Point{X: dx, Y: dy})
	}
}

func (g *Game) updateAccelState() {
	g.accel = false

	check := func(dx, dy int, action mng.Action) bool {
		if mng.Input.IsDown(action) {
			if g.dir.X == dx && g.dir.Y == dy {
				return true
			}
			if len(g.pendingDir) > 0 && g.pendingDir[0].X == dx && g.pendingDir[0].Y == dy {
				return true
			}
		}
		return false
	}

	g.accel = check(1, 0, mng.ActionRight) ||
		check(-1, 0, mng.ActionLeft) ||
		check(0, 1, mng.ActionDown) ||
		check(0, -1, mng.ActionUp)
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

	g.queueDirection(mng.ActionUp, 0, -1)
	g.queueDirection(mng.ActionDown, 0, 1)
	g.queueDirection(mng.ActionLeft, -1, 0)
	g.queueDirection(mng.ActionRight, 1, 0)

	g.updateAccelState()
	return nil
}

func (g *Game) computeInterval() time.Duration {
	base, min := 200.0, 50.0
	k := 0.0366
	score := float64(len(g.Snake.Segments))
	interval := min + (base-min)*math.Exp(-k*score)
	if interval < min {
		interval = min
	}

	if g.accel {
		g.accelTimer += 16 * time.Millisecond
		if g.accelTimer >= g.accelDelay {
			interval = 50
		}
	} else {
		g.accelTimer = 0
	}

	return time.Duration(interval) * time.Millisecond
}

func (g *Game) applyPendingDirection() {
	if len(g.pendingDir) > 0 {
		g.dir = g.pendingDir[0]
		g.pendingDir = g.pendingDir[1:]
	}
}

func (g *Game) Update() error {
	if err := g.HandleInput(); err != nil {
		return err
	}

	mng.Particle.Update(1.0 / 60)

	g.timer += 16 * time.Millisecond
	if g.timer < g.computeInterval() {
		return nil
	}
	g.timer = 0

	g.applyPendingDirection()

	selfCollision, ateFood := g.Snake.Move(g.dir, &g.freeCells, g.food)
	if selfCollision {
		g.gameOver = true
	}

	if ateFood {
		g.spawnFood()

		px := g.Snake.Segments[0].X * consts.CellSize
		py := g.Snake.Segments[0].Y * consts.CellSize
		mng.Particle.Explode(assets.ParticleImage, px, py, 20)
	}

	return nil
}
