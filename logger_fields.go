package log

// Fields to contains a batch field to log
type Fields map[string]interface{}

// fixedField json tag should keep pace with logger_formatter.go constant
type fixedField struct {
	File          string `json:"_filepath"` // filename "xxx.go:132"
	Fn            string `json:"_func"`     // func name
	Timestamp     int64  `json:"_ts"`       // timestamp
	FormattedTime string `json:"_fmt_time"` // formatted time
}

// copyFields copy all fields in src to dst
func copyFields(dst, src Fields) {
	for k := range src {
		dst[k] = src[k]
	}
}
