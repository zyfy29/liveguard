package cm

import (
	"testing"
)

func TestJsonMarshal(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test JsonMarshal",
			args: args{
				v: map[string]string{
					"key": "value",
				},
			},
			want: "{\"key\":\"value\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JsonMarshal(tt.args.v); got != tt.want {
				t.Errorf("JsonMarshal() = %v, want %v", got, tt.want)
			}
		})
	}
}
