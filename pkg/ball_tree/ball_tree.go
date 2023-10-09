package balltree

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/KrishanBhalla/space-partitioning-trees/pkg/common"
)

type BallTree struct {
	Root      *BallTreeNode `json:"root"`
	Left      *BallTree     `json:"left"`
	Right     *BallTree     `json:"right"`
	Dimension int           `json:"dimension"`
}

var _tree common.SpacePartitioningTree = BallTree{}

func (tree BallTree) Construct(points []common.Point, dimension int) error {
	points = common.Filter(points, func(p common.Point) bool {
		return p.Dimension() == dimension
	})
	tree.Dimension = dimension
	err := tree.recursivelyConstruct(points)
	if err != nil {
		return err
	}
	return nil
}

func (tree *BallTree) recursivelyConstruct(points []common.Point) error {
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
	tree.Root = &BallTreeNode{Vector: pivot.Vector(), Data: pivot, Radius: radius}
	if len(smaller) > 0 {
		tree.Left = &BallTree{Dimension: tree.Dimension}
		tree.Left.recursivelyConstruct(smaller)
	}
	if len(larger) > 0 {
		tree.Right = &BallTree{Dimension: tree.Dimension}
		tree.Right.recursivelyConstruct(larger)
	}
	return nil
}

type pointWithDistance struct {
	distance float64
	point    common.PointVector
}

// Approximate axis of maximal variation
func (tree *BallTree) bouncingBallAxis(points []common.Point) []float64 {
	if len(points) == 0 {
		return nil
	}
	vectors := make([]common.PointVector, len(points))
	for i, p := range points {
		vectors[i] = p.Vector()
	}
	start := vectors[rand.Int()%len(vectors)]
	axisStart := common.Reduce(vectors, pointWithDistance{0, start}, furthestPoint).point
	axisEnd := common.Reduce(vectors, pointWithDistance{0, axisStart}, furthestPoint).point
	axis := make(common.PointVector, len(axisStart))
	for i, v := range axisEnd {
		axis[i] = v - axisStart[i]
	}
	return common.Map(vectors, func(v common.PointVector) float64 {
		dotProduct, _ := common.DotProduct(v, axis)
		return dotProduct
	})
}

func furthestPoint(r pointWithDistance, p common.PointVector) pointWithDistance {
	d, _ := common.Distance(p, r.point)
	if d > r.distance {
		return pointWithDistance{distance: d, point: p}
	}
	return r
}

func (tree BallTree) Search(point common.PointVector, distance float64) ([]common.Point, error) {
	if len(point) != tree.Dimension {
		return nil, fmt.Errorf("The query point has Dimension %d, but the nodes of the tree are of Dimension %d", len(point), tree.Dimension)
	}
	queryStack := []*BallTree{}
	result := []common.Point{}
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
			d, err := common.Distance(point, currentNode.Root.Vector)
			if err != nil {
				return nil, err
			}
			if d < distance {
				result = append(result, currentNode.Root.Data)
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

// func (tree *BallTree) KNearestNeighbors(point common.Point, k int) ([]*BallTreeNode, error) {
// 	return nil, nil
// }

func (tree BallTree) NodeDimension() int {
	return tree.Dimension
}

func (tree BallTree) Size() int {
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

func (tree BallTree) Depth() int {
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
