package kdtree

import "github.com/KrishanBhalla/space-partitioning-trees/pkg/common"

type KdTreeNode[T any] struct {
	point common.Point
	data  T
	// The index of the ordinate on which this node is split
	OrdinateIndex  int
	SplittingValue float64
}

var _node common.Node[float64] = KdTreeNode[float64]{}

func (node KdTreeNode[T]) Point() common.Point {
	return node.point
}

func (node KdTreeNode[T]) Data() T {
	return node.data
}

func (node KdTreeNode[T]) SearchLeft(point common.Point, distance float64) bool {
	return point[node.OrdinateIndex]-distance <= node.SplittingValue
}

func (node KdTreeNode[T]) SearchRight(point common.Point, distance float64) bool {
	return point[node.OrdinateIndex]+distance >= node.SplittingValue
}
