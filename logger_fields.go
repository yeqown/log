package log

// Fields to contains a batch field to log
type Fields map[string]interface{}

type fixedField struct {
	File          string `json:"file"`          // filename "xxx.go:132"
	Fn            string `json:"fn"`            // func name
	Timestamp     int64  `json:"timestamp"`     // timestamp
	FormattedTime string `json:"formattedTime"` // formatted time
}
