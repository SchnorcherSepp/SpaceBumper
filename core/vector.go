package core

import "math"

// Vector is a geometric object that has length and direction.
// Stores the X and Y values inside.
type Vector struct {
	x float64
	y float64
}

// NewVector return a new Vector
func NewVector(x float64, y float64) *Vector {
	return &Vector{
		x: x,
		y: y,
	}
}

//--------  Getter  --------------------------------------------------------------------------------------------------//

// X return the x-dimensional of this vector
func (v *Vector) X() float64 {
	return v.x
}

// Y return the y-dimensional of this vector
func (v *Vector) Y() float64 {
	return v.y
}

// Clone returns a new, identical vector
func (v *Vector) Clone() *Vector {
	return &Vector{v.x, v.y}
}

// Angle return the angle in radians (rad).
func (v *Vector) Angle() float64 {
	return math.Atan2(v.y, v.x)
}

// Length returns the vector length:
// Sqrt(X*X + Y*Y)
func (v *Vector) Length() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y)
}

//--------  Setter  --------------------------------------------------------------------------------------------------//

// Normalize this vector maintains its direction
// but its length becomes 1.
//
// Attention: The method modifies this vector.
func (v *Vector) Normalize() {
	l := v.Length()
	v.x = v.x / l
	v.y = v.y / l
}

// Add multiplies the passed vector and then adds it to this vector:
// this = this + other*multi
//
// Attention: The method modifies this vector.
func (v *Vector) Add(o *Vector, multi float64) {
	v.x += o.x * multi
	v.y += o.y * multi
}

// Multi this vector:
// this = this * multi
//
// Attention: The method modifies this vector.
func (v *Vector) Multi(m float64) {
	v.x *= m
	v.y *= m
}
