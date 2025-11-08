package models

import (
	"baolhq/snake/internal/consts"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Snake struct {
	Segments      []Point
	PendingGrowth []Point
	CellSize      int
}

func NewSnake(initial []Point, cell int) *Snake {
	return &Snake{
		Segments: initial,
		CellSize: cell,
	}
}

func removePoints(list []Point, points []Point) []Point {
	out := list[:0]
	for _, p := range list {
		if !slices.Contains(points, p) {
			out = append(out, p)
		}
	}
	return out
}

func (s *Snake) Move(dir Point, freeCells *[]Point, food Point) (selfCollision bool, ateFood bool) {
	newHead := Point{X: s.Segments[0].X + dir.X, Y: s.Segments[0].Y + dir.Y}

	// Wrap around edges
	totalX := consts.ScreenWidth / consts.CellSize
	totalY := consts.ScreenHeight / consts.CellSize
	if newHead.X < 0 {
		newHead.X = totalX - 1
	} else if newHead.X >= totalX {
		newHead.X = 0
	}
	if newHead.Y < 0 {
		newHead.Y = totalY - 1
	} else if newHead.Y >= totalY {
		newHead.Y = 0
	}

	*freeCells = removePoints(*freeCells, []Point{newHead})

	// Check self-collision
	willGrow := newHead == food
	if slices.Contains(s.Segments, newHead) {
		tail := s.Segments[len(s.Segments)-1]
		tailWillBeRemoved := !willGrow && !(len(s.PendingGrowth) > 0 && tail == s.PendingGrowth[0])
		if !(newHead == tail && tailWillBeRemoved) {
			return true, willGrow
		}
	}

	// Add new head
	s.Segments = append([]Point{newHead}, s.Segments...)

	// Handle growth
	if willGrow {
		s.PendingGrowth = append(s.PendingGrowth, food)
	}

	// Remove tail if not growing
	tail := s.Segments[len(s.Segments)-1]
	if len(s.PendingGrowth) > 0 && tail == s.PendingGrowth[0] {
		s.PendingGrowth = s.PendingGrowth[1:]
	} else {
		s.Segments = s.Segments[:len(s.Segments)-1]
		*freeCells = append(*freeCells, tail)
	}

	return false, willGrow
}

func (s *Snake) Draw(screen *ebiten.Image) {
	if len(s.Segments) == 0 {
		return
	}

	vector.FillRect(
		screen,
		float32(s.Segments[0].X*s.CellSize),
		float32(s.Segments[0].Y*s.CellSize),
		float32(s.CellSize),
		float32(s.CellSize),
		consts.SnakeHeadColor,
		false,
	)

	for i := 1; i < len(s.Segments); i++ {
		vector.FillRect(
			screen,
			float32(s.Segments[i].X*s.CellSize),
			float32(s.Segments[i].Y*s.CellSize),
			float32(s.CellSize),
			float32(s.CellSize),
			consts.SnakeBodyColor,
			false,
		)
	}
}
