package common

import (
	"fmt"
	"math/rand"
)

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
	var pivotIndex int
	for l < r {

		pivotIndex = rand.Intn(r+1-l) + l
		pivotIndex = partition(ordering, points, l, r, pivotIndex)

		if k < pivotIndex {
			r = pivotIndex - 1
		} else if k > pivotIndex {
			l = pivotIndex + 1
		} else {
			break
		}
	}
	return points[k], points[:k], points[k+1:]
}

func partition(ordering []float64, points []Point, l, r, pivotIndex int) int {
	pivot := ordering[pivotIndex]
	partitionIndex := l

	ordering[pivotIndex], ordering[r] = ordering[r], ordering[pivotIndex]
	points[pivotIndex], points[r] = points[r], points[pivotIndex]

	for i := l; i < r; i++ {

		if ordering[i] <= pivot {
			ordering[partitionIndex], ordering[i] = ordering[i], ordering[partitionIndex]
			points[partitionIndex], points[i] = points[i], points[partitionIndex]
			partitionIndex++
		}
	}

	ordering[partitionIndex], ordering[r] = ordering[r], ordering[partitionIndex]
	points[partitionIndex], points[r] = points[r], points[partitionIndex]
	return partitionIndex
}
