package main

import "gopkg.in/ini.v1"

var logContextConfig LoggerContext

type ConfigInterface interface {
	Get(section string, key string, defValue string) string
	GetLogLevel(section string, defValue int) int
}

type ConfigModule struct {
	conf *ini.File
}

// ConfigInit initializes the config module
func ConfigInit(cm *ConfigModule, filename string) {
	// init LogContex
	logContextConfig = LoggerContext{
		level: 1,
		name:  "CONFIG"}

	// attempt to load config file
	filename += ".cfg"
	var err error
	cm.conf, err = ini.Load(filename)
	if err != nil {
		LogWar(logContextConfig, "No config file found: using defaults")
		cm.conf = ini.Empty()
		return
	}

	return
}

func (cm *ConfigModule) Get(section string, key string, defValue string) string {
	value, err := cm.conf.Section(section).GetKey(key)
	if err != nil {
		LogWar(logContextConfig, "No value for key %s found: using default %s", key, defValue)
		return defValue
	}

	return value.String()
}

func (cm *ConfigModule) GetLogLevel(section string, defValue int) int {
	value, err := cm.conf.Section(section).GetKey("loglevel")
	if err != nil {
		LogWar(logContextConfig, "No value for loglevel found: using default %s", defValue)
		return defValue
	}

	v, err := value.Int()
	if err != nil {
		LogWar(logContextConfig, "Cannot parse loglevel value: using default. %s", err)
		return defValue
	}

	if v < 0 || v > 3 {
		LogWar(logContextConfig, "Invalid value for key loglevel %d: using default", value)
		return defValue
	}

	return v
}
