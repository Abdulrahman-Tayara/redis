package utils

import (
	"reflect"
	"strconv"
)

func BoolPtr(b bool) *bool {
	return &b
}

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

func Count[T any](array []T, predicate func(T, int) bool) (count int) {
	for i, item := range array {
		if predicate(item, i) {
			count++
		}
	}

	return
}

func Values[K comparable, V any](m map[K]V) []V {
	var values []V

	for _, v := range m {
		values = append(values, v)
	}

	return values
}

func SliceToMap[K comparable, V any](slice []V, keyFunc func(V) K) map[K]V {
	m := make(map[K]V)

	for _, v := range slice {
		m[keyFunc(v)] = v
	}

	return m
}

func ToString(v any) string {
	if v == nil {
		return ""
	}

	switch v.(type) {
	case string:
		return v.(string)
	case int:
		return strconv.Itoa(v.(int))
	case int32:
		return strconv.Itoa(int(v.(int32)))
	case int64:
		return strconv.Itoa(int(v.(int64)))
	case float32:
		return strconv.FormatFloat(float64(v.(float32)), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v.(float64), 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v.(bool))
	}

	return ""
}

func MustParseInt(v any) int {
	switch v.(type) {
	case int:
		return v.(int)
	case int32:
		return int(v.(int32))
	case int64:
		return int(v.(int64))
	case string:
		i, err := strconv.Atoi(v.(string))
		if err != nil {
			panic(err)
		}
		return i
	default:
		return 0
	}
}
