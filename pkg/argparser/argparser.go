package argparser

import (
	"fmt"
	"redis/pkg/utils"
)

const (
	ArgTagName        = "arg"
	RequiredTagProp   = "required"
	PositionalTagProp = "positional"
)

type ArgInfo struct {
	Name   string
	IsFlag bool
}

func Parse(args []any, schema []ArgInfo) (map[string]any, error) {
	argsMap := utils.SliceToMap(schema, func(v ArgInfo) string {
		return v.Name
	})

	result := make(map[string]any)

	for i := 0; i < len(args); i++ {
		arg := args[i]
		argStr := utils.ToString(arg)
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
			if _, ok := argsMap[valueStr]; ok {
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

	return result, nil
}
