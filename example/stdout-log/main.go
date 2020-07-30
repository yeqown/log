package main

import (
	"github.com/yeqown/log"
)

type embed struct {
	FieldA string
	FieldB int
}

func main() {
	// using builtin logger
	log.Debug("this is a Debug log")
	log.Info("this is a Info log")
	log.Warn("this is a Warn log")
	log.Error("this is a Error log")
	log.
		WithField("key1", "value1").
		WithFields(log.Fields{
			"key2": "value2",
			"key3": "value3",
			"key4": "value4",
			"key5": "value5",
			"key6": "value6",
			"key7": "value7",
			"key8": "value8",
		}).Error("test error")

	// using new logger
	newLoggerOuput()
}

func newLoggerOuput() {
	logger, _ := log.NewLogger(
		log.WithLevel(log.LevelInfo),
		log.WithGlobalFields(log.Fields{"global_key": "global_value"}),
	)

	logger.Debug("this is a Debug log") // this will not output
	logger.Info("this is a Info log")
	logger.Warn("this is a Warn log")
	logger.Error("this is a Error log")
	logger.WithField("logger", "it's me").
		WithFields(log.Fields{
			"key2": "value2",
			"key3": "value3",
			"key4": "value4",
			"key5": "value5",
			"key6": "value6",
			"key7": "value7",
			"embed": embed{
				FieldA: "aaa",
				FieldB: 112091,
			},
			"embed_ptr": &embed{
				FieldA: "aaa",
				FieldB: 112091,
			},
		}).Error("test error")
}
