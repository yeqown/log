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
	lv           Level     // only log.lv is lte than lv, then it would be written into Writer
	globalFields Fields    // global fields

	callerReporter bool          // log caller or not.
	ctxParser      ContextParser // ContextParser for parse Context

	// formatTime format time or not.
	formatTime bool
	// formatTimeLayout format time layout.
	formatTimeLayout string
	// sortField print fields in order of fields' keys lexicographical order.
	sortField bool
}

func (o *options) level() Level {
	if o == nil {
		return LevelDebug
	}

	return o.lv
}

// isTerminal indicates the w (io.Writer) is a byte output device.
func (o *options) isTerminal() bool {
	if o == nil {
		return true
	}

	return isTerminal(o.w)
}

func isTerminal(w io.Writer) bool {
	if w == nil {
		return false
	}

	fd, ok := w.(*os.File)
	if !ok {
		return false
	}

	fi, err := fd.Stat()
	if err != nil {
		return false
	}

	// os.Stdout is named pipe to /dev/fd/1 (char device)
	// os.Stderr is named pipe to /dev/fd/2 (char device)
	return fi.Mode()&os.ModeNamedPipe == os.ModeNamedPipe || fi.Mode()&os.ModeCharDevice == os.ModeCharDevice
}

func (o *options) writer() io.Writer {
	if o == nil {
		return os.Stdout
	}

	return o.w
}

// withDefault sets os.Stdout as write, debug level,
// terminal open and no global fields.
func withDefault(lo *options) error {
	lo.w = os.Stdout
	lo.lv = LevelDebug
	// lo.stdout = true
	//lo.isTerminal = true
	lo.globalFields = nil
	// using `nonParser` as default to help user to define their own parser
	lo.ctxParser = DefaultContextParserFunc(nonParser)
	lo.formatTime = false
	lo.formatTimeLayout = ""
	lo.sortField = false

	return nil
}

// WithLevel setting the level, this could change dynamic
func WithLevel(lv Level) LoggerOption {
	return func(lo *options) error {
		lo.lv = lv
		return nil
	}
}

// WithStdout output to os.Stdout also.
func WithStdout(v bool) LoggerOption {
	return func(lo *options) error {
		if lo.w != nil && lo.w != os.Stdout {
			// If has set a witer, and the writer isn't os.Stdout, use io.MultiWriter
			// to merge old writer and os.Stdout.
			lo.w = io.MultiWriter(lo.w, os.Stdout)
		}
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

// WithContextParser set a custom ContextParser for parsing context.
// maybe you want to auto log opentracing traceId, this could help you.
func WithContextParser(parser ContextParser) LoggerOption {
	return func(lo *options) error {
		lo.ctxParser = parser
		return nil
	}
}

// WithFieldsSort print fields in order of fields' keys lexicographical order.
func WithFieldsSort(sortField bool) LoggerOption {
	return func(lo *options) error {
		lo.sortField = sortField
		return nil
	}
}
