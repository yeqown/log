package main

import "github.com/yeqown/log"

func main() {
	logger := log.NewLogger(
		log.WithLevel(log.LevelInfo),             // specify the log level
		log.WithFileLog("./logs/app.log", false), // file log
		log.WithStdout(true),                     // this will affect only FileLog mode
	)

	logger.
		WithField("example", "std-and-file").
		Info("this is an example")
}
