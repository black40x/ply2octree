package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
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
	reader := octree.NewPlyReader(1)
	err := reader.ReadPly(c.plyFile)

	if octree.DiagonalFraction != 0 {
		octree.Spacing = reader.GetAABB().Size.Len() / octree.DiagonalFraction
	} else if octree.Spacing == 0 {
		octree.DiagonalFraction = 200
	}

	if err != nil {
		return err
	} else {
		PrintInfo(fmt.Sprintf("Ply file reded with %d points", len(reader.GetPoints())))

		c.root = octree.NewRootNode(reader.GetAABB())
		tightAABB := octree.NewEmptyAABB()
		for i, point := range reader.GetPoints() {
			acceptedBy := c.root.Add(point)
			if acceptedBy != nil {
				tightAABB.Update(point.Pos)
			}
			if i%1000 == 0 {
				PrintInfo(fmt.Sprintf("Processed: %d points", i))
			}
		}

		err := c.flush(tightAABB)

		if err != nil {
			return err
		}
	}
	PrintInfo("Convert finished!")
	return nil
}
