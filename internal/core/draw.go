package core

import (
	"baolhq/snake/internal/assets"
	"baolhq/snake/internal/consts"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var gameOverFace = assets.LoadFont(assets.MainFont, 32)

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(consts.BackgroundColor)

	for y := range consts.ScreenHeight / consts.CellSize {
		for x := range consts.ScreenWidth / consts.CellSize {
			vector.StrokeRect(
				screen,
				float32(x*consts.CellSize), float32(y*consts.CellSize),
				consts.CellSize, consts.CellSize,
				1,
				consts.BorderColor,
				false,
			)
		}
	}

	g.Snake.Draw(screen)

	// draw food
	vector.FillRect(
		screen,
		float32(g.food.X*consts.CellSize),
		float32(g.food.Y*consts.CellSize),
		consts.CellSize, consts.CellSize,
		consts.FoodColor,
		false,
	)

	// game over text
	if g.gameOver {
		t := "GAME OVER"
		opt := &text.DrawOptions{}
		w, h := text.Measure(t, gameOverFace, 0)
		opt.GeoM.Translate(float64(consts.ScreenWidth/2-w/2), float64(consts.ScreenHeight/2-h/2))
		text.Draw(screen, t, gameOverFace, opt)
	}
}

func (g *Game) Layout(outW, outH int) (int, int) {
	return consts.ScreenWidth, consts.ScreenHeight
}
