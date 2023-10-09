package common

type PointVector []float64

type Point[T any] interface {
	Dimension() int
	Vector() PointVector
}

type SpacePartitioningTree[T any] interface {
	// constructors
	// Insert(point *PointWithData[T]) error
	Construct(points []*Point[T], dimension int) error
	// accessors
	Search(point PointVector, radius float64) ([]*Point[T], error)
	// KNearestNeighbors(point Point, k int) ([]Node[T], error)
	NodeDimension() int
	Size() int
	Depth() int
}
