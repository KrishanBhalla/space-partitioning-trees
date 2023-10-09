package common

type Point []float64

type PointWithData[T any] struct {
	Point
	Data T
}

type Node[T any] interface {
	Point() Point
	Data() T
}

type SpacePartitioningTree[T any, N Node[T]] interface {
	// constructors
	// Insert(point *PointWithData[T]) error
	Construct(points []*PointWithData[T], dimension int) error
	// accessors
	Search(point Point, radius float64) ([]*N, error)
	// KNearestNeighbors(point Point, k int) ([]Node[T], error)
	Dimension() int
	Size() int
	Depth() int
}
