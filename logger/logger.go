package logger

import (
	"log"
)

func init(){
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func Error(args ...interface{}) {
	log.Println("[ERRO]", args)
}

func Warn(args ...interface{}){
	log.Println("[WARN]", args)
}

func Info(args ...interface{}){
	log.Println("[INFO]", args)
}

func Debug(args ...interface{}){
	log.Println("[DEBUG]", args)
}