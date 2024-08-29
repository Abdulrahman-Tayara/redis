package server

import (
	"errors"
	"redis/pkg/resp"
)

func ReadRespCommand(b []byte) (string, []any, error) {
	respValue, err := resp.Unmarshal(b)
	if err != nil {
		return "", nil, err
	}

	switch respValue.(type) {
	case string:
		return respValue.(string), []any{}, nil
	case []any:
		arr := respValue.([]any)
		if len(arr) == 0 {
			return "", nil, errors.New("invalid command, empty array")
		}

		command, ok := arr[0].(string)
		if !ok {
			return "", nil, errors.New("invalid command, command name isn't valid")
		}

		return command, arr[1:], nil
	}

	return "", nil, errors.New("invalid command structure")
}
