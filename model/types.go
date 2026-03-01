package model

type Point struct {
	X, Y int
}

type Rect struct {
	X, Y, W, H int
}

type Sprite struct {
	W, H   int
	Pixels []byte
}
