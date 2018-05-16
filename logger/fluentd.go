package logger

import (
    "github.com/fluent/fluent-logger-golang/fluent"
    "time"
    "log"
)

type FluentdLogger struct {
    Name string
    Host string
    Port int
}

var f *fluent.Fluent

func (l FluentdLogger) Start() {
    config := fluent.Config{
        FluentHost: l.Host,
        FluentPort: l.Port,
        TagPrefix: PREFIX + l.Name,
        Timeout: 3*time.Second }
    client, err := fluent.New(config)
    if err != nil {
        log.Println(err)
    }

    f = client
}

func (l FluentdLogger) Debug(msg ...interface{}) {

}

func (l FluentdLogger) Info(msg ...interface{}) {
    err := f.Post(TAG_INFO, msg)
    if err != nil {
        log.Println(err.Error())
    }
}

func (l FluentdLogger) Error(msg ...interface{}) {
    err := f.Post(TAG_ERROR, msg)
    if err != nil {
        log.Println(err.Error())
    }
}