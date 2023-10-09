package balltree

import "github.com/KrishanBhalla/space-partitioning-trees/pkg/common"

type BallTreeNode struct {
	Centroid common.PointVector `json:"Centroid"`
	Data     common.Point       `json:"Data"`
	Radius   float64            `json:"Radius"`
}

// Triangle inequality - query the children only if the distance between query points minus
// the radius is less than the search distance
func (node BallTreeNode) SearchChildren(point common.PointVector, distance float64) bool {
	d, _ := common.Distance(point, node.Centroid)
	return d-node.Radius <= distance
}
