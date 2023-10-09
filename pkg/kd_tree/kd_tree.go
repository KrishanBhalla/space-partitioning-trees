package kdtree

import (
	"fmt"

	"github.com/KrishanBhalla/space-partitioning-trees/pkg/common"
)

type KdTree[T any] struct {
	Root      *KdTreeNode[T] `json:"root"`
	Left      *KdTree[T]     `json:"left"`
	Right     *KdTree[T]     `json:"right"`
	Dimension int            `json:"dimension"`
}

var _tree common.SpacePartitioningTree[float64, KdTreeNode[float64]] = KdTree[float64]{}

// func (tree KdTree[T]) Insert(point *common.PointWithData[T]) error {
// 	for tree.Root != nil {
// 		if tree.Left != nil && tree.Root.SearchLeft(point.Point, 0) {
// 			tree = *tree.Left
// 		} else if tree.Right != nil && tree.Root.SearchRight(point.Point, 0) {
// 			tree = *tree.Right
// 		}
// 	}
// 	ordinateIndex := (tree.Root.OrdinateIndex + 1) % tree.Dimension
// 	node := &KdTreeNode[T]{point: point.Point, data: point.Data, OrdinateIndex: ordinateIndex, SplittingValue: point.Point[ordinateIndex]}
// 	newTree := &KdTree[T]{Root: node, Dimension: tree.Dimension}
// 	if tree.Left == nil {
// 		tree.Left = newTree
// 	} else if tree.Right == nil {
// 		tree.Right = newTree
// 	}
// 	return nil
// }

func (tree KdTree[T]) Construct(points []*common.PointWithData[T], Dimension int) error {
	points = common.Filter(points, func(p *common.PointWithData[T]) bool {
		return len(p.Point) == Dimension
	})
	tree.Dimension = Dimension
	ordinateIndex := 0
	err := tree.recursivelyConstruct(points, ordinateIndex)
	if err != nil {
		return err
	}
	return nil
}

func (tree *KdTree[T]) recursivelyConstruct(points []*common.PointWithData[T], ordinateIndex int) error {
	if len(points) == 0 {
		return nil
	}
	ordinateValues := common.Map(points, func(p *common.PointWithData[T]) float64 { return p.Point[ordinateIndex] })
	pivot, smaller, larger, err := common.FindMedianByOrdering(ordinateValues, points)
	if err != nil {
		return err
	}
	tree.Root = &KdTreeNode[T]{Point: pivot.Point, Data: pivot.Data, OrdinateIndex: ordinateIndex, SplittingValue: pivot.Point[ordinateIndex]}
	if len(smaller) > 0 {
		tree.Left = &KdTree[T]{Dimension: tree.Dimension}
		tree.Left.recursivelyConstruct(smaller, (ordinateIndex+1)%tree.Dimension)
	}
	if len(larger) > 0 {
		tree.Right = &KdTree[T]{Dimension: tree.Dimension}
		tree.Right.recursivelyConstruct(larger, (ordinateIndex+1)%tree.Dimension)
	}
	return nil
}

func (tree KdTree[T]) Search(point common.Point, distance float64) ([]*KdTreeNode[T], error) {
	if len(point) != tree.Dimension {
		return nil, fmt.Errorf("The query point has Dimension %d, but the nodes of the tree are of Dimension %d", len(point), tree.Dimension)
	}
	queryStack := []*KdTree[T]{}
	result := []*KdTreeNode[T]{}
	currentNode := &tree
	for currentNode != nil || len(queryStack) > 0 {
		if currentNode != nil {
			queryStack = append(queryStack, currentNode)
			if currentNode.Left != nil && currentNode.Root.SearchLeft(point, distance) {
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
			if currentNode.Root.SearchRight(point, distance) {
				currentNode = currentNode.Right
			} else {
				currentNode = nil
			}
		}
	}
	return result, nil
}

// func (tree *KdTree[T]) KNearestNeighbors(point common.Point, k int) ([]*KdTreeNode[T], error) {
// 	return nil, nil
// }

func (tree KdTree[T]) NodeDimension() int {
	return tree.Dimension
}

func (tree KdTree[T]) Size() int {
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

func (tree KdTree[T]) Depth() int {
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
