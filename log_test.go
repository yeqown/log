package log

import (
	"testing"
)

func Test_NewLogger(t *testing.T) {
	l := NewLogger()
	l.SetFileOutput("./testdata", "app")

	// l.Fatal("faltal", 2, 3)
	l.Error("Error")
	l.Warn("Warn")
	l.Info("Info")
	l.Debug("Debug")

	t.Log(timeToSplit())
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

func Test_DefaultLoggerOutput(t *testing.T) {
	SetFileOutput("./testdata", "default.log")

	a := new(int)
	*a = 10

	b := struct {
		Name string
		Age  int
	}{"test", 10}

	// Fatal("faltal", 2, 3)
	Error("Error")
	Warn("Warn")
	Info("Info", *a, b)
	Debug("Debug")
}
