package model

type Point struct {
	X, Y int
}

type Size struct {
	W, H int
}

type Rect struct {
	Point Point
	Size  Size
	//	X, Y, W, H int
}

type Sprite struct {
	W, H   int
	Pixels []byte
}
