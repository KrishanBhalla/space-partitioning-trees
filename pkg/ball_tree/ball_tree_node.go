package balltree

import "github.com/KrishanBhalla/space-partitioning-trees/pkg/common"

type BallTreeNode struct {
	Point  common.PointVector `json:"Vector"`
	Data   common.Point       `json:"Data"`
	Radius float64            `json:"Radius"`
}

// Triangle inequality - query the children only if the distance between query points minus
// the radius is less than the search distance
func (node BallTreeNode) SearchChildren(point common.Point, distance float64) bool {
	d, _ := common.Distance(point, node.Point)
	return d-node.Radius <= distance
}
