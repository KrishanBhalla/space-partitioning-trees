package kdtree

import "github.com/KrishanBhalla/space-partitioning-trees/pkg/common"

type KdTreeNode struct {
	Vector common.PointVector `json:"Vector"`
	Data   common.Point       `json:"Data"`
	// The index of the ordinate on which this node is split
	OrdinateIndex  int     `json:"OrdinateIndex"`
	SplittingValue float64 `json:"SplittingValue"`
}

func (node KdTreeNode) Node() common.Point {
	return node.Data
}

func (node KdTreeNode) SearchLeft(point common.PointVector, distance float64) bool {
	return point[node.OrdinateIndex]-distance <= node.SplittingValue
}

func (node KdTreeNode) SearchRight(point common.PointVector, distance float64) bool {
	return point[node.OrdinateIndex]+distance >= node.SplittingValue
}
