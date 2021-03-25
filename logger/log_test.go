package logger

import (
	"os"
	"log"
	"testing"
)

func TestLog(t *testing.T) {
	SetLevel(LEVEL_DEBUG)
	SetLevelFlag(LEVEL_DEBUG,log.Lshortfile|log.Ldate|log.Ltime)

	handler,_ := NewStreamHandler(os.Stdout)
	Set(LEVEL_DEBUG,log.Lshortfile|log.Ldate|log.Ltime,handler)

	Debug("default global log test")

	Global = NewDefault(handler)
	Global.SetLevel(LEVEL_DEBUG)
	Global.SetLevelFlag(LEVEL_DEBUG,log.Lshortfile|log.Ldate|log.Ltime)
	Global.Set(LEVEL_DEBUG,log.Lshortfile|log.Ldate|log.Ltime,handler)
	Global.Debug("default global log test")
	Global.Close()
}