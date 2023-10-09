package balltree

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/KrishanBhalla/space-partitioning-trees/pkg/common"
)

type BallTree[T any] struct {
	root      *BallTreeNode[T]
	left      *BallTree[T]
	right     *BallTree[T]
	dimension int
}

var _tree common.SpacePartitioningTree[float64, BallTreeNode[float64]] = BallTree[float64]{}

func (tree BallTree[T]) Construct(points []*common.PointWithData[T], dimension int) error {
	points = common.Filter(points, func(p *common.PointWithData[T]) bool {
		return len(p.Point) == dimension
	})
	tree.dimension = dimension
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
	tree.root = &BallTreeNode[T]{point: pivot.Point, data: pivot.Data, Radius: radius}
	if len(smaller) > 0 {
		tree.left = &BallTree[T]{dimension: tree.dimension}
		tree.left.recursivelyConstruct(smaller)
	}
	if len(larger) > 0 {
		tree.right = &BallTree[T]{dimension: tree.dimension}
		tree.right.recursivelyConstruct(larger)
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
	if len(point) != tree.dimension {
		return nil, fmt.Errorf("The query point has dimension %d, but the nodes of the tree are of dimension %d", len(point), tree.dimension)
	}
	queryStack := []*BallTree[T]{}
	result := []*BallTreeNode[T]{}
	currentNode := &tree
	for currentNode != nil || len(queryStack) > 0 {
		if currentNode != nil {
			queryStack = append(queryStack, currentNode)
			if currentNode.left != nil && currentNode.left.root.SearchChildren(point, distance) {
				currentNode = currentNode.left
			} else {
				currentNode = nil
			}
		} else {
			currentNode, queryStack = queryStack[len(queryStack)-1], queryStack[:len(queryStack)-1]
			d, err := common.Distance(point, currentNode.root.point)
			if err != nil {
				return nil, err
			}
			if d < distance {
				result = append(result, currentNode.root)
			}
			if currentNode.right != nil && currentNode.right.root.SearchChildren(point, distance) {
				currentNode = currentNode.right
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

func (tree BallTree[T]) Dimension() int {
	return tree.dimension
}

func (tree BallTree[T]) Size() int {
	if tree.root == nil {
		return 0
	} else if tree.left == nil && tree.right == nil {
		return 1
	} else if tree.left == nil {
		return 1 + tree.right.Size()
	} else if tree.right == nil {
		return 1 + tree.left.Size()
	}
	return 1 + tree.left.Size() + tree.right.Size()
}

func (tree BallTree[T]) Depth() int {
	if tree.root == nil {
		return 0
	} else if tree.left == nil && tree.right == nil {
		return 1
	} else if tree.left == nil {
		return 1 + tree.right.Depth()
	} else if tree.right == nil {
		return 1 + tree.left.Depth()
	}
	return 1 + max(tree.left.Depth(), tree.right.Depth())
}
