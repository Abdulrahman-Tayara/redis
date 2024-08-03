package utils

import "testing"

func TestIsNull(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  bool
	}{
		{
			"Null",
			nil,
			true,
		},
		{
			"Not Null",
			0,
			false,
		},
		{
			"Null array",
			[]int(nil),
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNull(tt.input); got != tt.want {
				t.Errorf("IsNull() = %v, want %v", got, tt.want)
			}
		})
	}
}
