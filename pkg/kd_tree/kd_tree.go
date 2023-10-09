package kdtree

import (
	"fmt"

	"github.com/KrishanBhalla/space-partitioning-trees/pkg/common"
)

type KdTree[T any] struct {
	root      *KdTreeNode[T]
	left      *KdTree[T]
	right     *KdTree[T]
	dimension int
}

var _tree common.SpacePartitioningTree[float64, KdTreeNode[float64]] = KdTree[float64]{}

func (tree KdTree[T]) Insert(point *common.PointWithData[T]) error {
	for tree.root != nil {
		if tree.left != nil && tree.root.SearchLeft(point.Point, 0) {
			tree = *tree.left
		} else if tree.right != nil && tree.root.SearchRight(point.Point, 0) {
			tree = *tree.right
		}
	}
	ordinateIndex := (tree.root.OrdinateIndex + 1) % tree.dimension
	node := &KdTreeNode[T]{point: point.Point, data: point.Data, OrdinateIndex: ordinateIndex, SplittingValue: point.Point[ordinateIndex]}
	newTree := &KdTree[T]{root: node, dimension: tree.dimension}
	if tree.left == nil {
		tree.left = newTree
	} else if tree.right == nil {
		tree.right = newTree
	}
	return nil
}

func (tree KdTree[T]) Construct(points []*common.PointWithData[T], dimension int) error {
	points = common.Filter(points, func(p *common.PointWithData[T]) bool {
		return len(p.Point) == dimension
	})
	tree.dimension = dimension
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
	tree.root = &KdTreeNode[T]{point: pivot.Point, data: pivot.Data, OrdinateIndex: ordinateIndex, SplittingValue: pivot.Point[ordinateIndex]}
	if len(smaller) > 0 {
		tree.left = &KdTree[T]{dimension: tree.dimension}
		tree.left.recursivelyConstruct(smaller, (ordinateIndex+1)%tree.dimension)
	}
	if len(larger) > 0 {
		tree.right = &KdTree[T]{dimension: tree.dimension}
		tree.right.recursivelyConstruct(larger, (ordinateIndex+1)%tree.dimension)
	}
	return nil
}

func (tree KdTree[T]) Search(point common.Point, distance float64) ([]*KdTreeNode[T], error) {
	if len(point) != tree.dimension {
		return nil, fmt.Errorf("The query point has dimension %d, but the nodes of the tree are of dimension %d", len(point), tree.dimension)
	}
	queryStack := []*KdTree[T]{}
	result := []*KdTreeNode[T]{}
	currentNode := &tree
	for currentNode != nil || len(queryStack) > 0 {
		if currentNode != nil {
			queryStack = append(queryStack, currentNode)
			if currentNode.left != nil && currentNode.root.SearchLeft(point, distance) {
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
			if currentNode.root.SearchRight(point, distance) {
				currentNode = currentNode.right
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

func (tree KdTree[T]) Dimension() int {
	return tree.dimension
}

func (tree KdTree[T]) Size() int {
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

func (tree KdTree[T]) Depth() int {
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
