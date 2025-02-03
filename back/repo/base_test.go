package repo

import (
	"testing"
)

func Test_getDBMust(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test to call init()",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = getDbMust()
		})
	}
}
