package log

// Fields to contains a batch field to log
type Fields map[string]interface{}

type fixedField struct {
	File          string `json:"file"`          // filename
	Fn            string `json:"fn"`            // func name
	Line          int    `json:"line"`          // line no
	Timestamp     int64  `json:"timestamp"`     // timestamp
	FormattedTime string `json:"formattedTime"` // formatted time
}
