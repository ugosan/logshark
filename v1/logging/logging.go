package logging

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
	"github.com/ugosan/logshark/v1/config"
	"io"
	"os"
	"sync"
)

type logmanager struct {
	logger *log.Logger
}

var (
	singleton *logmanager
	once      sync.Once
)

func GetManager() *logmanager {

	once.Do(func() {
		singleton = &logmanager{logger: log.New()}
		singleton.logger.SetOutput(io.Discard)
	})

	return singleton
}

func (lm *logmanager) Log(s interface{}) {
	lm.logger.Println(s)
}

func (lm *logmanager) InitLogger(config config.Config) {

	f, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	lm.logger.SetFormatter(&nested.Formatter{
		HideKeys:    true,
		FieldsOrder: []string{"component", "category"},
	})

	lm.logger.SetOutput(f)

}
