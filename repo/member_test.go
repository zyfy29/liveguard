package repo

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDBMembers(t *testing.T) {
	tests := []struct {
		name    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "Get members",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDBMembers()
			if !tt.wantErr(t, err, fmt.Sprintf("GetDBMembers()")) {
				return
			}
			for i, v := range got {
				t.Log(i, v)
			}
		})
	}
}
