package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

// LoggerOption to apply single function into `lo`
type LoggerOption func(lo *options) error

// logger options to construct logger
type options struct {
	// variables
	w            io.Writer // writer
	lv           Level     // only log.LV is lte than lv, then it would be written into Writer
	globalFields Fields    // global fields

	// flags
	isTerminal       bool   // to mark the output is file or stdout
	stdout           bool   // output to stdout, only affect when file log mode
	callerReporter   bool   // log caller or not.
	formatTime       bool   // format time or nor.
	formatTimeLayout string // format time layout.

	ctxParser ContextParser // ContextParser for parse Context
}

func (o *options) level() Level {
	if o == nil {
		return LevelDebug
	}

	return o.lv
}

func (o *options) terminal() bool {
	if o == nil {
		return true
	}

	return o.isTerminal
}

func (o *options) writer() io.Writer {
	if o == nil {
		return os.Stdout
	}

	if !o.isTerminal && o.stdout {
		return io.MultiWriter(os.Stdout, o.w)
	}

	return o.w
}

// defaultLoggerOption sets os.Stdout as write, debug level,
// terminal open and no global fields.
func defaultLoggerOption(lo *options) error {
	lo.w = os.Stdout
	lo.lv = LevelDebug
	lo.stdout = true
	lo.isTerminal = true
	lo.globalFields = nil
	// using `nonParser` as default to help user to define their own parser
	lo.ctxParser = DefaultContextParserFunc(nonParser)

	return nil
}

// WithLevel setting the level, this could change dynamic
func WithLevel(lv Level) LoggerOption {
	return func(lo *options) error {
		lo.lv = lv
		return nil
	}
}

// WithStdout output to os.Stdout this only affect when file log is opening
func WithStdout(v bool) LoggerOption {
	return func(lo *options) error {
		lo.stdout = v
		return nil
	}
}

// WithGlobalFields set global fields those would be logged in every log.
func WithGlobalFields(fields Fields) LoggerOption {
	return func(lo *options) error {
		lo.globalFields = fields
		return nil
	}
}

// WithCustomWriter using custom writer to log
func WithCustomWriter(w io.Writer) LoggerOption {
	return func(lo *options) error {
		if w != nil {
			lo.w = w
			lo.isTerminal = false
			lo.stdout = false
		}

		return nil
	}
}

// WithReportCaller b is a switch to open print caller or not.
func WithReportCaller(b bool) LoggerOption {
	return func(lo *options) error {
		lo.callerReporter = b

		return nil
	}
}

// WithTimeFormat to output time as the layout you want.
func WithTimeFormat(b bool, layout string) LoggerOption {
	return func(lo *options) error {
		lo.formatTime = b
		lo.formatTimeLayout = layout
		if lo.formatTimeLayout == "" {
			lo.formatTimeLayout = time.RFC3339
		}
		return nil
	}
}

// WithFileLog store log into file, if autoRotate is set,
// it will start a goroutine to split log file by day.
// TODO(@yeqown): using time round instead of ticker
func WithFileLog(file string, autoRotate bool) LoggerOption {
	return func(lo *options) error {
		abs, err := filepath.Abs(file)
		if err != nil {
			return errors.Wrapf(err, "WithFileLog.Abs file: %s", file)
		}

		dir, pureFilename := filepath.Split(abs)
		if lo.w, err = open(abs); err != nil {
			return errors.Wrapf(err, "WithFileLog.open abs: %s", abs)
		}
		lo.isTerminal = false

		// support autoRotate
		if autoRotate {
			go func() {
				ticker := time.NewTicker(1 * time.Minute)
				for tick := range ticker.C {
					if !shouldSplitByTime(tick) {
						continue
					}

					// rename file to old filename
					if err = rename(dir, pureFilename); err != nil {
						fmt.Printf("rename failed dir: %s, filename: %s err: %v \n", dir, pureFilename, err)
						continue
					}
					// open new file
					if lo.w, err = open(assembleFilename(dir, file, true)); err != nil {
						fmt.Printf("open failed file: %s, err: %v \n", assembleFilename(dir, file, true), err)
						continue
					}

					// record the splitting time
					lastSplitTimestamp = time.Now()
				}
			}()
		}

		return nil
	}
}

// WithContextParser set an custom ContextParser for parsing context.
// maybe you want to auto log opentracing traceId, this could help you.
func WithContextParser(parser ContextParser) LoggerOption {
	return func(lo *options) error {
		lo.ctxParser = parser
		return nil
	}
}
