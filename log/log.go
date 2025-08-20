package log

const (
	TRACE_ID string = "traceID"
)

type Logger interface {
	Info(message string, args ...interface{})
	Debug(message string, args ...interface{})
	Error(message string, args ...interface{})
	Warn(message string, args ...interface{})
}
