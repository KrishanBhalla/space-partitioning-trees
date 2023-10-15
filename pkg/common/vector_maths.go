package common

import (
	"fmt"
	"math"
)

/**
* Returns the L2 (Euclidean) distance between two vectors
 */
func Distance(vec1, vec2 PointVector) (float64, error) {
	if len(vec1) != len(vec2) {
		return 0.0, fmt.Errorf("Points have differing lengths: %d and %d", len(vec1), len(vec2))
	}
	distance := 0.
	for i, v1 := range vec1 {
		v2 := vec2[i]
		distance += (v1 - v2) * (v1 - v2)
	}
	return math.Sqrt(distance), nil
}

/**
* Returns the dot product of two vectors
 */
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

/**
* Computes the result of vec2 - vec1
 */
func Difference(vec1, vec2 PointVector) (PointVector, error) {
	if len(vec1) != len(vec2) {
		return nil, fmt.Errorf("Point Vectors have differing lengths: %d and %d", len(vec1), len(vec2))
	}
	result := make(PointVector, len(vec1))
	for i, v := range vec2 {
		result[i] = v - vec1[i]
	}
	return result, nil
}
