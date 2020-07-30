package log_test

import (
	"testing"

	"github.com/yeqown/log"

	"github.com/stretchr/testify/assert"
)

func Test_findCaller(t *testing.T) {
	type args struct {
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
			args:         args{},
			wantFile:     "log/caller_test.go",
			wantFunction: "github.com/yeqown/log.testCallFn",
			wantLine:     10,
		},
	}
	for _, tt := range tests {
		frm := log.GetCallerForTest()
		gotFile, gotFunction, gotLine := frm.File, frm.Function, frm.Line

		t.Log(gotFile, gotFunction, gotLine)
		assert.Contains(t, gotFile, tt.wantFile)
		assert.Contains(t, gotFunction, tt.wantFunction)
		assert.Equal(t, tt.wantLine, gotLine)
	}
}
