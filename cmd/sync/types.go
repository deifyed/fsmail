package sync

type logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
}
