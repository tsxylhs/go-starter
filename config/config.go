package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

var confRoot *string
var Config *viper.Viper

const (
	PREFIX_REGISTRY_CONSUL = "registry.consul"
	PREFIX_DB              = "db"
)

type Pair struct {
	Key    string
	Target interface{}
}

var (
	once sync.Once
)

func Configer(names ...string) *viper.Viper {
	once.Do(func() {
		if Config == nil {
			Config = viper.New()
			fn := "config"
			if len(names) > 0 {
				fn = names[0]
			}
			Config.SetConfigName(fn)
			Config.AddConfigPath("$HOME/.lncios.cn/")
			Config.AddConfigPath(".")
			err := Config.ReadInConfig()
			if err != nil {
				log.Fatal(err)
			}

		}
	})
	return Config
}
