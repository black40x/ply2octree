package octree

import "math"

type AABB struct {
	Min  *Vector3
	Max  *Vector3
	Size *Vector3
}

func NewEmptyAABB() *AABB {
	return &AABB{
		Max:  NewVector3(-math.MaxFloat64),
		Min:  NewVector3(math.MaxFloat64),
		Size: NewVector3(math.MaxFloat64),
	}
}

func NewAABB(max, min *Vector3) *AABB {
	return &AABB{
		Max:  max,
		Min:  min,
		Size: max.Minus(min),
	}
}

func (a *AABB) IsInside(v *Vector3) bool {
	if a.Min.X <= v.X && v.X <= a.Max.X {
		if a.Min.Y <= v.Y && v.Y <= a.Max.Y {
			if a.Min.Z <= v.Z && v.Z <= a.Max.Z {
				return true
			}
		}
	}

	return false
}

func (a *AABB) Update(v *Vector3) {
	a.Min.X = math.Min(a.Min.X, v.X)
	a.Min.Y = math.Min(a.Min.Y, v.Y)
	a.Min.Z = math.Min(a.Min.Z, v.Z)

	a.Max.X = math.Max(a.Max.X, v.X)
	a.Max.Y = math.Max(a.Max.Y, v.Y)
	a.Max.Z = math.Max(a.Max.Z, v.Z)

	a.Size = a.Max.Minus(a.Min)
}

func (a *AABB) UpdateAABB(ab *AABB) {
	a.Update(ab.Min)
	a.Update(ab.Max)
}

func (a *AABB) MakeCubic() {
	a.Max = a.Min.PlusValue(a.Size.MaxValue())
	a.Size = a.Max.Minus(a.Min)
}
