package util

type LoggerDummy struct{}

func NewLoggerDummy() Logger {
	return &LoggerDummy{}
}
func (l *LoggerDummy) Log(key string, values map[string]string) {
}
func (l *LoggerDummy) PushSection(s string) {
}

func (l *LoggerDummy) PopSection() {
}

func (l *LoggerDummy) PushRelease(s string) {
}

func (l *LoggerDummy) PopRelease() {
}

func (l *LoggerDummy) SetLogLevel(logLevel string) {
}
