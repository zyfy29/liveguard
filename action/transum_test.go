package action

import (
	"bearguard/repo"
	"fmt"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_doTranscriptAndSummarize(t *testing.T) {
	type args struct {
		task repo.Task
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Test doTranscriptAndSummarize",
			args: args{
				task: lo.Must(repo.GetDBTaskByLiveID("1036045407548674048")),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, doTranscriptAndSummarize(tt.args.task), fmt.Sprintf("doTranscriptAndSummarize(%v)", tt.args.task))
		})
	}
}
