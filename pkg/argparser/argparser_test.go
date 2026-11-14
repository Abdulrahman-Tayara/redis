package argparser

import (
	"redis/pkg/ds/mapx"
	"redis/pkg/utils"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		args        []any
		schema      []ArgInfo
		expected    mapx.IMap[string, any]
		expectedErr string
	}{
		{
			name:     "empty args",
			args:     []any{},
			expected: mapx.NewMap[string, any](),
		},
		{
			name: "flags",
			args: []any{"GET"},
			schema: []ArgInfo{
				{Name: "GET", IsFlag: true},
				{Name: "SAVE", IsFlag: true},
			},
			expected: mapx.NewMapFromSource(map[string]any{
				"GET":  true,
				"SAVE": false,
			}),
		},
		{
			name: "multiple args",
			args: []any{"SET", "key", "TTL", 12},
			schema: []ArgInfo{
				{Name: "SET", IsFlag: false},
				{Name: "TTL", IsFlag: false},
			},
			expected: mapx.NewMapFromSource(map[string]any{
				"SET": "key",
				"TTL": 12,
			}),
		},
		{
			name: "args with flags",
			args: []any{"SET", "key", "IGNORE", "TTL", 12},
			schema: []ArgInfo{
				{Name: "SET", IsFlag: false},
				{Name: "IGNORE", IsFlag: true},
				{Name: "TTL", IsFlag: false},
			},
			expected: mapx.NewMapFromSource(map[string]any{
				"SET":    "key",
				"IGNORE": true,
				"TTL":    12,
			}),
		},
		{
			name:        "missing value",
			args:        []any{"SET", "key", "TTL"},
			schema:      []ArgInfo{{Name: "SET", IsFlag: false}, {Name: "TTL", IsFlag: false}},
			expectedErr: "missing value for argument: ttl",
		},
		{
			name: "missing value 2",
			args: []any{"SET", "key", "TTL", "IGNORE"},
			schema: []ArgInfo{
				{Name: "SET", IsFlag: false},
				{Name: "TTL", IsFlag: false},
				{Name: "IGNORE", IsFlag: true},
			},
			expectedErr: "missing value for argument: ttl",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := Parse(test.args, test.schema)

			utils.AssertTest(t, test.name, test.expected, actual, test.expectedErr, err)
		})
	}
}
