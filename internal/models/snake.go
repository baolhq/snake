package models

func NewSnake(screenW, cells int) []Point {
	center := (screenW / cells) / 2

	return []Point{
		{center, center},
		{center - 1, center},
		{center - 2, center},
	}
}

func MoveSnake(snake []Point, newHead Point, grow bool) []Point {
	snake = append([]Point{newHead}, snake...)
	if !grow {
		snake = snake[:len(snake)-1]
	}
	return snake
}
