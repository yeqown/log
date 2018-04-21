package log

import (
	"testing"
	// "time"
)

func Test_NewLogger(t *testing.T) {
	l := NewLogger()
	if err := l.SetFileOutput("./testdata", "app"); err != nil {
		t.Error(err)
		t.Fail()
	}

	l.Info("Info")
	l.Error("Error")
	l.Debug("Debug")
	l.Warn("Warn")
	// time.Sleep(100 * time.Second)
	//
	// for i := 0; i < 100; i++ {
	//	l.Info("info loop")
	//	time.Sleep(1 * time.Nanosecond)
	// }
}

func Test_FileSplit(t *testing.T) {
	logPath := "./testdata"
	filename := "app"
	if _, err := openOrCreate(assembleFilepath(logPath, filename)); err != nil {
		t.Error(err)
		t.FailNow()
	}

	if err := renameLogfile(logPath, filename); err != nil {
		t.Error(err)
		t.FailNow()
	}
	// renew file
	if _, err := openOrCreate(assembleFilepath(logPath, filename)); err != nil {
		t.Error(err)
		t.FailNow()
	}
}
