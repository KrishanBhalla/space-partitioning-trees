package kdtree

import (
	"fmt"

	"github.com/KrishanBhalla/space-partitioning-trees/pkg/common"
)

type KdTree struct {
	Root      *KdTreeNode `json:"root"`
	Left      *KdTree     `json:"left"`
	Right     *KdTree     `json:"right"`
	Dimension int         `json:"dimension"`
}

var _tree common.SpacePartitioningTree = KdTree{}

// func (tree KdTree) Insert(point *common.Point) error {
// 	for tree.Root != nil {
// 		if tree.Left != nil && tree.Root.SearchLeft(point.Point, 0) {
// 			tree = *tree.Left
// 		} else if tree.Right != nil && tree.Root.SearchRight(point.Point, 0) {
// 			tree = *tree.Right
// 		}
// 	}
// 	ordinateIndex := (tree.Root.OrdinateIndex + 1) % tree.Dimension
// 	node := &KdTreeNode{point: point.Point, data: point.Data, OrdinateIndex: ordinateIndex, SplittingValue: point.Point[ordinateIndex]}
// 	newTree := &KdTree{Root: node, Dimension: tree.Dimension}
// 	if tree.Left == nil {
// 		tree.Left = newTree
// 	} else if tree.Right == nil {
// 		tree.Right = newTree
// 	}
// 	return nil
// }

func (tree KdTree) Construct(points []common.Point, dimension int) error {
	points = common.Filter(points, func(p common.Point) bool {
		return p.Dimension() == dimension
	})
	tree.Dimension = dimension
	ordinateIndex := 0
	err := tree.recursivelyConstruct(points, ordinateIndex)
	if err != nil {
		return err
	}
	return nil
}

func (tree *KdTree) recursivelyConstruct(points []common.Point, ordinateIndex int) error {
	if len(points) == 0 {
		return nil
	}
	ordinateValues := common.Map(points, func(p common.Point) float64 { return p.Vector()[ordinateIndex] })
	pivot, smaller, larger, err := common.FindMedianByOrdering(ordinateValues, points)
	if err != nil {
		return err
	}
	tree.Root = &KdTreeNode{Vector: pivot.Vector(), Data: pivot, OrdinateIndex: ordinateIndex, SplittingValue: pivot.Vector()[ordinateIndex]}
	if len(smaller) > 0 {
		tree.Left = &KdTree{Dimension: tree.Dimension}
		tree.Left.recursivelyConstruct(smaller, (ordinateIndex+1)%tree.Dimension)
	}
	if len(larger) > 0 {
		tree.Right = &KdTree{Dimension: tree.Dimension}
		tree.Right.recursivelyConstruct(larger, (ordinateIndex+1)%tree.Dimension)
	}
	return nil
}

func (tree KdTree) Search(point common.PointVector, distance float64) ([]common.Point, error) {
	if len(point) != tree.Dimension {
		return nil, fmt.Errorf("The query point has dimension %d, but the nodes of the tree are of dimension %d", len(point), tree.Dimension)
	}
	queryStack := []*KdTree{}
	result := []common.Point{}
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
			d, err := common.Distance(point, currentNode.Root.Vector)
			if err != nil {
				return nil, err
			}
			if d < distance {
				result = append(result, currentNode.Root.Data)
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

// func (tree *KdTree) KNearestNeighbors(point common.Point, k int) ([]*KdTreeNode, error) {
// 	return nil, nil
// }

func (tree KdTree) NodeDimension() int {
	return tree.Dimension
}

func (tree KdTree) Size() int {
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

func (tree KdTree) Depth() int {
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
