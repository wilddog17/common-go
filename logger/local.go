package logger

import (
    "log"
)

type LocalLogger struct {

}

func (l LocalLogger) Start() {

}

func (l LocalLogger) Debug(msg ...interface{}) {
    m := append([]interface{} { TAG_DEBUG }, msg)
    log.Println(m...)
}

func (l LocalLogger) Info(msg ...interface{}) {
    m := append([]interface{} { TAG_INFO }, msg)
    log.Println(m...)
}

func (l LocalLogger) Error(msg ...interface{}) {
    m := append([]interface{} { TAG_ERROR }, msg)
    log.Println(m...)
}