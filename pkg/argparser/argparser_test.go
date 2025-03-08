package argparser

import (
	"redis/pkg/ds/mapx"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		args        []any
		schema      []ArgInfo
		expected    map[string]any
		expectedErr string
	}{
		{
			name:     "empty args",
			args:     []any{},
			expected: make(map[string]any),
		},
		{
			name: "flags",
			args: []any{"GET"},
			schema: []ArgInfo{
				{Name: "GET", IsFlag: true},
				{Name: "SAVE", IsFlag: true},
			},
			expected: map[string]any{
				"GET":  true,
				"SAVE": false,
			},
		},
		{
			name: "multiple args",
			args: []any{"SET", "key", "TTL", 12},
			schema: []ArgInfo{
				{Name: "SET", IsFlag: false},
				{Name: "TTL", IsFlag: false},
			},
			expected: map[string]any{
				"SET": "key",
				"TTL": 12,
			},
		},
		{
			name: "args with flags",
			args: []any{"SET", "key", "IGNORE", "TTL", 12},
			schema: []ArgInfo{
				{Name: "SET", IsFlag: false},
				{Name: "IGNORE", IsFlag: true},
				{Name: "TTL", IsFlag: false},
			},
			expected: map[string]any{
				"SET":    "key",
				"IGNORE": true,
				"TTL":    12,
			},
		},
		{
			name:        "missing value",
			args:        []any{"SET", "key", "TTL"},
			schema:      []ArgInfo{{Name: "SET", IsFlag: false}, {Name: "TTL", IsFlag: false}},
			expectedErr: "missing value for argument: TTL",
		},
		{
			name: "missing value 2",
			args: []any{"SET", "key", "TTL", "IGNORE"},
			schema: []ArgInfo{
				{Name: "SET", IsFlag: false},
				{Name: "TTL", IsFlag: false},
				{Name: "IGNORE", IsFlag: true},
			},
			expectedErr: "missing value for argument: TTL",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := Parse(test.args, test.schema)
			if err != nil {
				if err.Error() != test.expectedErr {
					t.Errorf("expected error: %s, got: %s", test.expectedErr, err.Error())
				}
			}

			if !actual.IsEqual(mapx.NewMapFromSource(test.expected)) {
				t.Errorf("expected: %v, got: %v", test.expected, actual)
			}
		})
	}
}
