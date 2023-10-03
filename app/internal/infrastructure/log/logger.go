package log

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
)

type Logger struct {
	infoLog    *log.Logger
	debugLog   *log.Logger
	warningLog *log.Logger
	errorLog   *log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		infoLog:    log.New(os.Stdout, "INFO\t", log.Ltime),
		debugLog:   log.New(os.Stdout, "DEBUG\t", log.Ltime),
		warningLog: log.New(os.Stderr, color.YellowString("WARING")+"\t", log.Ltime),
		errorLog:   log.New(os.Stderr, color.RedString("ERROR")+"\t", log.Ltime),
	}
}

func (l *Logger) Info(message string, parts ...interface{}) {
	l.infoLog.Println(message, parts)
}

func (l *Logger) Debug(message string, parts ...interface{}) {
	l.debugLog.Println(message, parts)
}

func (l *Logger) Warning(message string, parts ...interface{}) {
	l.warningLog.Println(color.YellowString(fmt.Sprint(message, parts)))
}

func (l *Logger) Error(message string, parts ...interface{}) {
	l.errorLog.Println(color.RedString(fmt.Sprint(message, parts)))
}

func (l *Logger) Fatal(message string, parts ...interface{}) {
	l.errorLog.Fatalln(color.HiRedString(fmt.Sprint(message, parts)))
}
