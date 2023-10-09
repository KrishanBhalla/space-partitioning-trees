package common

type PointVector []float64

type Point interface {
	Dimension() int
	Vector() PointVector
}

type SpacePartitioningTree interface {
	// constructors
	// Insert(point *PointWithData[T]) error
	Construct(points []Point, dimension int) error
	// accessors
	Search(point PointVector, radius float64) ([]Point, error)
	// KNearestNeighbors(point Point, k int) ([]Node[T], error)
	NodeDimension() int
	Size() int
	Depth() int
}
