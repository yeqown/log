package log

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_appendValue(t *testing.T) {
	type args struct {
		b     *bytes.Buffer
		value interface{}
	}
	tests := []struct {
		name     string
		args     args
		expected string
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
			expected: `case 0="{Str:12123 Int:2222 Bool:true unexported:unexported}"`,
		},
		{
			name: "case 1",
			args: args{
				b:     bytes.NewBuffer(nil),
				value: []string{"1111", "2222", "3333"},
			},
			expected: `case 1="[1111 2222 3333]"`,
		},
		{
			name: "case 2",
			args: args{
				b:     bytes.NewBuffer(nil),
				value: &struct{ Str string }{Str: "string"},
			},
			expected: `case 2="&{Str:string}"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appendKeyValue(tt.args.b, tt.name, tt.args.value, true, true)
			assert.Equal(t, tt.expected, tt.args.b.String())
		})
	}
}

func Test_format(t *testing.T) {
	formatter := newTextFormatter(
		false, false, true, time.RFC3339)
	entry := entry{
		logger:     nil,
		out:        nil,
		formatter:  nil,
		lv:         0,
		withCaller: false,
		fixedField: &fixedField{Timestamp: 1747750112},
		fields:     Fields{"a": "a", "b": "b", "c": "c"},
		ctx:        nil,
		ctxParser:  nil,
	}
	out, err := formatter.Format(&entry, "This is a test message")
	assert.NoError(t, err)
	assert.Equal(t, "[FTL] 2025-05-20T22:08:32+08:00 Fields{a=\"a\" b=\"b\" c=\"c\"} This is a test message\n", string(out))
}
