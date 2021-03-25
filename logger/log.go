package logger

import (
	"log"
	"fmt"
	"os"
)

const (
	LEVEL_SILENT = -1
	LEVEL_FATAL = iota
	LEVEL_ALERT
	LEVEL_ERROR
	LEVEL_WARN
	LEVEL_CONFIG
	LEVEL_INFO
	LEVEL_DEBUG
	NR_LEVELS
)


type LoggerIF interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Config(v ...interface{})
	Configf(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Alert(v ...interface{})
	Alertf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})

	Close()
}

type logger struct {
	logLevel int
	loggers  []*log.Logger
	prefix   string
	handler  Handler
	closed   bool
	flags    int
}

func (l *logger) Close() {
	if l.closed {
		return
	}
	l.closed = true

	l.handler.Close()
}


func (l *logger) Debug(v ...interface{}) {
	l.loggers[LEVEL_DEBUG].Print(v...)
}

func (l *logger) Debugf(format string, v ...interface{}) {
	l.loggers[LEVEL_DEBUG].Printf(format, v...)
}

func (l *logger) Info(v ...interface{}) {
	l.loggers[LEVEL_INFO].Print(v...)
}

func (l *logger) Infof(format string, v ...interface{}) {
	l.loggers[LEVEL_INFO].Printf(format, v...)
}

func (l *logger) Config(v ...interface{}) {
	l.loggers[LEVEL_CONFIG].Print(v...)
}

func (l *logger) Configf(format string, v ...interface{}) {
	l.loggers[LEVEL_CONFIG].Printf(format, v...)
}

func (l *logger) Warn(v ...interface{}) {
	l.loggers[LEVEL_WARN].Print(v...)
}

func (l *logger) Warnf(format string, v ...interface{}) {
	l.loggers[LEVEL_WARN].Printf(format, v...)
}

func (l *logger) Error(v ...interface{}) {
	l.loggers[LEVEL_ERROR].Print(v...)
}

func (l *logger) Errorf(format string, v ...interface{}) {
	l.loggers[LEVEL_ERROR].Printf(format, v...)
}

func (l *logger) Alert(v ...interface{}) {
	l.loggers[LEVEL_ALERT].Print(v...)
}

func (l *logger) Alertf(format string, v ...interface{}) {
	l.loggers[LEVEL_ALERT].Printf(format, v...)
}

func (l *logger) Fatal(v ...interface{}) {
	l.loggers[LEVEL_FATAL].Fatal(v...)
}

func (l *logger) Fatalf(format string, v ...interface{}) {
	l.loggers[LEVEL_FATAL].Fatalf(format, v...)
}

var logLevelToName map[int]string

func init() {
	logLevelToName = make(map[int]string, NR_LEVELS)
	logLevelToName[LEVEL_DEBUG] = "[Debug]"
	logLevelToName[LEVEL_INFO] = "[Info]"
	logLevelToName[LEVEL_CONFIG] = "[Config]"
	logLevelToName[LEVEL_WARN] = "[Warning]"
	logLevelToName[LEVEL_ERROR] = "[Error]"
	logLevelToName[LEVEL_ALERT] = "[Alert]"
	logLevelToName[LEVEL_FATAL] = "[Fatal]"
}

func NewLogger(handler Handler, prefix string, logLevel int,flags int) *logger {
	p := new(logger)
	p.loggers = make([]*log.Logger, NR_LEVELS)
	p.prefix = prefix
	p.Set(logLevel,flags,handler)
	return p
}

func (p *logger) Set(logLevel int,flags int, handler Handler) {
	if logLevel > LEVEL_DEBUG {
		logLevel = LEVEL_DEBUG
	}

	if handler == nil {
		p.handler,_= NewNullHandler()
	} else {
		p.handler = handler
	}

	p.flags = flags
	p.logLevel = logLevel

	for i := 0; i <= logLevel; i++ {
		p.loggers[i] = log.New(p.handler, p.prefix+logLevelToName[i]+" ", p.flags)
	}
	nullwriter,_:= NewNullHandler()
	for i := logLevel + 1; i < NR_LEVELS; i++ {
		p.loggers[i] = log.New(nullwriter, p.prefix+logLevelToName[i]+" ", p.flags)
	}
}

func (p *logger)  SetLevel(level int) {
	p.Set(level,p.flags,p.handler)
}

func (p *logger)  SetLevelFlag(level int, flags int) {
	p.Set(level,flags,p.handler)
}


func newStdHandler() *StreamHandler {
	p,_:= NewStreamHandler(os.Stdout)
	return p
}

//new a default logger with specified handler and flags: Ltime|Lfile|Llevel
func NewDefault(handler Handler) *logger {
	return NewLogger(handler,"",LEVEL_SILENT,log.LstdFlags)
}

var Global = NewDefault(newStdHandler())

func Set(level int,flags int, handler Handler){
	Global.Set(level,flags,handler)
}

func SetLevel(level int) {
	SetLevelFlag(level,Global.flags)
}

func SetLevelFlag(level int, flags int) {
	Set(level,flags,Global.handler)
}


func GetLevel() int {
	return Global.logLevel
}

var CallDepth int = 2;

func Debug(v ...interface{}) {
	if(GetLevel() < LEVEL_DEBUG){
		return
	}
	Global.loggers[LEVEL_DEBUG].Output(CallDepth, fmt.Sprintln(v...))
}

func Debugf(format string, v ...interface{}) {
	if(GetLevel() < LEVEL_DEBUG){
		return
	}
	fmt.Println(Global.handler)
	Global.loggers[LEVEL_DEBUG].Output(CallDepth, fmt.Sprintf(format, v...))
}

func Config(v ...interface{}) {
	if(GetLevel() < LEVEL_CONFIG){
		return
	}
	Global.loggers[LEVEL_CONFIG].Output(CallDepth, fmt.Sprintln(v...))
}

func Configf(format string, v ...interface{}) {
	if(GetLevel() < LEVEL_CONFIG){
		return
	}
	Global.loggers[LEVEL_CONFIG].Output(CallDepth, fmt.Sprintf(format, v...))
}

func Info(v ...interface{}) {
	if(GetLevel() < LEVEL_INFO){
		return
	}
	Global.loggers[LEVEL_INFO].Output(CallDepth, fmt.Sprintln(v...))
}

func Infof(format string, v ...interface{}) {
	if(GetLevel() < LEVEL_INFO){
		return
	}
	Global.loggers[LEVEL_INFO].Output(CallDepth, fmt.Sprintf(format, v...))
}

func Warn(v ...interface{}) {
	if(GetLevel() < LEVEL_WARN){
		return
	}
	Global.loggers[LEVEL_WARN].Output(CallDepth, fmt.Sprintln(v...))
}

func Warnf(format string, v ...interface{}) {
	if(GetLevel() < LEVEL_WARN){
		return
	}
	Global.loggers[LEVEL_WARN].Output(CallDepth, fmt.Sprintf(format, v...))
}

func Alert(v ...interface{}) {
	if(GetLevel() < LEVEL_ALERT){
		return
	}
	Global.loggers[LEVEL_ALERT].Output(CallDepth, fmt.Sprintln(v...))
}

func Alertf(format string, v ...interface{}) {
	if(GetLevel() < LEVEL_ALERT){
		return
	}
	Global.loggers[LEVEL_ALERT].Output(CallDepth, fmt.Sprintf(format, v...))
}

func Error(v ...interface{}) {
	if(GetLevel() < LEVEL_ERROR){
		return
	}
	Global.loggers[LEVEL_ERROR].Output(CallDepth, fmt.Sprintln(v...))
}

func Errorf(format string, v ...interface{}) {
	if(GetLevel() < LEVEL_ERROR){
		return
	}
	Global.loggers[LEVEL_ERROR].Output(CallDepth, fmt.Sprintf(format, v...))
}

func Fatal(v ...interface{}) {
	if(GetLevel() < LEVEL_FATAL){
		return
	}
	Global.loggers[LEVEL_FATAL].Output(CallDepth, fmt.Sprintln(v...))
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	if(GetLevel() < LEVEL_FATAL){
		return
	}
	Global.loggers[LEVEL_FATAL].Output(CallDepth, fmt.Sprintf(format, v...))
	os.Exit(1)
}
