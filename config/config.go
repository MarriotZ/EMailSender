package config

import "github.com/spf13/viper"

var Conf *config
var Secret *config

type config struct {
	viper *viper.Viper
}

func init() {
	Conf = &config{
		viper: getConf("conf", "config/conf"),
	}
	Secret = &config{
		viper: getConf("secret", "config/secret"),
	}
}

func getConf(configName, configPath string) *viper.Viper {
	v := viper.New()
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)
	v.SetConfigName(configName)
	v.ReadInConfig()
	return v
}

func (c *config) GetString(key string) string {
	return c.viper.GetString(key)
}
func (c *config) Get(key string) interface{} {
	return c.viper.Get(key)
}
