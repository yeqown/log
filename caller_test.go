package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testCallFn(skip int) (file, fn string, line int) {
	return findCaller(skip)
}

func Test_findCaller(t *testing.T) {
	type args struct {
		skip int
	}
	tests := []struct {
		name         string
		args         args
		wantFile     string
		wantFunction string
		wantLine     int
	}{
		{
			name:         "case 0",
			args:         args{skip: 1},
			wantFile:     "log/caller_test.go",
			wantFunction: "github.com/yeqown/log.testCallFn",
			wantLine:     10,
		},
	}
	for _, tt := range tests {
		gotFile, gotFunction, gotLine := testCallFn(tt.args.skip)
		t.Log(gotFile, gotFunction, gotLine)
		assert.Equal(t, tt.wantFile, gotFile)
		assert.Equal(t, tt.wantFunction, gotFunction)
		assert.Equal(t, tt.wantLine, gotLine)
	}
}
