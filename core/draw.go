package core

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x11, 0x11, 0x11, 0xff})

	for y := range ScreenH/cell {
		for x := range ScreenW/cell {
			vector.StrokeRect(
				screen,
				float32(x*cell), float32(y*cell),
				cell, cell,
				1,
				color.RGBA{0x33, 0x33, 0x33, 0xff},
				false,
			)
		}
	}

	for _, p := range g.snake {
		vector.FillRect(
			screen,
			float32(p.X*cell), float32(p.Y*cell),
			cell, cell,
			color.RGBA{0x55, 0x55, 0x55, 0xff},
			false,
		)
	}

	vector.FillRect(
		screen,
		float32(g.food.X*cell), float32(g.food.Y*cell),
		cell, cell,
		color.RGBA{0xbb, 0xbb, 0xbb, 0xff},
		false,
	)

	if g.gameOver {
		ebitenutil.DebugPrint(screen, "Game Over! Press Enter to restart")
	}
}

func (g *Game) Layout(outW, outH int) (int, int) {
	return ScreenW, ScreenH
}
