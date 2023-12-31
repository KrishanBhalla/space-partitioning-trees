package balltree

import (
	"fmt"
	"log"
	"math"
	"math/rand"

	"github.com/KrishanBhalla/space-partitioning-trees/pkg/common"
)

type BallTree struct {
	Root      *BallTreeNode `json:"root"`
	Left      *BallTree     `json:"left"`
	Right     *BallTree     `json:"right"`
	Dimension int           `json:"dimension"`
}

var _tree common.SpacePartitioningTree = &BallTree{}

func (tree *BallTree) Construct(points []common.Point, dimension int) error {
	points = common.Filter(points, func(p common.Point) bool {
		return p.Dimension() == dimension
	})
	tree.Dimension = dimension
	err := tree.recursivelyConstruct(points)
	if err != nil {
		return err
	}
	return nil
}

func (tree *BallTree) recursivelyConstruct(points []common.Point) error {
	if len(points) == 0 {
		return nil
	}
	midPoint, orderingAxis := tree.bouncingBallAxis(points)
	radius := findRadiusOfBall(points, midPoint)
	pivot, smaller, larger, err := common.FindMedianByOrdering(orderingAxis, points)
	if err != nil {
		return err
	}
	tree.Root = &BallTreeNode{Centroid: midPoint, Data: pivot, Radius: radius}
	if len(smaller) > 0 {
		tree.Left = &BallTree{Dimension: tree.Dimension}
		tree.Left.recursivelyConstruct(smaller)
	}
	if len(larger) > 0 {
		tree.Right = &BallTree{Dimension: tree.Dimension}
		tree.Right.recursivelyConstruct(larger)
	}
	return nil
}

func findRadiusOfBall(points []common.Point, midPoint common.PointVector) float64 {
	return common.Reduce[common.Point, float64](points, math.Inf(-1), func(r float64, p common.Point) float64 {
		d, _ := common.Distance(p.Vector(), midPoint)
		return math.Max(r, d)
	})
}

// Approximate axis of maximal variation
func (tree *BallTree) bouncingBallAxis(points []common.Point) (common.PointVector, []float64) {
	if len(points) == 0 {
		return nil, nil
	}
	vectors := make([]common.PointVector, len(points))
	for i, p := range points {
		vectors[i] = p.Vector()
	}
	start := vectors[rand.Intn(len(vectors))]
	axisStart := furthestPoint(start, vectors)
	axisEnd := furthestPoint(axisStart, vectors)
	axis, _ := common.Difference(axisStart, axisEnd)

	midPoint := make(common.PointVector, len(axisStart))
	for i, v := range axisStart {
		midPoint[i] = v + axis[i]/2
	}

	dotProduct := common.Map(vectors, func(v common.PointVector) float64 {
		dotProduct, _ := common.DotProduct(v, axis)
		return dotProduct
	})
	return midPoint, dotProduct
}

func furthestPoint(startVec common.PointVector, vecs []common.PointVector) common.PointVector {
	d := 0.
	result := startVec
	for _, v := range vecs {
		new_d, err := common.Distance(v, startVec)
		if err != nil {
			log.Fatal(err)
		}
		if new_d > d {
			result = v
			d = new_d
		}
	}
	return result
}

func (tree BallTree) Search(point common.Point, distance float64) ([]common.Point, error) {
	if point.Dimension() != tree.Dimension {
		return nil, fmt.Errorf("The query point has dimension %d, but the nodes of the tree are of dimension %d", point.Dimension(), tree.Dimension)
	}
	queryStack := []*BallTree{}
	result := []common.Point{}
	pointVector := point.Vector()
	currentNode := &tree
	for currentNode != nil || len(queryStack) > 0 {
		if currentNode != nil {
			queryStack = append(queryStack, currentNode)
			if currentNode.Left != nil && currentNode.Left.Root.SearchChildren(pointVector, distance) {
				currentNode = currentNode.Left
			} else {
				currentNode = nil
			}
		} else {
			currentNode, queryStack = queryStack[len(queryStack)-1], queryStack[:len(queryStack)-1]
			// At this point you must use the vector associated with the data, not with the centroid of the ball
			d, err := common.Distance(pointVector, currentNode.Root.Data.Vector())
			if err != nil {
				return nil, err
			}
			if d < distance {
				result = append(result, currentNode.Root.Data)
			}
			if currentNode.Right != nil && currentNode.Right.Root.SearchChildren(pointVector, distance) {
				currentNode = currentNode.Right
			} else {
				currentNode = nil
			}
		}
	}
	return result, nil
}

func (tree BallTree) KNearestNeighbors(point common.Point, k int) ([]common.Point, error) {

	if point.Dimension() != tree.Dimension {
		return nil, fmt.Errorf("The query point has dimension %d, but the nodes of the tree are of dimension %d", point.Dimension(), tree.Dimension)
	}
	pointVector := point.Vector()
	var candidateNeighbours []common.Point
	var newNode *BallTree
	currentNode := &tree
	for currentNode != nil {
		leftDistance := math.Inf(1)
		rightDistance := math.Inf(1)
		if currentNode.Left != nil {
			leftDistance, _ = common.Distance(currentNode.Left.Root.Centroid, pointVector)
		}
		if currentNode.Right != nil {
			rightDistance, _ = common.Distance(currentNode.Right.Root.Centroid, pointVector)
		}
		if leftDistance < rightDistance {
			newNode = currentNode.Left
		} else {
			newNode = currentNode.Right
		}
		if newNode.Size() < k {
			candidateNeighbours = currentNode.Points()
			break
		} else {
			currentNode = newNode
		}
	}
	distances := common.Map(candidateNeighbours, func(candidate common.Point) float64 {
		d, _ := common.Distance(candidate.Vector(), pointVector)
		return d
	})
	candidateDistance := common.Reduce(distances, 0., math.Max)

	candidateNeighbours, err := tree.Search(point, candidateDistance)
	if err != nil {
		return nil, err
	}
	var candidatesWithDistances common.PointWithDistanceHeap
	candidatesWithDistances = common.Map(candidateNeighbours, func(candidate common.Point) common.PointWithDistance {
		d, _ := common.Distance(candidate.Vector(), pointVector)
		return common.PointWithDistance{Point: candidate, Distance: d}
	})

	for len(candidatesWithDistances) > k {
		candidatesWithDistances.Pop()
	}

	result := common.Map(candidatesWithDistances, func(c common.PointWithDistance) common.Point {
		return c.Point
	})
	return result, nil
}

func (tree BallTree) NodeDimension() int {
	return tree.Dimension
}

func (tree BallTree) Size() int {
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

func (tree BallTree) Depth() int {
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

func (tree BallTree) Points() []common.Point {
	if tree.Root == nil {
		return []common.Point{}
	}
	result := []common.Point{tree.Root.Data}
	if tree.Left != nil {
		result = append(tree.Left.Points(), result...)
	}
	if tree.Right != nil {
		result = append(tree.Right.Points(), result...)
	}
	return result
}
