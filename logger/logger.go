package logger

const PREFIX = "live."
const TAG_DEBUG = "debug"
const TAG_INFO = "info"
const TAG_ERROR = "error"

type Logger interface {
    Start()
    Debug(...interface{})
    Info(...interface{})
    Error(...interface{})
}

var l Logger
func Debug(logs ...interface{}) {
    l.Debug(logs...)
}

func Info(logs ...interface{}) {
    l.Info(logs...)
}

func Error(logs... interface{}) {
    l.Error(logs)
}

func Run(config LoggerConfig) {
    if config.Env == "production" {
        l = FluentdLogger{ config.Fluentd.Name, config.Fluentd.Host, config.Fluentd.Port }
    } else {
        l = LocalLogger{}
    }

    l.Start()
}
