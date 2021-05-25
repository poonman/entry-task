package log

import (
	"fmt"
	"log"
	"os"

)

type defaultLogger struct {
	*log.Logger
}

func (l *defaultLogger) Debug(v ...interface{}) {
	_ = l.Output(calldepth, header("DEBUG", fmt.Sprint(v...)))
}

func (l *defaultLogger) Debugf(format string, v ...interface{}) {
	_ = l.Output(calldepth, header("DEBUG", fmt.Sprintf(format, v...)))
}

func (l *defaultLogger) Info(v ...interface{}) {
	_ = l.Output(calldepth, header("INFO ", fmt.Sprint(v...)))
}

func (l *defaultLogger) Infof(format string, v ...interface{}) {
	_ = l.Output(calldepth, header("INFO ", fmt.Sprintf(format, v...)))
}

func (l *defaultLogger) Warn(v ...interface{}) {
	_ = l.Output(calldepth, header("WARN ", fmt.Sprint(v...)))
}

func (l *defaultLogger) Warnf(format string, v ...interface{}) {
	_ = l.Output(calldepth, header("WARN ", fmt.Sprintf(format, v...)))
}

func (l *defaultLogger) Error(v ...interface{}) {
	_ = l.Output(calldepth, header("ERROR", fmt.Sprint(v...)))
}

func (l *defaultLogger) Errorf(format string, v ...interface{}) {
	_ = l.Output(calldepth, header("ERROR", fmt.Sprintf(format, v...)))
}

func (l *defaultLogger) Fatal(v ...interface{}) {
	_ = l.Output(calldepth, header("FATAL", fmt.Sprint(v...)))
	os.Exit(1)
}

func (l *defaultLogger) Fatalf(format string, v ...interface{}) {
	_ = l.Output(calldepth, header("FATAL", fmt.Sprintf(format, v...)))
	os.Exit(1)
}

func (l *defaultLogger) Panic(v ...interface{}) {
	l.Logger.Panic(v...)
}

func (l *defaultLogger) Panicf(format string, v ...interface{}) {
	l.Logger.Panicf(format, v...)
}

func header(lvl, msg string) string {
	return fmt.Sprintf("%s: %s", lvl, msg)
}

