package core

import (
	"baolhq/snake/internal/assets"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var gameOverFont font.Face
var gameOverFace text.Face

func init() {
	fnt, err := opentype.Parse(assets.MainFont)
	if err != nil {
		log.Fatal(err)
	}
	gameOverFont, err = opentype.NewFace(fnt, &opentype.FaceOptions{
		Size: 32, // adjust size
		DPI:  72,
	})
	gameOverFace = text.NewGoXFace(gameOverFont)
	if err != nil {
		log.Fatal(err)
	}
}

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
		opt := &text.DrawOptions{}
		opt.GeoM.Translate(ScreenW/2-80, ScreenH/2-80)
		text.Draw(screen, "GAME OVER", gameOverFace, opt)
	}
}

func (g *Game) Layout(outW, outH int) (int, int) {
	return ScreenW, ScreenH
}
