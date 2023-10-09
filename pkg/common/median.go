package common

import "fmt"

func FindMedianByOrdering(ordering []float64, points []Point) (Point, []Point, []Point, error) {
	if len(ordering) != len(points) {
		return nil, nil, nil, fmt.Errorf("The ordering slice and points slice must have the same length")
	}
	if len(ordering) == 0 {
		return nil, nil, nil, fmt.Errorf("No data passed to FindMedianByOrdering")
	}

	n := len(ordering)
	k := len(ordering) / 2
	pivot, smaller, larger := quickSelect(ordering, points, 0, n-1, k)
	return pivot, smaller, larger, nil
}

func quickSelect(ordering []float64, points []Point, l, r, k int) (Point, []Point, []Point) {
	if l == r { // If the list contains only one element,
		return points[l], points[:l], points[l+1:] // return that element
	}
	pivotIndex := partition(ordering, points, l, r)
	// The pivot is in its final sorted position
	if pivotIndex == k {
		return points[k], points[:k], points[k+1:]
	}
	if pivotIndex < k {
		return quickSelect(ordering, points, pivotIndex+1, r, k)
	}
	return quickSelect(ordering, points, l, pivotIndex-1, k)
}

func partition(ordering []float64, points []Point, l, r int) int {

	x := ordering[r]
	i := l
	for j := l; j < r; j++ {

		if ordering[j] <= x {
			ordering[i], ordering[j] = ordering[j], ordering[i]
			points[i], points[j] = points[j], points[i]
			i++
		}
	}

	ordering[i], ordering[r] = ordering[r], ordering[i]
	points[i], points[r] = points[r], points[i]
	return i
}
