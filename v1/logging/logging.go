package logging

import (
	"io/ioutil"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/ugosan/logshark/v1/config"
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
	})

	return singleton
}

func (sm *logmanager) Log(s interface{}) {
	sm.logger.Println(s)
}

func (sm *logmanager) InitLogger(config config.Config) {

	if config.LogFile == "/dev/null" {
		sm.logger.SetOutput(ioutil.Discard)
	} else {
		f, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}

		sm.logger.SetOutput(f)
	}

}
