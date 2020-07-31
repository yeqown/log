package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_copyFields(t *testing.T) {
	type args struct {
		dst Fields
		src Fields
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "case 0",
			args: args{
				dst: Fields{
					"key1": "value1",
				},
				src: Fields{
					"key1": "value2",
					"key2": "value2",
					"key3": "value2",
					"key4": "value2",
					"key5": "value2",
					"key6": "value2",
					"key7": "value2",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			copyFields(tt.args.dst, tt.args.src)
			assert.Len(t, tt.args.dst, 7)
			assert.Contains(t, tt.args.dst, "key1")
			assert.Equal(t, tt.args.dst["key1"], "value2")
		})
	}
}
