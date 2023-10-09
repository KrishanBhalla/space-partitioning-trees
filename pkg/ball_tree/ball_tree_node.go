package balltree

import "github.com/KrishanBhalla/space-partitioning-trees/pkg/common"

type BallTreeNode[T any] struct {
	Point  common.Point
	Data   T
	Radius float64
}

var _node common.INode[float64] = BallTreeNode[float64]{}

func (node BallTreeNode[T]) Node() common.Node[T] {
	return common.Node[T]{Point: node.Point, Data: node.Data}
}

// Triangle inequality - query the children only if the distance between query points minus
// the radius is less than the search distance
func (node BallTreeNode[T]) SearchChildren(point common.Point, distance float64) bool {
	d, _ := common.Distance(point, node.Point)
	return d-node.Radius <= distance
}
