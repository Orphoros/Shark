package cmd

import "github.com/phuslu/log"

func RegisterLogger(level string) {
	var lvl log.Level

	switch level {
	case "debug":
		lvl = log.DebugLevel
	case "info":
		lvl = log.InfoLevel
	case "warn":
		lvl = log.WarnLevel
	case "error":
		lvl = log.ErrorLevel
	case "fatal":
		lvl = log.FatalLevel
	case "panic":
		lvl = log.PanicLevel
	case "trace":
		lvl = log.TraceLevel
	default:
		lvl = log.PanicLevel
	}

	log.DefaultLogger = log.Logger{
		Level:      lvl,
		TimeFormat: "15:04:05",
		Caller:     0,
		Writer: &log.ConsoleWriter{
			ColorOutput:    true,
			QuoteString:    true,
			EndWithMessage: true,
		},
	}

}
