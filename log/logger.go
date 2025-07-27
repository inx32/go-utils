package log

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/inx32/go-utils/hooks"
)

const (
	LEVEL_NAME_FILL     = 5
	INVALID_LEVEL_COLOR = "\033[1;30m\033[41m"
	RESET_COLOR         = "\033[0m"
	ARG_COLOR           = "\033[2;37m"
)

type Logger struct {
	MinLevel uint8
	SigLoop  hooks.SigLoop
	LevelMap LevelMap
	Streams  []*Stream
}

func (l *Logger) exit(code int) {
	if l.SigLoop == nil {
		l.SigLoop = hooks.DefaultSigLoop()
	}
	l.SigLoop.Exit(code)
}

func (l *Logger) AddStream(s *Stream) {
	l.Streams = append(l.Streams, s)
}

func (l *Logger) Log(level uint8, msg any, args ...any) {
	if level < l.MinLevel {
		return
	}

	lvl, ok := l.LevelMap[level]
	if !ok {
		lvl = Level{
			Name:  strings.Repeat("?", LEVEL_NAME_FILL),
			Color: INVALID_LEVEL_COLOR,
		}
	}

	if len(lvl.Name) > LEVEL_NAME_FILL {
		panic("length of the level name is longer than LEVEL_NAME_FILL")
	}

	for _, s := range l.Streams {
		if slices.Contains(s.Levels, level) {
			s.WriteFunc([]byte(formatMsg(s.UseColors, lvl, msg, args...)))
		}
	}
}

func (l *Logger) Logf(level uint8, format string, args ...any) {
	l.Log(level, fmt.Sprintf(format, args...))
}

func (l *Logger) Debug(msg any, args ...any) {
	l.Log(LEVEL_DEBUG, msg, args...)
}

func (l *Logger) Debugf(format string, args ...any) {
	l.Logf(LEVEL_DEBUG, format, args...)
}

func (l *Logger) Info(msg any, args ...any) {
	l.Log(LEVEL_INFO, msg, args...)
}

func (l *Logger) Infof(format string, args ...any) {
	l.Logf(LEVEL_INFO, format, args...)
}

func (l *Logger) Warn(msg any, args ...any) {
	l.Log(LEVEL_WARN, msg, args...)
}

func (l *Logger) Warnf(format string, args ...any) {
	l.Logf(LEVEL_WARN, format, args...)
}

func (l *Logger) Error(msg any, args ...any) {
	l.Log(LEVEL_ERROR, msg, args...)
}

func (l *Logger) Errorf(format string, args ...any) {
	l.Logf(LEVEL_ERROR, format, args...)
}

func (l *Logger) Fatal(code int, msg any, args ...any) {
	l.Log(LEVEL_FATAL, msg, args...)
	l.exit(code)
}

func (l *Logger) Fatalf(code int, format string, args ...any) {
	l.Logf(LEVEL_FATAL, format, args...)
	l.exit(code)
}

func NewLogger() *Logger {
	l := Logger{LevelMap: StdLevelMap, MinLevel: LEVEL_INFO}
	l.AddStream(&StreamStdout)
	l.AddStream(&StreamStderr)
	if os.Getenv("DEBUG") == "1" {
		l.MinLevel = LEVEL_DEBUG
	}
	return &l
}

var DefaultLogger = Logger{
	MinLevel: LEVEL_INFO,
	SigLoop:  nil,
	LevelMap: StdLevelMap,
	Streams:  []*Stream{&StreamStdout, &StreamStderr},
}

func Debug(msg any, args ...any) {
	DefaultLogger.Debug(msg, args...)
}

func Info(msg any, args ...any) {
	DefaultLogger.Info(msg, args...)
}

func Warn(msg any, args ...any) {
	DefaultLogger.Warn(msg, args...)
}

func Error(msg any, args ...any) {
	DefaultLogger.Error(msg, args...)
}

func Fatal(code int, msg any, args ...any) {
	DefaultLogger.Fatal(code, msg, args...)
}

func Debugf(format string, args ...any) {
	DefaultLogger.Debugf(format, args...)
}

func Infof(format string, args ...any) {
	DefaultLogger.Infof(format, args...)
}

func Warnf(format string, args ...any) {
	DefaultLogger.Warnf(format, args...)
}

func Errorf(format string, args ...any) {
	DefaultLogger.Errorf(format, args...)
}

func Fatalf(code int, format string, args ...any) {
	DefaultLogger.Fatalf(code, format, args...)
}

func init() {
	if os.Getenv("DEBUG") == "1" {
		DefaultLogger.MinLevel = LEVEL_DEBUG
	}
}
