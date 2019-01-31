package lib

import (
	"github.com/spf13/viper"
)

type Config struct {
	DspUrlOne string `yaml:"dspurlOne"`
	DspUrlTwo string `yaml:"dspurlTwo"`
	Port      int    `yaml:"port"`
}

//Load configration based on path and filename
// return Config
func LoadConfig(path string, fileName string) (*Config, error) {
	var conf *Config = &Config{}

	viper.AddConfigPath(path)
	viper.SetConfigName(fileName)
	err := viper.ReadInConfig()

	if err != nil {
		return conf, err
	}

	err = viper.Unmarshal(conf)
	if err != nil {
		return conf, err
	}

	return conf, nil
}
