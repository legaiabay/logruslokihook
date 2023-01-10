package main

import (
	"github.com/legaiabay/logruslokihook"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	lokiHook, err := logruslokihook.NewLogrusLoki(logruslokihook.LogrusLokiConfig{
		Url:    "http://172.31.24.210:3100/loki/api/v1/push",
		Job:    "someJob",
		Source: "local abay",
		Formatter: &logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.999999999Z07:00",
		},
	})

	if err != nil {
		log.Fatalf("Failed to init hook: %v", err)
	}

	log.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.999999999Z07:00",
	})

	log.SetReportCaller(true)
	log.AddHook(lokiHook)
}

func main() {
	log.Info("info log")
	log.Error("error log")
	log.Fatal("fatal log")
}
