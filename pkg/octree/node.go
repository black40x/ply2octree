package octree

import (
	"fmt"
	"math"
	"path"
)

var maxDepth = -1
var maxStoreSize = 10000

var Spacing = 0.0
var DiagonalFraction = 200.0

type Node struct {
	Index                   int
	Level                   int
	AABB                    *AABB
	AcceptedAABB            *AABB
	Parent                  *Node
	Children                map[int]*Node
	Grid                    *Grid
	Store                   []*Point
	Cache                   []*Point
	HasChildren             bool
	AddedSinceLastFlush     bool
	AddCalledSinceLastFlush bool
	NumAccepted             int
}

func (n *Node) Name() string {
	if n.Parent == nil {
		return "r"
	} else {
		return fmt.Sprintf("%s%d", n.Parent.Name(), n.Index)
	}
}

func NewRootNode(aabb *AABB) *Node {
	root := &Node{
		Index:                   -1,
		AABB:                    aabb,
		Level:                   0,
		Parent:                  nil,
		Children:                make(map[int]*Node),
		HasChildren:             false,
		AddedSinceLastFlush:     true,
		AddCalledSinceLastFlush: false,
		NumAccepted:             0,
	}

	root.Grid = NewGrid(aabb, root.Spacing())
	root.AcceptedAABB = NewEmptyAABB()

	return root
}

func (n *Node) CreateChild(index int) *Node {
	aabb := ChildAABB(n.AABB, index)
	child := &Node{
		Index:                   index,
		AABB:                    aabb,
		Level:                   n.Level + 1,
		Parent:                  n,
		Children:                make(map[int]*Node),
		HasChildren:             false,
		AddedSinceLastFlush:     true,
		AddCalledSinceLastFlush: false,
		NumAccepted:             0,
	}

	if child.Level%2 == 0 {
		child.Index = 7 - child.Index
	}

	child.Grid = NewGrid(aabb, child.Spacing())
	child.AcceptedAABB = NewEmptyAABB()

	n.Children[index] = child

	return child
}

func (n *Node) IsLeafNode() bool {
	return !n.HasChildren
}

func (n *Node) Spacing() float64 {
	return Spacing / math.Pow(2.0, float64(n.Level))
}

func (n *Node) PointsCount() int {
	if n.NumAccepted == 0 {
		return len(n.Store)
	} else {
		return n.NumAccepted
	}
}

func (n *Node) Add(point *Point) *Node {
	n.AddCalledSinceLastFlush = true

	if n.IsLeafNode() {
		n.Store = append(n.Store, point)
		if len(n.Store) >= maxStoreSize {
			n.Split()
		}

		return n
	} else {
		var accepted = false
		accepted = n.Grid.Add(point)

		if accepted {
			n.Cache = append(n.Cache, point)
			n.AcceptedAABB.Update(point.Pos)
			n.NumAccepted++
			return n
		} else {
			if maxDepth != -1 && n.Level >= maxDepth {
				return nil
			}

			childIndex := NodeIndex(*n.AABB, *point)

			if childIndex >= 0 {
				if n.IsLeafNode() {
					n.HasChildren = true
				}

				child := n.Children[childIndex]

				if child == nil {
					child = n.CreateChild(childIndex)
					n.Children[childIndex] = child
				}

				return child.Add(point)
			} else {
				return nil
			}
		}
	}

	return nil
}

func (n *Node) Split() {
	n.HasChildren = true

	for _, point := range n.Store {
		n.Add(point)
	}

	n.Store = make([]*Point, 0)
}

func (n *Node) GetHierarchy(levels int) []*Node {
	var hierarchy []*Node
	var stack []*Node
	stack = append(stack, n)

	for len(stack) != 0 {
		node := stack[0]
		stack = stack[1:]

		if node.Level >= n.Level+levels {
			break
		}

		hierarchy = append(hierarchy, node)

		for _, child := range node.Children {
			if child != nil {
				stack = append(stack, child)
			}
		}
	}

	return hierarchy
}

func (n *Node) Flush(out string) {
	writeToDisk := func(points []*Point) {
		writer := &PointWriter{
			File:      path.Join(out, n.Name()+".bin"),
			AABB:      n.AABB,
			Scale:     0,
			NumPoints: 0,
		}

		for _, point := range points {
			writer.Write(point)
		}
	}

	if n.IsLeafNode() {
		if n.AddCalledSinceLastFlush {
			writeToDisk(n.Store)
		} else {
			n.Store = make([]*Point, 0)
		}
	} else {
		if n.AddCalledSinceLastFlush {
			writeToDisk(n.Cache)
			n.Cache = make([]*Point, 0)
		}
	}

	n.AddCalledSinceLastFlush = false

	for _, child := range n.Children {
		if child != nil {
			child.Flush(out)
		}
	}
}
