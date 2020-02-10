package logger

import (
	"log"
)

func Error(args ...interface{}) {
	log.Println("[ERRO] ", args)
}

func Warn(args ...interface{}){
	log.Println("[WARN] ", args)
}

func Info(args ...interface{}){
	log.Println("[INFO] ", args)
}