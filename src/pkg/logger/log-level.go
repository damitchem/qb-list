package logger

type LogLevel int

var levelNames = map[LogLevel]string{
	Trace:    "Trace",
	Debug:    "Debug",
	Info:     "Info",
	Warn:     "Warn",
	Error:    "Error",
	Critical: "Critical",
	None:     "None",
}

func init() {
	nameLevel = make(map[string]LogLevel)
	for level, name := range levelNames {
		nameLevel[name] = level
	}
}

var nameLevel map[string]LogLevel

func GetLevel(name string) LogLevel {
	if level, ok := nameLevel[name]; ok {
		return level
	}
	return Info
}

const (
	Trace    LogLevel = 0
	Debug    LogLevel = 1
	Info     LogLevel = 2
	Warn     LogLevel = 3
	Error    LogLevel = 4
	Critical LogLevel = 5
	None     LogLevel = 6
)
