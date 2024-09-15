package resp

import (
	"errors"
	"fmt"
	"redis/pkg/utils"
	"strconv"
	"strings"
)

var (
	ErrValueNotEndWithCrlf = errors.New("value doesn't end with crlf suffix")
)

type unmarshaller func(v string) (any, error)

var (
	unmarshallers map[string]unmarshaller

	crlfCountsPerType = map[string]int{
		stringSuffix:      1,
		errorSuffix:       1,
		intSuffix:         1,
		bulkStringsSuffix: 2,
	}
)

func init() {
	if unmarshallers == nil {
		unmarshallers = map[string]unmarshaller{
			stringSuffix:      simpleStringUnmarshaller,
			errorSuffix:       errorUnmarshaller,
			intSuffix:         intUnmarshaller,
			bulkStringsSuffix: bulkStringsUnmarshaller,
			arraySuffix:       arrayUnmarshaller,
		}
	}
}

func Unmarshal(v any) (any, error) {
	switch v.(type) {
	case string:
		return unmarshal(v.(string))
	case []byte:
		return unmarshal(string(v.([]byte)))
	}

	return nil, fmt.Errorf("can't unmarshal this type %T", v)
}

func unmarshal(v string) (any, error) {
	if len(v) == 0 {
		return nil, errors.New("invalid empty input")
	}

	if u, ok := unmarshallers[string(v[0])]; !ok {
		return nil, fmt.Errorf("ERR syntax error")
	} else {
		return u(v)
	}
}

func simpleStringUnmarshaller(v string) (any, error) {
	if !strings.HasSuffix(v, CRLF) {
		return nil, ErrValueNotEndWithCrlf
	}

	str := removeCrlfSuffix(v[1:])

	if strings.Contains(str, CR) || strings.Contains(str, LF) {
		return nil, errors.New("simple string mustn't contain a CR or LF character")
	}

	return str, nil
}

func errorUnmarshaller(v string) (any, error) {
	if !strings.HasSuffix(v, CRLF) {
		return nil, ErrValueNotEndWithCrlf
	}

	str := removeCrlfSuffix(v[1:])

	if strings.Contains(str, CR) || strings.Contains(str, LF) {
		return nil, errors.New("error string mustn't contain a CR or LF character")
	}

	return str, nil
}

func intUnmarshaller(v string) (any, error) {
	if !strings.HasSuffix(v, CRLF) {
		return nil, ErrValueNotEndWithCrlf
	}

	str := removeCrlfSuffix(v[1:])

	i, err := strconv.Atoi(str)
	if err != nil {
		return nil, errors.New("invalid int parsing")
	}

	return i, nil
}

func bulkStringsUnmarshaller(v string) (any, error) {
	if !strings.HasSuffix(v, CRLF) {
		return nil, ErrValueNotEndWithCrlf
	}

	str := removeCrlfSuffix(v[1:])

	// https://redis.io/docs/latest/develop/reference/protocol-spec/#bulk-strings
	if str == "-1" { // Null bulk strings
		return nil, nil
	}

	elements := strings.Split(str, CRLF)

	if len(elements) != 2 {
		return nil, errors.New("invalid bulk strings structure")
	}

	stringLengthStr, content := elements[0], elements[1]

	stringLength, err := strconv.Atoi(stringLengthStr)
	if err != nil {
		return nil, errors.New("invalid string length value, expected one or more decimal digits (0..9)")
	}

	if stringLength != len(content) {
		return nil, errors.New("content length doesn't match the passed string length")
	}

	return content, nil
}

func arrayUnmarshaller(v string) (any, error) {
	if !strings.HasSuffix(v, CRLF) {
		return nil, ErrValueNotEndWithCrlf
	}

	elements := strings.Split(v, CRLF)

	elementsWithCrlfSuffix := utils.Map(
		utils.Filter(elements, func(s string, i int) bool {
			return s != ""
		}),
		func(t string, i int) string {
			return fmt.Sprintf("%s%s", t, CRLF)
		},
	)

	arrayRespElements, err := parseArrayElements(elementsWithCrlfSuffix)
	if err != nil {
		return nil, err
	}

	if len(arrayRespElements) == 0 {
		return []any{}, nil
	}

	var arrayElements []any

	for _, element := range arrayRespElements {
		if elementValue, err := Unmarshal(element); err == nil {
			arrayElements = append(arrayElements, elementValue)
		} else {
			return nil, err
		}
	}

	return arrayElements, nil
}

func parseArrayElements(parts []string) ([]string, error) {
	var elements []string

	arrayPart := parts[0]

	arrayLength, err := parseArrayLength(arrayPart)
	if err != nil {
		return nil, err
	}

	for i := 1; i < len(parts); i++ {
		part := parts[i]

		opType := string(part[0])

		if opType == arraySuffix {
			nestedArrayElements, err := parseArrayElements(parts[i:])
			if err != nil {
				return nil, err
			}

			nestedArrayStr := strings.Join(
				append(
					[]string{part}, // Add the array part to the beginning of the string
					nestedArrayElements...,
				),
				"",
			)

			elements = append(elements, nestedArrayStr)

			i += len(nestedArrayElements) // skip the nested array elements
		} else {
			partsCount, ok := crlfCountsPerType[opType]
			if !ok {
				partsCount = 1
			}

			elements = append(elements, strings.Join(parts[i:i+partsCount], ""))

			i += partsCount - 1
		}

		if len(elements) == arrayLength {
			break
		}
	}

	return elements, nil
}

func parseArrayLength(v string) (int, error) {
	firstCrlfIdx := strings.Index(v, CRLF)
	if firstCrlfIdx <= 0 {
		return 0, errors.New("invalid array structure")
	}

	arrayLength, err := strconv.Atoi(v[1:firstCrlfIdx])
	if err != nil {
		return 0, err
	}

	return arrayLength, nil
}

func removeCrlfSuffix(v string) string {
	return v[:len(v)-len(CRLF)]
}
