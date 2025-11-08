package main

import (
	"baolhq/snake/internal/core"
	"baolhq/snake/internal/consts"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(consts.ScreenWidth, consts.ScreenHeight)
	ebiten.SetWindowTitle("Snake - Ebiten v2")

	if err := ebiten.RunGame(core.NewGame()); err != nil {
		log.Fatal(err)
	}
}
