package hup

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var (
	debugLog = log.New(os.Stdout, "\033[32m[Debug]\033[0m ", log.LstdFlags|log.Lshortfile)
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile)
	infoLog  = log.New(os.Stdout, "\033[34m[info ]\033[0m ", log.LstdFlags|log.Lshortfile)
	loggers  = []*log.Logger{errorLog, infoLog}
	mu       sync.Mutex
)

var (
	Debug  = debugLog.Println
	DebugF = debugLog.Printf
	Error  = errorLog.Println
	ErrorF = errorLog.Printf
	Info   = infoLog.Println
	InfoF  = infoLog.Printf
)

const (
	ErrorLevel = iota
	DebugLevel
	InfoLevel
	Disabled
)

func SetLevel(level int) {
	mu.Lock()
	defer mu.Unlock()

	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}

	if ErrorLevel < level {
		errorLog.SetOutput(ioutil.Discard)
	}

	if DebugLevel < level {
		debugLog.SetOutput(ioutil.Discard)
	}

	if InfoLevel < level {
		infoLog.SetOutput(ioutil.Discard)
	}
}
