package config

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/spf13/viper"
)

var Conf ConfYaml

var defaultConf = []byte(`
jwt:
    key: "YOUR_JWT_KEY"
    age: 3000 # UNIX timestamp

log_level: "debug"

web_server:
    port: ":8080"

`)

type ConfYaml struct {
	JWT SectionJWT `yaml: "jwt"`
	WebServer SectionWebServer `yaml: "web_server"`
	LogLevel string `yaml: "loglevel"`
}

type SectionWebServer struct {
	Port string `yaml: "port"`
}

type SectionJWT struct {
	Key string `yaml: "key"`
	Age int    `yaml: "age"`
}

func LoadConf(confPath string) error {

	viper.SetConfigType("yaml")

	if confPath != "" {
		content, err := ioutil.ReadFile(confPath)

		if err != nil {
			return err
		}

		if err := viper.ReadConfig(bytes.NewBuffer(content)); err != nil {
			return err
		}
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		} else {
			// load default config
			if err := viper.ReadConfig(bytes.NewBuffer(defaultConf)); err != nil {
				return err
			}
		}
	}
	Conf.JWT.Key = viper.GetString("jwt.key")
	Conf.JWT.Age = viper.GetInt("jwt.age")
	Conf.LogLevel = viper.GetString("log_level")
	Conf.WebServer.Port = viper.GetString("web_server.port")
	return nil
}
