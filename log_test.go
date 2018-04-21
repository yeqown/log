package log

import (
	"testing"
	"time"
)

func Test_NewLogger(t *testing.T) {
	l := NewLogger()
	l.SetFileOutput("./testdata", "app")

	l.Info("Info")
	l.Error("Error")
	l.Debug("Debug")
	l.Warn("Warn")

	time.Sleep(100 * time.Second)
	//
	// for i := 0; i < 100; i++ {
	//	l.Info("info loop")
	//	time.Sleep(1 * time.Nanosecond)
	// }
}

func Test_FileSplit(t *testing.T) {
	logPath := "./testdata"
	filename := "split"

	l := NewLogger()
	l.SetFileOutput(logPath, filename)
	l.Info("split before")

	renameLogfile(logPath, filename)

	// renew file
	openOrCreate(assembleFilepath(logPath, filename))
	l.SetFileOutput(logPath, filename)
	l.Info("split after")
}
