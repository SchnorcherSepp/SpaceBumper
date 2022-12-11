package megagrid

import "math"

type Vector struct {
	X float64
	Y float64
}

func NewVector(x float64, y float64) *Vector {
	return &Vector{
		X: x,
		Y: y,
	}
}

func (v *Vector) Angle() float64 {
	return math.Atan2(v.Y, v.X)
}

func (v *Vector) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}
