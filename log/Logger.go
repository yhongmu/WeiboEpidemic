package log

import (
	"io"
	"log"
	"os"
	"sync"
	"time"
)

var loggerInstance *logger
var once sync.Once

func GetLog() *logger {
	once.Do(func() {
		loggerInstance = &logger{}
		loggerInstance.Init("main")
	})
	return loggerInstance
}

type logger struct {
	Trace *log.Logger
	Info *log.Logger
	Warning *log.Logger
	Error *log.Logger
}

func (l *logger) Init(name string) {
	file := "./log/file/" + time.Now().Format("2006-01-02_15:04:05") + "_" + name + "_log" + ".txt"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	l.Trace = log.New(io.MultiWriter(logFile, os.Stdout),
		"[TRACE]: ",
		log.Ldate | log.Ltime)

	l.Info = log.New(io.MultiWriter(logFile, os.Stdout),
		"[INFO]: ",
		log.Ldate | log.Ltime)

	l.Warning = log.New(io.MultiWriter(logFile, os.Stdout),
		"[WARNING]: ",
		log.Ldate | log.Ltime | log.Lshortfile)

	l.Error = log.New(io.MultiWriter(logFile, os.Stdout),
		"[ERROR]: ",
		log.Ldate | log.Ltime | log.Lshortfile)
}

