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
	log.Info(1, 2, 3, 4, 5)
	log.Infof("this is format: %d", 2)

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
	logger, _ := log.NewLogger(
		log.WithLevel(log.LevelError),
		log.WithGlobalFields(log.Fields{"global_key": "global_value"}),
	)
	logger.Info(1, 2, 3, 4, 5)
	logger.Infof("this is format: %d", 2)
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
