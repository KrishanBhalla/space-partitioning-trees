package common

type PointVector []float64

type Point interface {
	Dimension() int
	Vector() PointVector
}

type SpacePartitioningTree interface {
	// constructors
	Construct(points []Point, dimension int) error
	// accessors
	Search(point Point, radius float64) ([]Point, error)
	KNearestNeighbors(point Point, k int) ([]Point, error)
	// Helpers
	NodeDimension() int
	Size() int
	Depth() int
	Points() []Point
}
