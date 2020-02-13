package logger

import (
	"log"
	"runtime/debug"
)

const (
	_ = iota
	DEBUG
	INFO
	WARN
	ERROR
)

var LEVEL = INFO

func init(){
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func Error(args ...interface{}) {
	if LEVEL > ERROR {
		return
	}

	log.Println("[ERRO]", args)
	debug.PrintStack()
}

func Warn(args ...interface{}){
	if LEVEL > WARN {
		return
	}

	log.Println("[WARN]", args)
	debug.PrintStack()
}

func Info(args ...interface{}){
	if LEVEL > INFO {
		return
	}

	log.Println("[INFO]", args)
}

func Debug(args ...interface{}){
	if LEVEL > DEBUG {
		return 
	}
	
	log.Println("[DEBUG]", args)
}