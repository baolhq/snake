package core

import (
	"baolhq/snake/internal/assets"
	"baolhq/snake/internal/consts"
	mng "baolhq/snake/internal/managers"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var mainFont = assets.LoadFont(assets.MainFont, 48)
var subFont = assets.LoadFont(assets.MainFont, 18)

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

	g.snake.Draw(screen)

	// draw food
	vector.FillRect(
		screen,
		float32(g.food.X*consts.CellSize),
		float32(g.food.Y*consts.CellSize),
		consts.CellSize, consts.CellSize,
		consts.FoodColor,
		false,
	)

	mng.Particle.Draw(screen)

	if mng.State.Is(mng.GamePaused) {
		drawTextBackground(screen)
		drawCenteredText(screen, mainFont, "PAUSED", -12)
		drawCenteredText(screen, subFont, "PRESS <ENTER> TO CONTINUE", 24)
	}

	if mng.State.Is(mng.GameOver) {
		drawTextBackground(screen)
		drawCenteredText(screen, mainFont, "GAME OVER", -24)
		drawCenteredText(screen, subFont, "PRESS <ENTER> TO RESTART", 24)
	}
}

func drawTextBackground(screen *ebiten.Image) {
	w, h := float32(consts.ScreenWidth), float32(100)

	vector.FillRect(
		screen,
		0,
		consts.ScreenHeight/2-h/2,
		w, h,
		consts.PanelColor,
		false,
	)
}

func drawCenteredText(screen *ebiten.Image, font text.Face, content string, yOffset float64) {
	t := content
	opt := &text.DrawOptions{}
	w, h := text.Measure(t, font, 0)
	opt.GeoM.Translate(float64(consts.ScreenWidth/2-w/2), float64(consts.ScreenHeight/2-h/2)+yOffset)
	text.Draw(screen, t, font, opt)
}

func (g *Game) Layout(outW, outH int) (int, int) {
	return consts.ScreenWidth, consts.ScreenHeight
}
