package octree

import (
	"os"
	"unsafe"
)

type PointWriter struct {
	File      string
	AABB      *AABB
	Scale     float64
	NumPoints int
}

func (w *PointWriter) Write(p *Point) {
	f, _ := os.OpenFile(w.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	var pos = [3]float64{p.Pos.X, p.Pos.Y, p.Pos.Z}
	bits := unsafe.Slice((*byte)(unsafe.Pointer(&pos[0])), len(pos)*8)

	var colors = [3]byte{p.R, p.G, p.B}
	bitsColors := unsafe.Slice((*byte)(unsafe.Pointer(&colors[0])), len(colors))

	f.Write(bits)
	f.Write(bitsColors)

	// f2, _ := os.OpenFile(w.File+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// defer f2.Close()

	// f2.WriteString(fmt.Sprintf("x: %f, y: %f, z: %f, r: %d, g: %d, b: %d \n", p.Pos.X, p.Pos.Y, p.Pos.Z, p.R, p.G, p.B))
}
