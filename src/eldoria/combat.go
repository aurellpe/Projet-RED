package main

type Rect struct {
	X float64
	Y float64
	W float64
	H float64
}

func Intersects(a Rect, b Rect) bool {
	return a.X < b.X+b.W &&
		a.X+a.W > b.X &&
		a.Y < b.Y+b.H &&
		a.Y+a.H > b.Y
}