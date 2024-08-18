package utils

import "reflect"

func IsNull(v any) (ret bool) {
	defer func() {
		if r := recover(); r != nil {
			ret = false
		}
	}()

	if v == nil {
		return true
	}

	value := reflect.ValueOf(v)

	return value.IsNil() || value.IsZero()
}

func Filter[T any](array []T, predicate func(T, int) bool) []T {
	var newArray []T

	for i, item := range array {
		if predicate(item, i) {
			newArray = append(newArray, item)
		}
	}

	return newArray
}

func Map[T any, K any](array []T, mapFunc func(T, int) K) []K {
	var newArray []K

	for i, item := range array {
		newArray = append(newArray, mapFunc(item, i))
	}

	return newArray
}
