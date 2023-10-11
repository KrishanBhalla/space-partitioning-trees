package kdtree_test

import (
	"math/rand"
	"testing"

	"github.com/KrishanBhalla/space-partitioning-trees/pkg/common"
	kdtree "github.com/KrishanBhalla/space-partitioning-trees/pkg/kd_tree"
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

func TestCanCreateTree(t *testing.T) {
	nPoints := 1000
	dimension := 3
	points := createPoints(nPoints, dimension, -100, 100)
	tree := kdtree.KdTree{}
	tree.Construct(points, dimension)
	assert.Equal(t, nPoints, tree.Size(), "Expecting tree size to match the number of nodes. Tree size: %d, expected: %d", tree.Size(), nPoints)
}

func TestCanSearchTree(t *testing.T) {
	nPoints := 1000
	dimension := 3
	points := createPoints(nPoints, dimension, -100, 100)
	testPoint := createPoint(dimension, -100, 100)
	tree := kdtree.KdTree{}
	tree.Construct(points, dimension)
	result, err := tree.Search(testPoint.Vector(), 500)
	assert.Nil(t, err, "No error should be returned")
	assert.NotEmpty(t, nPoints, result, "Expecting a non empty search result")
}
