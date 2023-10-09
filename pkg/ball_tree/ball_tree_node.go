package balltree

import "github.com/KrishanBhalla/space-partitioning-trees/pkg/common"

type BallTreeNode[T any] struct {
	point  common.Point
	data   T
	Radius float64
}

var _node common.Node[float64] = BallTreeNode[float64]{}

func (node BallTreeNode[T]) Point() common.Point {
	return node.point
}

func (node BallTreeNode[T]) Data() T {
	return node.data
}

// Triangle inequality - query the children only if the distance between query points minus
// the radius is less than the search distance
func (node BallTreeNode[T]) QueryChildren(point common.Point, distance float64) bool {
	d, _ := common.Distance(point, node.point)
	return d-node.Radius <= distance
}
