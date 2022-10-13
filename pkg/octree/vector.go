package octree

import (
	"fmt"
	"math"
)

type Point struct {
	Pos     *Vector3
	R, G, B byte
}

type Vector3 struct {
	X, Y, Z float64
}

func NewVector3(value float64) *Vector3 {
	return &Vector3{
		X: value,
		Y: value,
		Z: value,
	}
}

func (v *Vector3) Len() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v *Vector3) SquaredLen() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v *Vector3) DistanceTo(vec *Vector3) float64 {
	return v.Minus(vec).Len()
}

func (v *Vector3) SquaredDistanceTo(vec *Vector3) float64 {
	return v.Minus(vec).SquaredLen()
}

func (v *Vector3) MaxValue() float64 {
	return math.Max(v.X, math.Max(v.Y, v.Z))
}

func (v *Vector3) Minus(vec *Vector3) *Vector3 {
	return &Vector3{X: v.X - vec.X, Y: v.Y - vec.Y, Z: v.Z - vec.Z}
}

func (v *Vector3) Plus(vec *Vector3) *Vector3 {
	return &Vector3{X: v.X + vec.X, Y: v.Y + vec.Y, Z: v.Z + vec.Z}
}

func (v *Vector3) Div(vec *Vector3) *Vector3 {
	return &Vector3{X: v.X / vec.X, Y: v.Y / vec.Y, Z: v.Z / vec.Z}
}

func (v *Vector3) PlusValue(val float64) *Vector3 {
	return &Vector3{X: v.X + val, Y: v.X + val, Z: v.X + val}
}

func (v *Vector3) ToString() string {
	return fmt.Sprintf("{%f, %f, %f}", v.X, v.Y, v.Z)
}
