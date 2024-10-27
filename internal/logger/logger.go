package logger

import "log"

var Debug bool

func DebugPrintf(format string, v ...interface{}) {
	if Debug {
		log.Printf(format, v...)
	}
}

func Println(v ...interface{}) {
	log.Println(v...)
}

func Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}
