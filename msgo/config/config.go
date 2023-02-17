package config

import (
	"flag"
	"github.com/BurntSushi/toml"
	msLog "github.com/jinouy/msgo/log"
	"os"
)

var Conf = &MsConfig{
	logger: msLog.Default(),
}

type MsConfig struct {
	logger   *msLog.Logger
	Log      map[string]any
	Pool     map[string]any
	Template map[string]any
}

func init() {
	loadToml()
}

func loadToml() {
	configFile := flag.String("conf", "conf/app.toml", "app config file")
	flag.Parse()
	if _, err := os.Stat(*configFile); err != nil {
		Conf.logger.Info("conf/app.toml file not load, because not exist")
		return
	}
	_, err := toml.DecodeFile(*configFile, Conf)
	if err != nil {
		Conf.logger.Info("conf/app.toml decode fail check format")
		return
	}

}
