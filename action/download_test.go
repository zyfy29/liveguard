package action

import (
	"bearguard/repo"
	"fmt"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_doDownloadLive(t *testing.T) {
	type args struct {
		task repo.Task
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Test doDownloadLive",
			args: args{
				task: lo.Must(repo.GetDBTaskByLiveID("1036045407548674048")),
			},
			wantErr: assert.NoError,
		},
		{
			name: "Test doDownloadLive not prepared",
			args: args{
				task: lo.Must(repo.GetDBTaskByLiveID("1037118894660980736")),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, doDownloadLive(tt.args.task), fmt.Sprintf("doDownloadLive(%v)", tt.args.task))
		})
	}
}
