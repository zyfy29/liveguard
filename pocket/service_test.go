package pocket

import (
	"bearguard/cm"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testClient = GetClient()

func TestValidateToken(t *testing.T) {
	tests := []struct {
		name    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "nihao",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testClient.ValidateToken()
			if !tt.wantErr(t, err, fmt.Sprintf("ValidateToken()")) {
				return
			}
			t.Log(got)
		})
	}
}

func TestGetHistoryLives(t *testing.T) {
	type args struct {
		memberID string
		nextTime int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "normal",
			args: args{
				memberID: "63566",
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := testClient.GetHistoryLives(tt.args.memberID, tt.args.nextTime)
			if !tt.wantErr(t, err, fmt.Sprintf("GetHistoryLives(%v, %v)", tt.args.memberID, tt.args.nextTime)) {
				return
			}
			for i, v := range got {
				t.Log(i, cm.JsonMarshal(v))
			}
			t.Log(got1)
		})
	}
}

func TestClient_GetLiveInfo(t *testing.T) {
	type args struct {
		liveId string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "normal",
			args:    args{liveId: "1037118894660980736"},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testClient.GetLiveInfo(tt.args.liveId)
			if !tt.wantErr(t, err, fmt.Sprintf("GetLiveInfo(%v)", tt.args.liveId)) {
				return
			}
			t.Log(cm.JsonMarshal(got))
		})
	}
}
