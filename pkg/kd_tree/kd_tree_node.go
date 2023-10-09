package kdtree

import "github.com/KrishanBhalla/space-partitioning-trees/pkg/common"

type KdTreeNode[T any] struct {
	Point common.Point `json:"Point"`
	Data  T            `json:"Data"`
	// The index of the ordinate on which this node is split
	OrdinateIndex  int     `json:"OrdinateIndex"`
	SplittingValue float64 `json:"SplittingValue"`
}

var _node common.INode[float64] = KdTreeNode[float64]{}

func (node KdTreeNode[T]) Node() common.Node[T] {
	return common.Node[T]{Point: node.Point, Data: node.Data}
}

func (node KdTreeNode[T]) SearchLeft(point common.Point, distance float64) bool {
	return point[node.OrdinateIndex]-distance <= node.SplittingValue
}

func (node KdTreeNode[T]) SearchRight(point common.Point, distance float64) bool {
	return point[node.OrdinateIndex]+distance >= node.SplittingValue
}
