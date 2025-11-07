package main

import (
	"baolhq/snake/core"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(core.ScreenW, core.ScreenH)
	ebiten.SetWindowTitle("Snake - Ebiten v2")

	if err := ebiten.RunGame(core.NewGame()); err != nil {
		log.Fatal(err)
	}
}
