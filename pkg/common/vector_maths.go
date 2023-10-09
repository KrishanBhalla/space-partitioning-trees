package common

import "fmt"

func Distance(point1, point2 Point) (float64, error) {
	if len(point1) != len(point2) {
		return 0.0, fmt.Errorf("Points have differing lengths: %d and %d", len(point1), len(point2))
	}
	distance := 0.
	for i, v1 := range point1 {
		v2 := point2[i]
		distance += (v1 - v2) * (v1 - v2)
	}
	return distance, nil
}

func DotProduct(point1, point2 Point) (float64, error) {
	if len(point1) != len(point2) {
		return 0.0, fmt.Errorf("Points have differing lengths: %d and %d", len(point1), len(point2))
	}
	dotProduct := 0.0
	for i, v1 := range point1 {
		v2 := point2[i]
		dotProduct += (v1 * v2)
	}
	return dotProduct, nil
}
