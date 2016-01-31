package utils

import (
	log "github.com/Sirupsen/logrus"
)

var Logger *log.Logger

func init() {
	Logger = log.New()
	Logger.Formatter = &log.JSONFormatter{}
}