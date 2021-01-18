package log

import (
	"bytes"
	"testing"
)

func Test_appendValue(t *testing.T) {
	type args struct {
		b     *bytes.Buffer
		value interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "case 0",
			args: args{
				b: bytes.NewBuffer(nil),
				value: struct {
					Str        string
					Int        int
					Bool       bool
					unexported string
				}{
					Str:        "12123",
					Int:        2222,
					Bool:       true,
					unexported: "unexported",
				},
			},
		},
		{
			name: "case 1",
			args: args{
				b:     bytes.NewBuffer(nil),
				value: []string{"1111", "2222", "3333"},
			},
		},
		{
			name: "case 2",
			args: args{
				b:     bytes.NewBuffer(nil),
				value: &struct{ Str string }{Str: "string"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appendKeyValue(tt.args.b, tt.name, tt.args.value)
			t.Log(tt.args.b.String())
		})
	}
}
