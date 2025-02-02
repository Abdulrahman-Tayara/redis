package resp

import (
	"bufio"
	"errors"
	"strconv"
	"strings"
)

var (
	ErrValueNotEndWithCrlf = errors.New("value doesn't end with crlf suffix")
)

func Unmarshal(r *bufio.Reader) ([]any, error) {
	bufReader := r

	values := []any{}

	for {
		value, err := read(bufReader)
		if err != nil {
			return nil, err
		}
		values = append(values, value)

		if bufReader.Buffered() == 0 {
			break
		}
	}

	return values, nil
}

func read(reader *bufio.Reader) (any, error) {
	dataType, _, err := reader.ReadRune()
	if err != nil {
		return nil, err
	}

	switch string(dataType) {
	case stringSuffix:
		return readString(reader)
	case errorSuffix:
		return readString(reader)
	case bulkStringsSuffix:
		return readBulkStrings(reader)
	case intSuffix:
		return readInt(reader)
	case doubleSuffix:
		return readDouble(reader)
	case boolSuffix:
		return readBool(reader)
	case nullSuffix:
		return nil, nil
	case arraySuffix:
		return readArray(reader)
	case mapSuffix:
		return readMap(reader)
	}

	return nil, errors.New("invalid data format")
}

func readString(reader *bufio.Reader) (string, error) {
	value, err := readUntilCRLF(reader, bytesToString)
	if err != nil {
		return "", err
	}
	strValue := string(value)
	if strings.Contains(strValue, string(LF)) || strings.Contains(strValue, string(CR)) {
		return "", errors.New("simple string mustn't contain a CR or LF character")
	}
	return strValue, nil
}

func readBulkStrings(reader *bufio.Reader) (string, error) {
	length, err := readUntilCRLF(reader, bytesToInt)
	if err != nil {
		return "", errors.New("invalid string length value, expected one or more decimal digits (0..9)")
	}
	if length == 0 {
		reader.Discard(len(CRLF))
		return "", nil
	}

	bytes := make([]byte, length)
	size, err := reader.Read(bytes)
	if err != nil {
		return "", nil
	}

	reader.Discard(len(CRLF))

	return string(bytes[:size]), nil
}

func readInt(reader *bufio.Reader) (int64, error) {
	return readUntilCRLF(reader, bytesToInt)
}

func readDouble(reader *bufio.Reader) (float64, error) {
	return readUntilCRLF(reader, bytesToFloat64)
}

func readBool(reader *bufio.Reader) (bool, error) {
	str, err := readUntilCRLF(reader, bytesToString)
	if err != nil {
		return false, err
	}

	return strconv.ParseBool(str)
}

func readArray(reader *bufio.Reader) ([]any, error) {
	length, err := readUntilCRLF(reader, bytesToInt)
	if err != nil {
		return nil, err
	}
	if length == 0 {
		return []any{}, nil
	}

	elements := make([]any, length)
	for i := 0; i < int(length); i++ {
		element, err := read(reader)
		if err != nil {
			return nil, err
		}
		elements[i] = element
	}
	return elements, nil
}

func readMap(reader *bufio.Reader) (map[any]any, error) {
	length, err := readUntilCRLF(reader, bytesToInt)
	if err != nil {
		return nil, err
	}
	if length == 0 {
		return map[any]any{}, nil
	}

	elements := make(map[any]any, length)
	for i := 0; i < int(length); i++ {
		key, err := read(reader)
		if err != nil {
			return nil, err
		}
		value, err := read(reader)
		if err != nil {
			return nil, err
		}
		elements[key] = value
	}
	return elements, nil
}

// ------ Helpers

func readUntilCRLF[T any](reader *bufio.Reader, formatter func([]byte) (T, error)) (T, error) {
	var ret T
	var err error

	value, err := reader.ReadSlice(CR)
	if err != nil {
		return ret, ErrValueNotEndWithCrlf
	}
	nextRune, err := reader.Peek(1)
	if err != nil {
		return ret, err
	}
	if string(nextRune) != string(LF) {
		return ret, errors.New("invalid format")
	}

	// Discard LF
	if n, _ := reader.Discard(1); n <= 0 {
		return ret, ErrValueNotEndWithCrlf
	}

	value = value[:len(value)-1] // Remove CR

	ret, err = formatter(value)
	return ret, err
}

func bytesToString(bytes []byte) (string, error) {
	return string(bytes), nil
}

func bytesToInt(bytes []byte) (int64, error) {
	v, err := strconv.ParseInt(string(bytes), 10, 64)
	if err != nil {
		return 0, errors.New("invalid int parsing")
	}
	return v, nil
}

func bytesToFloat64(bytes []byte) (float64, error) {
	return strconv.ParseFloat(string(bytes), 64)
}
