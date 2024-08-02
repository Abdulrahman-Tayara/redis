package resp

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	CRLF = "\r\n"
	CR   = "\r"
	LF   = "\n"
)

const (
	stringSuffix      = "+"
	bulkStringsSuffix = "$"
	errorSuffix       = "-"
	intSuffix         = ":"
	arraySuffix       = "*"
)

func Marshal(v any) (string, error) {
	var res string
	var err error

	indirectValue := func() any {
		return reflect.Indirect(reflect.ValueOf(v)).Interface()
	}

	kind := reflect.TypeOf(v).Kind()

	if kind == reflect.Slice || kind == reflect.Array {
		return marshalArray(indirectValue())
	}

	switch v.(type) {
	case string:
		res, err = marshalString(indirectValue())
		break
	case error:
		res, err = marshalError(v)
		break
	case int, int8, int16, int32, int64:
		res, err = marshalInt(indirectValue())
		break
	default:
		return "", errors.New("unknown type")
	}

	if err != nil {
		return "", err
	}

	return res, nil
}

func marshalString(v any) (string, error) {
	vStr := v.(string)

	if strings.Contains(vStr, CR) || strings.Contains(vStr, LF) {
		return marshalBulkStrings(vStr)
	}

	return marshalSimpleString(vStr)
}

func marshalBulkStrings(v string) (string, error) {
	bytesLen := len(v)

	return fmt.Sprintf("%s%d%s%s%s", bulkStringsSuffix, bytesLen, CRLF, v, CRLF), nil
}

func marshalSimpleString(v string) (string, error) {
	if strings.Contains(v, CR) {
		return "", errors.New("simple string mustn't contain \\r")
	}
	if strings.Contains(v, LF) {
		return "", errors.New("simple string mustn't contain \\n")
	}

	return fmt.Sprintf("%s%s%s", stringSuffix, v, CRLF), nil
}

func marshalError(v any) (string, error) {
	vErr := v.(error)
	res, err := marshalSimpleString(vErr.Error())
	if err != nil {
		return "", err
	}

	return strings.Replace(res, stringSuffix, errorSuffix, 1), nil
}

func marshalInt(v any) (ret string, retErr error) {
	return fmt.Sprintf("%v%s%s", intSuffix, fmt.Sprintf("%d", v), CRLF), nil
}

func marshalArray(v any) (string, error) {
	value := reflect.ValueOf(v)

	arrayLen := value.Len()

	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf("%s%d%s", arraySuffix, arrayLen, CRLF))

	for i := range arrayLen {
		item := reflect.Indirect(value.Index(i))

		itemStr, err := Marshal(item.Interface())
		if err != nil {
			return "", err
		}

		builder.WriteString(itemStr)
	}

	return builder.String(), nil
}
