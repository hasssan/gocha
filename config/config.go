// Package config wraps the configuration setup and getters.
package config

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	ConfigFileName = ".gocha"
	ConfigFileType = "yaml"
	ConfigFilePath = "$HOME"
)

func init() {
	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType(ConfigFileType)
	viper.AddConfigPath(ConfigFilePath)
	err := viper.ReadInConfig()
	if err != nil {
		log.Debug(err)
	}
}

func Get(c string) interface{} {
	return viper.Get(c)
}

func GetCliOrConfig(c string, cli interface{}) interface{} {
	val := cli
	if cli == "" {
		if strings.Contains(c, "/") {
			s := strings.Split(c, "/")
			return viper.GetStringMap(s[0])[s[1]]
		} else {
			return Get(c)
		}
	}
	return val
}

func GetCliOrConfigString(c string, cli string) string {
	if val, ok := GetCliOrConfig(c, cli).(string); ok {
		return val
	}
	return ""
}
