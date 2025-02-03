package action

import (
	"bearguard/repo"
	"fmt"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_doSubmit(t *testing.T) {
	type args struct {
		task repo.Task
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Test doSubmit",
			args: args{
				task: lo.Must(repo.GetDBTaskByLiveID("1036045407548674048")),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, doSubmit(tt.args.task), fmt.Sprintf("doSubmit(%v)", tt.args.task))
		})
	}
}

func Test_retrySubmitTask(t *testing.T) {
	type args struct {
		task repo.Task
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Test retrySubmitFailedTask",
			args: args{
				task: lo.Must(repo.GetDBTaskByLiveID("1036045407548674048")),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, retrySubmitFailedTask(tt.args.task), fmt.Sprintf("retrySubmitFailedTask(%v)", tt.args.task))
		})
	}
}
