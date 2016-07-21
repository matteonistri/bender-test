package main

import "gopkg.in/ini.v1"

var logContextConfig LoggerContext

type config struct {
	generalLogLevel int
	daemonLogLevel  int
	statusName      string
}

func ConfigInit(cfgFileName string) config {
	// init LogContex
	logContextConfig = LoggerContext{
		level: 1,
		name:  "CONFIG"}

	// set defaults
	conf := config{
		generalLogLevel: 3,
		daemonLogLevel:  3,
		statusName:      "bender"}

	// attempt to load config file
	cfgFileName += ".cfg"
	cfg, err := ini.Load(cfgFileName)
	if err != nil {
		LogWar(logContextConfig, "No config file found: using defaults")
		return conf
	}

	// parse config file and set config
	sections := cfg.Sections()
	for _, s := range sections {
		switch s.Name() {
		case "general":
			if s.HasKey("loglevel") {
				v, err := s.Key("loglevel").Int()
				if err != nil {
					LogWar(logContextConfig, "cannot parse key log/level: %s", err)
				} else {
					if v >= 0 && v <= 3 {
						conf.generalLogLevel = v
					} else {
						LogWar(logContextConfig, "invalid value for key log/level: %d", v)
					}
				}
			}
		case "daemon":
			if s.HasKey("loglevel") {
				v, err := s.Key("loglevel").Int()
				if err != nil {
					LogWar(logContextConfig, "cannot parse key log/level: %s", err)
				} else {
					if v >= 0 && v <= 3 {
						conf.daemonLogLevel = v
					} else {
						LogWar(logContextConfig, "invalid value for key log/level: %d", v)
					}
				}
			}
		case "status":
			if s.HasKey("servername") {
				conf.statusName = s.Key("servername").String()
			}
		}
	}

	return conf
}
