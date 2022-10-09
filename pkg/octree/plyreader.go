package octree

import (
	"errors"
	"io"
	"os"
	"regexp"
	"strings"
	"unsafe"
)

type Point struct {
	Pos     *Vector3
	R, G, B byte
}

type PlyProperty struct {
	Name string
	Size int
}

type PlyProperties struct {
	Properties []PlyProperty
}

func (p *PlyProperties) pointByteSize() int {
	size := 0
	for _, prop := range p.Properties {
		size += prop.Size
	}

	return size
}

func (p *PlyProperties) addProperty(name, ptype string) {
	prop := PlyProperty{
		Name: name,
	}

	switch ptype {
	case "char":
		prop.Size = 1
	case "uchar":
		prop.Size = 1
	case "short":
		prop.Size = 2
	case "ushort":
		prop.Size = 2
	case "int":
		prop.Size = 4
	case "uint":
		prop.Size = 4
	case "float":
		prop.Size = 4
	case "double":
		prop.Size = 8
	}

	p.Properties = append(p.Properties, prop)
}

type PlyReader struct {
	points []*Point
	aabb   *AABB
	scale  float64
}

func NewPlyReader(scale float64) *PlyReader {
	return &PlyReader{
		scale: scale,
		aabb:  NewEmptyAABB(),
	}
}

func (r *PlyReader) GetAABB() *AABB {
	return r.aabb
}

func (r *PlyReader) GetPoints() []*Point {
	return r.points
}

func (r *PlyReader) ReadPly(fn string) error {
	f, err := os.Open(fn)
	defer f.Close()

	if err != nil {
		return errors.New("failed file open")
	}

	// Read header

	header := ""
	headerEnd := false
	var headerOffset int64 = 0
	buf := make([]byte, 100)

	for headerEnd != true {
		n, err := f.Read(buf)
		if err != nil {
			return errors.New("failed read header")
		}

		header = header + string(buf[:n])
		if pos := strings.Index(header, "end_header"); pos != -1 {
			headerEnd = true
			headerOffset = int64(pos + len("end_header") + 1)
			header = header[:headerOffset]
		}
	}

	props := PlyProperties{}

	headerProps := strings.Split(header, "\n")
	re, _ := regexp.Compile("^property (char|uchar|short|ushort|int|uint|float|double) (\\w*)")
	for _, p := range headerProps {
		pp := re.FindAllStringSubmatch(p, -1)
		if len(pp) > 0 && len(pp[0]) >= 3 {
			props.addProperty(pp[0][2], pp[0][1])
		}
	}

	if len(props.Properties) == 0 {
		return errors.New("failed read header properties")
	}

	// Read points
	buf = make([]byte, props.pointByteSize())
	_, err = f.Seek(headerOffset, 0)
	if err != nil {
		return errors.New("failed read body")
	}

	for {
		_, err = f.Read(buf)

		if err == io.EOF {
			break
		}

		offset := 0
		point := &Point{
			Pos: &Vector3{},
		}

		for i := 0; i < len(props.Properties); i++ {
			prop := props.Properties[i]
			if prop.Name == "x" {
				memcpy(buf[offset:offset+prop.Size], unsafe.Pointer(&point.Pos.X))
			}
			if prop.Name == "y" {
				memcpy(buf[offset:offset+prop.Size], unsafe.Pointer(&point.Pos.Y))
			}
			if prop.Name == "z" {
				memcpy(buf[offset:offset+prop.Size], unsafe.Pointer(&point.Pos.Z))
			}
			if prop.Name == "red" {
				memcpy(buf[offset:offset+prop.Size], unsafe.Pointer(&point.R))
			}
			if prop.Name == "green" {
				memcpy(buf[offset:offset+prop.Size], unsafe.Pointer(&point.G))
			}
			if prop.Name == "blue" {
				memcpy(buf[offset:offset+prop.Size], unsafe.Pointer(&point.B))
			}

			offset += prop.Size
		}

		r.aabb.Update(point.Pos)
		r.points = append(r.points, point)
	}

	return nil
}
