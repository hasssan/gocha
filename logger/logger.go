package logger

import (
	log "github.com/Sirupsen/logrus"
)

func SetLogLevel(lvl string) {
	log.SetLevel(getLogLevel(lvl))
}

// getLogLevel retrieves the log.Level corresponding to the passed string.
func getLogLevel(l string) log.Level {
	lvl, err := log.ParseLevel(l)
	if err != nil {
		log.WithFields(log.Fields{
			"passed":  l,
			"default": "fatal",
		}).Warn("Log level is not valid, fallback to default level")

		return log.FatalLevel
	}

	return lvl
}
