package logruslokihook

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type LokiPayload struct {
	Streams []LokiStreams `json:"streams"`
}

type LokiStreams struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

type LogrusLokiConfig struct {
	Url       string
	Labels    map[string]string
	Formatter logrus.Formatter
}

type LogrusLokiHook struct {
	Config LogrusLokiConfig
}

func NewLogrusLoki(config LogrusLokiConfig) (logrus.Hook, error) {
	hook := LogrusLokiHook{Config: config}

	return &hook, nil
}

func (hook *LogrusLokiHook) Levels() []logrus.Level {
	return logrus.AllLevels[:logrus.DebugLevel]
}

func (hook *LogrusLokiHook) Fire(entry *logrus.Entry) (err error) {
	b, err := hook.Config.Formatter.Format(entry)
	if err != nil {
		return err
	}

	go hook.Push(b)

	return
}

func (hook *LogrusLokiHook) Push(b []byte) {
	payload := LokiPayload{
		Streams: []LokiStreams{
			{
				Stream: hook.Config.Labels,
				Values: [][]string{{strconv.FormatInt(time.Now().UnixNano(), 10), string(b)}},
			},
		},
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", hook.Config.Url, bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")

	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return
	}
}
