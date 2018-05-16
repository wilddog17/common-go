package logger

type FluentdConfig struct {
	Name string
	Host string
	Port int
}

type LoggerConfig struct {
	Env string
	Fluentd FluentdConfig
}