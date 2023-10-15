package balltree_test

import (
	"math"
	"math/rand"
	"testing"

	balltree "github.com/KrishanBhalla/space-partitioning-trees/pkg/ball_tree"
	"github.com/KrishanBhalla/space-partitioning-trees/pkg/common"
	"github.com/stretchr/testify/assert"
)

type testPoint struct {
	dimension int
	vector    common.PointVector
}

func (t *testPoint) Dimension() int {
	return t.dimension
}

func (t *testPoint) Vector() common.PointVector {
	return t.vector
}

func createPoint(dimension int, lowerBound, upperBound float64) common.Point {
	vector := make([]float64, dimension)
	for i := range vector {
		vector[i] = rand.Float64() * (upperBound - lowerBound)
	}
	return &testPoint{dimension: dimension, vector: vector}
}

func createPoints(nPoints, dimension int, lowerBound, upperBound float64) []common.Point {
	result := make([]common.Point, nPoints)
	for i := range result {
		result[i] = createPoint(dimension, lowerBound, upperBound)
	}
	return result
}
func treeSizeValidator(t *testing.T, nPoints int, tree *balltree.BallTree) {
	assert.Equal(t, nPoints, tree.Size(), "Expecting tree size to match the number of nodes. Tree size: %d, expected: %d", tree.Size(), nPoints)
}

func treeDepthValidator(t *testing.T, nPoints int, tree *balltree.BallTree) {
	treeSizeLowerBound := int(math.Floor(math.Log2(float64(nPoints))))
	treeSizeUpperBound := treeSizeLowerBound + 1
	assert.GreaterOrEqual(t, tree.Depth(), treeSizeLowerBound, "Expecting tree depth to be at least log2(#nodes). Tree depth: %d, expected lower bound: %d", tree.Depth(), treeSizeLowerBound)
	assert.LessOrEqual(t, tree.Depth(), treeSizeUpperBound, "Expecting tree depth to be approximately log2(#nodes). Tree depth: %d, expected upper bound: %d", tree.Depth(), treeSizeUpperBound)
}

func TestCanCreateTree(t *testing.T) {
	nPoints := 1000
	dimension := 3
	points := createPoints(nPoints, dimension, -100, 100)
	tree := balltree.BallTree{}
	tree.Construct(points, dimension)
	treeSizeValidator(t, nPoints, &tree)
	treeDepthValidator(t, nPoints, &tree)
}

func TestCanCreateLargeTree(t *testing.T) {
	nPoints := 1_000_000
	dimension := 7
	points := createPoints(nPoints, dimension, -100, 100)
	tree := balltree.BallTree{}
	tree.Construct(points, dimension)
	treeSizeValidator(t, nPoints, &tree)
	treeDepthValidator(t, nPoints, &tree)
}

func TestCanSearchTree(t *testing.T) {
	nPoints := 1000
	dimension := 3
	points := createPoints(nPoints, dimension, -100, 100)
	testPoint := createPoint(dimension, -100, 100)
	tree := balltree.BallTree{}
	tree.Construct(points, dimension)
	result, err := tree.Search(testPoint.Vector(), 500)
	assert.Nil(t, err, "No error should be returned")
	assert.NotEmpty(t, nPoints, result, "Expecting a non empty search result")
}
