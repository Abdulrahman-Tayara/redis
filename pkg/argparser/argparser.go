package argparser

import (
	"fmt"
	"redis/pkg/ds/mapx"
	"redis/pkg/utils"
	"strings"
)

type ArgInfo struct {
	Name   string
	IsFlag bool
}

func Parse(args []any, schema []ArgInfo) (mapx.IMap[string, any], error) {
	argsMap := utils.SliceToMap(schema, func(v ArgInfo) string {
		return strings.ToLower(v.Name)
	})

	result := make(map[string]any)

	for i := 0; i < len(args); i++ {
		arg := args[i]
		argStr := strings.ToLower(utils.ToString(arg))
		argInfo, ok := argsMap[argStr]
		if !ok {
			return nil, fmt.Errorf("unknown argument: %s", argStr)
		}

		if argInfo.IsFlag {
			result[argInfo.Name] = true
			continue
		}

		if i+1 >= len(args) {
			return nil, fmt.Errorf("missing value for argument: %s", argStr)
		}

		value := args[i+1]

		// Check if the value is an argument
		if valueStr := utils.ToString(value); valueStr != "" {
			if _, ok := argsMap[strings.ToLower(valueStr)]; ok {
				return nil, fmt.Errorf("missing value for argument: %s", argStr)
			}
		}

		result[argInfo.Name] = value
		i++
	}

	for _, argInfo := range schema {
		if argInfo.IsFlag {
			if _, ok := result[argInfo.Name]; !ok {
				result[argInfo.Name] = false
			}
		}
	}

	return mapx.NewMapFromSource(result), nil
}
