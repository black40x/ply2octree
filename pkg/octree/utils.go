package octree

import "unsafe"

func memcpy(bits []byte, dest unsafe.Pointer) {
	copy(unsafe.Slice((*byte)(unsafe.Pointer(dest)), len(bits)), bits)
}

func NodeIndex(aabb AABB, point Point) int {
	mx := (int)(2.0 * (point.Pos.X - aabb.Min.X) / aabb.Size.X)
	my := (int)(2.0 * (point.Pos.Y - aabb.Min.Y) / aabb.Size.Y)
	mz := (int)(2.0 * (point.Pos.Z - aabb.Min.Z) / aabb.Size.Z)

	mx = Min(mx, 1)
	my = Min(my, 1)
	mz = Min(mz, 1)

	return (mx << 2) | (my << 1) | mz
}

func ChildAABB(aabb *AABB, index int) *AABB {
	min := *aabb.Min
	max := *aabb.Max

	if (index & 0b0001) > 0 {
		min.Z += aabb.Size.Z / 2
	} else {
		max.Z -= aabb.Size.Z / 2
	}

	if (index & 0b0010) > 0 {
		min.Y += aabb.Size.Y / 2
	} else {
		max.Y -= aabb.Size.Y / 2
	}

	if (index & 0b0100) > 0 {
		min.X += aabb.Size.X / 2
	} else {
		max.X -= aabb.Size.X / 2
	}

	return NewAABB(&min, &max)
}

func Max[T float64 | float32 | int](a, b T) T {
	if a > b {
		return a
	}

	return b
}

func Min[T float64 | float32 | int](a, b T) T {
	if a < b {
		return a
	}

	return b
}
