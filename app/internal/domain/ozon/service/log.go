package service

type Log interface {
	Error(message string, parts ...interface{})
	Fatal(message string, parts ...interface{})
}
