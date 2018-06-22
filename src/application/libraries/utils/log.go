package utils

import (
	"log"
	"os"
)

var logger = log.New(os.Stderr, "\r\n", log.Ldate|log.Ltime|log.Lshortfile)

func Logger() *log.Logger{
	return logger
}

func Mylog(logs ...interface{}) {

	logger.Println(logs...)

}

func Error(logs ...interface{}) {

	logger.Println(logs)

}
