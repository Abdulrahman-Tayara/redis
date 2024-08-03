package resp

import (
	"errors"
	"fmt"
	"redis/pkg/utils"
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
	nullSuffix        = "_"
	boolSuffix        = "#"
)

type marshaller interface {
	isValid(v any) bool
	marshal(input *marshalInput) (string, error)
}

type marshalInput struct {
	value             any
	indirectValueFunc func() any
}

var marshallers = []marshaller{
	&stringMarshaller{},
	&errorMarshaller{},
	&intMarshaller{},
	&booleanMarshaller{},
	&arrayMarshaller{},
	&nullMarshaller{},
}

func Marshal(v any) (string, error) {
	input := &marshalInput{
		value: v,
		indirectValueFunc: func() any {
			return reflect.Indirect(reflect.ValueOf(v)).Interface()
		},
	}

	for _, m := range marshallers {
		if m.isValid(v) {
			return m.marshal(input)
		}
	}

	return "", fmt.Errorf("unsupported type: %T", v)
}

type stringMarshaller struct{}

func (m *stringMarshaller) isValid(v any) bool {
	switch v.(type) {
	case string:
		return true
	default:
		return false
	}
}

func (m *stringMarshaller) marshal(input *marshalInput) (string, error) {
	vStr := input.value.(string)

	if strings.Contains(vStr, CR) || strings.Contains(vStr, LF) {
		return m.marshalBulkStrings(vStr)
	}

	return m.marshalSimpleString(vStr)
}

func (m *stringMarshaller) marshalBulkStrings(v string) (string, error) {
	bytesLen := len(v)

	return fmt.Sprintf("%s%d%s%s%s", bulkStringsSuffix, bytesLen, CRLF, v, CRLF), nil
}

func (m *stringMarshaller) marshalSimpleString(v string) (string, error) {
	if strings.Contains(v, CR) {
		return "", errors.New("simple string mustn't contain \\r")
	}
	if strings.Contains(v, LF) {
		return "", errors.New("simple string mustn't contain \\n")
	}

	return fmt.Sprintf("%s%s%s", stringSuffix, v, CRLF), nil
}

type errorMarshaller struct{}

func (m *errorMarshaller) isValid(v any) bool {
	switch v.(type) {
	case error:
		return true
	default:
		return false
	}
}

func (m *errorMarshaller) marshal(input *marshalInput) (string, error) {
	vErr := input.value.(error)

	if strings.Contains(vErr.Error(), CR) {
		return "", errors.New("error mustn't contain \\r")
	}
	if strings.Contains(vErr.Error(), LF) {
		return "", errors.New("simple string mustn't contain \\n")
	}

	return fmt.Sprintf("%s%s%s", errorSuffix, vErr.Error(), CRLF), nil
}

type intMarshaller struct{}

func (m *intMarshaller) isValid(v any) bool {
	switch v.(type) {
	case int, int8, int16, int32, int64:
		return true
	default:
		return false
	}
}

func (m *intMarshaller) marshal(input *marshalInput) (string, error) {
	return fmt.Sprintf("%v%s%s", intSuffix, fmt.Sprintf("%d", input.value), CRLF), nil
}

type arrayMarshaller struct{}

func (m *arrayMarshaller) isValid(v any) bool {
	if v == nil {
		return false
	}

	kind := reflect.TypeOf(v).Kind()

	return kind == reflect.Slice || kind == reflect.Array
}

func (m *arrayMarshaller) marshal(input *marshalInput) (string, error) {
	v := input.indirectValueFunc()

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

type nullMarshaller struct{}

func (m *nullMarshaller) isValid(v any) bool {
	return utils.IsNull(v)
}

func (m *nullMarshaller) marshal(input *marshalInput) (string, error) {
	return fmt.Sprintf("%s%s", nullSuffix, CRLF), nil
}

type booleanMarshaller struct{}

func (m *booleanMarshaller) isValid(v any) bool {
	switch v.(type) {
	case bool:
		return true
	default:
		return false
	}
}

func (m *booleanMarshaller) marshal(input *marshalInput) (string, error) {
	boolStr := "t"
	if !input.value.(bool) {
		boolStr = "f"
	}
	return fmt.Sprintf("%s%s%s", boolSuffix, boolStr, CRLF), nil
}
