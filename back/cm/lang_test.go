package cm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvert2Zhcn(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test Convert2Zhcn",
			args: args{
				text: "自然語言處理是人工智能領域中的一個重要方向。",
			},
			want: "自然语言处理是人工智能领域中的一个重要方向。",
		},
		{
			name: "english",
			args: args{
				text: "Natural language processing is an important direction in the field of artificial intelligence.",
			},
			want: "Natural language processing is an important direction in the field of artificial intelligence.",
		},
		{
			name: "mixed",
			args: args{
				text: "自然語言處理是人工智能領域中的一個重要方向。Natural language processing is an important direction in the field of artificial intelligence.",
			},
			want: "自然语言处理是人工智能领域中的一个重要方向。Natural language processing is an important direction in the field of artificial intelligence.",
		},
		{
			name: "日本語",
			args: args{
				text: "自然言語処理は、人工知能の分野で重要な方向です。",
			},
			want: "自然言语処理は、人工知能の分野で重要な方向です。",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, Convert2Zhcn(tt.args.text), "Convert2Zhcn(%v)", tt.args.text)
		})
	}
}
