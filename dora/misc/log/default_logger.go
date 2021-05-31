package log

import (
	"fmt"
	"log"
	"os"
)

type defaultLogger struct {
	*log.Logger
	level Level
}

func (l *defaultLogger) SetLevelByString(level string) {
	var lvl Level
	switch level{
	case "DEBUG":
		lvl = DEBUG
	case "INFO":
		lvl = INFO
	case "ERROR":
		lvl = ERROR
	case "FATAL":
		lvl = FATAL
	case "PANIC":
		lvl = PANIC
	}

	l.SetLevel(lvl)
}

func (l *defaultLogger) SetLevel(level Level) {
	l.level = level
}

func (l *defaultLogger) GetLevel() Level {
	return l.level
}

func (l *defaultLogger) Debug(v ...interface{}) {
	if l.level > DEBUG {
		return
	}
	_ = l.Output(calldepth, header("DEBUG", fmt.Sprint(v...)))
}

func (l *defaultLogger) Debugf(format string, v ...interface{}) {
	if l.level > DEBUG {
		return
	}
	_ = l.Output(calldepth, header("DEBUG", fmt.Sprintf(format, v...)))
}

func (l *defaultLogger) Info(v ...interface{}) {
	if l.level > INFO {
		return
	}
	_ = l.Output(calldepth, header("INFO ", fmt.Sprint(v...)))
}

func (l *defaultLogger) Infof(format string, v ...interface{}) {
	if l.level > INFO {
		return
	}
	_ = l.Output(calldepth, header("INFO ", fmt.Sprintf(format, v...)))
}

func (l *defaultLogger) Warn(v ...interface{}) {
	if l.level > WARN {
		return
	}
	_ = l.Output(calldepth, header("WARN ", fmt.Sprint(v...)))
}

func (l *defaultLogger) Warnf(format string, v ...interface{}) {
	if l.level > WARN {
		return
	}
	_ = l.Output(calldepth, header("WARN ", fmt.Sprintf(format, v...)))
}

func (l *defaultLogger) Error(v ...interface{}) {
	if l.level > ERROR {
		return
	}
	_ = l.Output(calldepth, header("ERROR", fmt.Sprint(v...)))
}

func (l *defaultLogger) Errorf(format string, v ...interface{}) {
	if l.level > ERROR {
		return
	}
	_ = l.Output(calldepth, header("ERROR", fmt.Sprintf(format, v...)))
}

func (l *defaultLogger) Fatal(v ...interface{}) {
	if l.level > FATAL {
		return
	}
	_ = l.Output(calldepth, header("FATAL", fmt.Sprint(v...)))
	os.Exit(1)
}

func (l *defaultLogger) Fatalf(format string, v ...interface{}) {
	if l.level > FATAL {
		return
	}
	_ = l.Output(calldepth, header("FATAL", fmt.Sprintf(format, v...)))
	os.Exit(1)
}

func (l *defaultLogger) Panic(v ...interface{}) {
	if l.level > PANIC {
		return
	}
	l.Logger.Panic(v...)
}

func (l *defaultLogger) Panicf(format string, v ...interface{}) {
	if l.level > PANIC {
		return
	}
	l.Logger.Panicf(format, v...)
}

func header(lvl, msg string) string {
	return fmt.Sprintf("%s: %s", lvl, msg)
}
