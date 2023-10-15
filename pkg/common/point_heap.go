package common

type PointWithDistance struct {
	Point    Point
	Distance float64
}

type PointWithDistanceHeap []PointWithDistance

func (h PointWithDistanceHeap) Len() int {
	return len(h)
}

func (h PointWithDistanceHeap) Less(i, j int) bool {
	return h[i].Distance < h[j].Distance
}

func (h PointWithDistanceHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *PointWithDistanceHeap) Push(x interface{}) {
	*h = append(*h, x.(PointWithDistance))
}

func (h *PointWithDistanceHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}
