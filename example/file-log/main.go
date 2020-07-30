package main

import "github.com/yeqown/log"

func main() {
	logger, err := log.NewLogger(
		log.WithLevel(log.LevelInfo),             // specify the log level
		log.WithFileLog("./logs/app.log", false), // file log
		log.WithStdout(true),                     // also output to stdout
	)
	if err != nil {
		panic(err)
	}

	logger.
		WithField("example", "std-and-file").
		Info("this is an example")

	logger.Debug("this is a Debug log") // this will not output
	logger.Info("this is a Info log")
	logger.Warn("this is a Warn log")
	logger.Error("this is a Error log")
}
