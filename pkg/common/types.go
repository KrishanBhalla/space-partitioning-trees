package common

type Point []float64

type PointWithData[T any] struct {
	Point
	Data T
}

type Node[T any] struct {
	Point Point
	Data  T
}
type INode[T any] interface {
	Node() Node[T]
}

type SpacePartitioningTree[T any, N INode[T]] interface {
	// constructors
	// Insert(point *PointWithData[T]) error
	Construct(points []*PointWithData[T], dimension int) error
	// accessors
	Search(point Point, radius float64) ([]*N, error)
	// KNearestNeighbors(point Point, k int) ([]Node[T], error)
	NodeDimension() int
	Size() int
	Depth() int
}
