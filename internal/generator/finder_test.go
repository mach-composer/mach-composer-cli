package generator

import (
	"reflect"
	"testing"
)

func Test_findDependsOn(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]any
		want  []string
	}{
		{
			name: "test",
			input: map[string]any{
				"api_management_name": "${module.azure-config.api_management_name}",
			},
			want: []string{
				"module.azure-config",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findDependsOn(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findDependsOn() = %v, want %v", got, tt.want)
			}
		})
	}
}
