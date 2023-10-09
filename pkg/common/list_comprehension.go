package common

type Predicate[X any] func(elem X) bool
type Mapper[X any, Y any] func(elem X) Y
type Reducer[X any, Y any] func(memo Y, elem X) Y

func Map[X any, Y any](collection []X, fn Mapper[X, Y]) []Y {
	result := make([]Y, len(collection))
	for i, item := range collection {
		result[i] = fn(item)
	}
	return result
}

func Reduce[X any, Y any](collection []X, init Y, fn Reducer[X, Y]) Y {
	result := init
	for _, item := range collection {
		result = fn(result, item)
	}
	return result
}

func Filter[X any](collection []X, fn Predicate[X]) []X {
	var result []X
	for _, item := range collection {
		if fn(item) {
			result = append(result, item)
		}
	}
	return result
}
