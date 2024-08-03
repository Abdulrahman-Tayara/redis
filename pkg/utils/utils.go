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
