package log

import "os"

var stdUseColors = os.Getenv("TERM") == "xterm" || os.Getenv("TERM") == "xterm-256color"

type Stream struct {
	WriteFunc func(p []byte) (n int, err error)
	Levels    []uint8
	UseColors bool
}

var StreamStdout = Stream{
	WriteFunc: os.Stdout.Write,
	Levels:    []uint8{LEVEL_DEBUG, LEVEL_INFO},
	UseColors: stdUseColors,
}

var StreamStderr = Stream{
	WriteFunc: os.Stderr.Write,
	Levels:    []uint8{LEVEL_WARN, LEVEL_ERROR, LEVEL_FATAL},
	UseColors: stdUseColors,
}
