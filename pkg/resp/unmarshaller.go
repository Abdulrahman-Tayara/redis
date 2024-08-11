package resp

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrValueNotEndWithCrlf = errors.New("value doesn't end with crlf suffix")

	ErrMultipleCrlfFound = errors.New("multiple crlf were found")
)

type unmarshaller func(v string) (any, error)

var (
	unmarshallers = map[string]unmarshaller{
		stringSuffix: simpleStringUnmarshaller,
		errorSuffix:  errorUnmarshaller,
		intSuffix:    intUnmarshaller,
	}
)

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
		return nil, fmt.Errorf("invalid value %v", v)
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

func removeCrlfSuffix(v string) string {
	return v[:len(v)-len(CRLF)]
}
