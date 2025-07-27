package log

type Level struct {
	Name  string
	Index uint8
	Color string
}

type LevelMap map[uint8]Level

const (
	LEVEL_DEBUG uint8 = iota
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
	LEVEL_FATAL
)

var StdLevelMap = LevelMap{
	LEVEL_DEBUG: {Name: "DEBUG", Index: LEVEL_DEBUG, Color: "\033[1;37m\033[44m"},
	LEVEL_INFO:  {Name: "INFO", Index: LEVEL_INFO, Color: "\033[1;30m\033[47m"},
	LEVEL_WARN:  {Name: "WARN", Index: LEVEL_WARN, Color: "\033[1;37m\033[43m"},
	LEVEL_ERROR: {Name: "ERROR", Index: LEVEL_ERROR, Color: "\033[1;37m\033[41m"},
	LEVEL_FATAL: {Name: "FATAL", Index: LEVEL_FATAL, Color: "\033[1;30m\033[41m"},
}
