package balltree

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/KrishanBhalla/space-partitioning-trees/pkg/common"
)

type BallTree[T any] struct {
	Root      *BallTreeNode[T] `json:"root"`
	Left      *BallTree[T]     `json:"left"`
	Right     *BallTree[T]     `json:"right"`
	Dimension int              `json:"dimension"`
}

var _tree common.SpacePartitioningTree[float64, BallTreeNode[float64]] = BallTree[float64]{}

func (tree BallTree[T]) Construct(points []*common.PointWithData[T], Dimension int) error {
	points = common.Filter(points, func(p *common.PointWithData[T]) bool {
		return len(p.Point) == Dimension
	})
	tree.Dimension = Dimension
	err := tree.recursivelyConstruct(points)
	if err != nil {
		return err
	}
	return nil
}

func (tree *BallTree[T]) recursivelyConstruct(points []*common.PointWithData[T]) error {
	if len(points) == 0 {
		return nil
	}
	orderingAxis := tree.bouncingBallAxis(points)
	minValueOnAxis := common.Reduce[float64, float64](orderingAxis, math.Inf(1), math.Min)
	maxValueOnAxis := common.Reduce[float64, float64](orderingAxis, math.Inf(-1), math.Max)
	radius := (maxValueOnAxis - minValueOnAxis) / 2
	pivot, smaller, larger, err := common.FindMedianByOrdering(orderingAxis, points)
	if err != nil {
		return err
	}
	tree.Root = &BallTreeNode[T]{Point: pivot.Point, Data: pivot.Data, Radius: radius}
	if len(smaller) > 0 {
		tree.Left = &BallTree[T]{Dimension: tree.Dimension}
		tree.Left.recursivelyConstruct(smaller)
	}
	if len(larger) > 0 {
		tree.Right = &BallTree[T]{Dimension: tree.Dimension}
		tree.Right.recursivelyConstruct(larger)
	}
	return nil
}

type pointWithDistance struct {
	distance float64
	point    common.Point
}

// Approximate axis of maximal variation
func (tree *BallTree[T]) bouncingBallAxis(pointsWithData []*common.PointWithData[T]) []float64 {
	if len(pointsWithData) == 0 {
		return nil
	}
	points := make([]common.Point, len(pointsWithData))
	for i, p := range pointsWithData {
		points[i] = p.Point
	}
	start := points[rand.Int()%len(points)]
	axisStart := common.Reduce(points, pointWithDistance{0, start}, furthestPoint).point
	axisEnd := common.Reduce(points, pointWithDistance{0, axisStart}, furthestPoint).point
	axis := make(common.Point, len(axisStart))
	for i, v := range axisEnd {
		axis[i] = v - axisStart[i]
	}
	return common.Map(points, func(p common.Point) float64 {
		dotProduct, _ := common.DotProduct(p, axis)
		return dotProduct
	})
}

func furthestPoint(r pointWithDistance, p common.Point) pointWithDistance {
	d, _ := common.Distance(p, r.point)
	if d > r.distance {
		return pointWithDistance{distance: d, point: p}
	}
	return r
}

func (tree BallTree[T]) Search(point common.Point, distance float64) ([]*BallTreeNode[T], error) {
	if len(point) != tree.Dimension {
		return nil, fmt.Errorf("The query point has Dimension %d, but the nodes of the tree are of Dimension %d", len(point), tree.Dimension)
	}
	queryStack := []*BallTree[T]{}
	result := []*BallTreeNode[T]{}
	currentNode := &tree
	for currentNode != nil || len(queryStack) > 0 {
		if currentNode != nil {
			queryStack = append(queryStack, currentNode)
			if currentNode.Left != nil && currentNode.Left.Root.SearchChildren(point, distance) {
				currentNode = currentNode.Left
			} else {
				currentNode = nil
			}
		} else {
			currentNode, queryStack = queryStack[len(queryStack)-1], queryStack[:len(queryStack)-1]
			d, err := common.Distance(point, currentNode.Root.Point)
			if err != nil {
				return nil, err
			}
			if d < distance {
				result = append(result, currentNode.Root)
			}
			if currentNode.Right != nil && currentNode.Right.Root.SearchChildren(point, distance) {
				currentNode = currentNode.Right
			} else {
				currentNode = nil
			}
		}
	}
	return result, nil
}

// func (tree *BallTree[T]) KNearestNeighbors(point common.Point, k int) ([]*BallTreeNode[T], error) {
// 	return nil, nil
// }

func (tree BallTree[T]) NodeDimension() int {
	return tree.Dimension
}

func (tree BallTree[T]) Size() int {
	if tree.Root == nil {
		return 0
	} else if tree.Left == nil && tree.Right == nil {
		return 1
	} else if tree.Left == nil {
		return 1 + tree.Right.Size()
	} else if tree.Right == nil {
		return 1 + tree.Left.Size()
	}
	return 1 + tree.Left.Size() + tree.Right.Size()
}

func (tree BallTree[T]) Depth() int {
	if tree.Root == nil {
		return 0
	} else if tree.Left == nil && tree.Right == nil {
		return 1
	} else if tree.Left == nil {
		return 1 + tree.Right.Depth()
	} else if tree.Right == nil {
		return 1 + tree.Left.Depth()
	}
	return 1 + max(tree.Left.Depth(), tree.Right.Depth())
}
