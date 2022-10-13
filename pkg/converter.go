package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/black40x/plyfile/plyfile"
	"github.com/schollz/progressbar/v3"
	"io"
	"os"
	"path"
	"ply2octree/pkg/octree"
)

type MetaBoundingBox struct {
	Lx, Ly, Lz, Ux, Uy, Uz float64
}

type MetaJson struct {
	Spacing     float64
	BoundingBox MetaBoundingBox
	Hierarchy   []interface{}
}

type Converter struct {
	root              *octree.Node
	hierarchyStepSize int
	plyFile           string
	outputDir         string
}

func NewConverter(plyFile, outputDir string, hierarchyStepSize int) *Converter {
	return &Converter{
		root:              nil,
		plyFile:           plyFile,
		outputDir:         outputDir,
		hierarchyStepSize: hierarchyStepSize,
	}
}

func (c *Converter) flush(aabb *octree.AABB) error {
	os.Mkdir(c.outputDir, 0777)

	if _, err := os.Stat(c.outputDir); os.IsNotExist(err) {
		return errors.New("output directory not exist")
	}

	hrcTotal := 0
	hrcFlushed := 0

	meta := MetaJson{
		Spacing: octree.Spacing,
		BoundingBox: MetaBoundingBox{
			Lx: aabb.Min.X,
			Ly: aabb.Min.Y,
			Lz: aabb.Min.Z,
			Ux: aabb.Max.X,
			Uy: aabb.Max.Y,
			Uz: aabb.Max.Z,
		},
	}

	var stack []*octree.Node
	stack = append(stack, c.root)

	for len(stack) != 0 {
		node := stack[0]
		stack = stack[1:]
		hrcTotal++

		hierarchy := node.GetHierarchy(c.hierarchyStepSize + 1)
		needsFlush := false

		for _, descendant := range hierarchy {
			if descendant.Level == node.Level+c.hierarchyStepSize {
				stack = append(stack, descendant)
			}
			needsFlush = needsFlush || descendant.AddedSinceLastFlush
		}

		if needsFlush {
			for _, descendant := range hierarchy {
				children := 0
				for i := 0; i < len(descendant.Children); i++ {
					if descendant.Children[i] != nil {
						children = children | (1 << i)
					}
				}

				meta.Hierarchy = append(meta.Hierarchy, [2]interface{}{descendant.Name(), descendant.PointsCount()})
			}

			hrcFlushed++
		}
	}

	PrintInfo("Save meta.json")

	metaData, _ := json.MarshalIndent(meta, "", " ")
	os.WriteFile(path.Join(c.outputDir, "meta.json"), metaData, 0666)

	PrintInfo("Save hierarchy binaries")

	c.root.Flush(c.outputDir)

	return nil
}

func (c *Converter) Convert() error {
	PrintInfo("Convert started...")

	ply, err := plyfile.Open(c.plyFile)
	if err != nil {
		return err
	}
	defer ply.Close()

	point := plyfile.Point{}
	aabb := octree.NewEmptyAABB()
	r, err := ply.GetElementReader("vertex")
	if err == nil {
		for {
			_, err := r.ReadNext(&point)
			if err == io.EOF {
				break
			}
			aabb.Update(&octree.Vector3{X: point.X, Y: point.Y, Z: point.Z})
		}
	} else {
		return err
	}

	if octree.DiagonalFraction != 0 {
		octree.Spacing = aabb.Size.Len() / octree.DiagonalFraction
	} else if octree.Spacing == 0 {
		octree.DiagonalFraction = 200
	}

	r.Reset()
	PrintInfo(fmt.Sprintf("Ply file reded with %d points", r.Count()))
	bar := progressbar.Default(r.Count(), "processing")

	c.root = octree.NewRootNode(aabb)
	tightAABB := octree.NewEmptyAABB()

	for {
		n, err := r.ReadNext(&point)
		if err == io.EOF {
			break
		}
		if n%100 == 0 {
			bar.Add(100)
		}
		ocPoint := &octree.Point{
			Pos: &octree.Vector3{
				X: point.X,
				Y: point.Y,
				Z: point.Z,
			},
			R: point.R,
			G: point.G,
			B: point.B,
		}
		acceptedBy := c.root.Add(ocPoint)
		if acceptedBy != nil {
			tightAABB.Update(ocPoint.Pos)
		}
		aabb.Update(&octree.Vector3{X: point.X, Y: point.Y, Z: point.Z})
	}

	bar.Finish()

	err = c.flush(tightAABB)
	if err != nil {
		return err
	}

	PrintInfo("Convert finished!")
	return nil
}
