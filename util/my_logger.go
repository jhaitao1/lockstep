package util

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

//Log 全局log 带颜色的
var Log *logrus.Logger

var once sync.Once

func init() {
	once.Do(func() {
		log := logrus.New()
		log.Formatter = &logrus.TextFormatter{
			DisableColors: false,
			ForceColors:   true,
		}
		log.Level = logrus.TraceLevel
		log.Out = os.Stdout
		// log.Debug("你好")
		// log.Warn("你好")
		// log.Error("你好")
	})

}
