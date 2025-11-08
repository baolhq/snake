package core

import (
	"baolhq/snake/internal/assets"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var gameOverFace = assets.LoadFont(assets.MainFont, 32)

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x11, 0x11, 0x11, 0xff})

	for y := range ScreenH / cell {
		for x := range ScreenW / cell {
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
		t := "GAME OVER"
		opt := &text.DrawOptions{}
		w, h := text.Measure(t, gameOverFace, 0)
		opt.GeoM.Translate(float64(ScreenW/2-w/2), float64(ScreenH/2-h/2))

		text.Draw(screen, t, gameOverFace, opt)
	}
}

func (g *Game) Layout(outW, outH int) (int, int) {
	return ScreenW, ScreenH
}
