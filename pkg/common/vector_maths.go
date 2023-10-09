package common

import "fmt"

func Distance(vec1, vec2 PointVector) (float64, error) {
	if len(vec1) != len(vec2) {
		return 0.0, fmt.Errorf("Points have differing lengths: %d and %d", len(vec1), len(vec2))
	}
	distance := 0.
	for i, v1 := range vec1 {
		v2 := vec2[i]
		distance += (v1 - v2) * (v1 - v2)
	}
	return distance, nil
}

func DotProduct(vec1, vec2 PointVector) (float64, error) {
	if len(vec1) != len(vec2) {
		return 0.0, fmt.Errorf("Point Vectors have differing lengths: %d and %d", len(vec1), len(vec2))
	}
	dotProduct := 0.0
	for i, v1 := range vec1 {
		v2 := vec2[i]
		dotProduct += (v1 * v2)
	}
	return dotProduct, nil
}
