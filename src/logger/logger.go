package logger

import (
	"github.com/Sirupsen/logrus"
	"log"
	"os"
)

var SystemLog = logrus.New()

func init() {
	file, err := os.Create("./var/system.log")
	if err != nil {
		log.Fatal(err)
	}
	SystemLog.Formatter = new(logrus.JSONFormatter)
	SystemLog.Out = file
}
