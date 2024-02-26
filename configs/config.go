package configs

import "github.com/spf13/viper"

var Cfg *conf

type conf struct {
	WebServerPort     string `mapstructure:"WEB_SERVER_PORT"`
	RateLimitIp       int    `mapstructure:"RATE_LIMIT_IP"`
	RateLimitToken    int    `mapstructure:"RATE_LIMIT_TOKEN"`
	RateLimitQtdBlock int    `mapstructure:"RATE_LIMIT_QTD_BLOCK"`
	BdAddress         string `mapstructure:"BD_ADDRESS"`
	BdPassword        string `mapstructure:"BD_PASSWORD"`
}

func LoadConfig(path string) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&Cfg)
	if err != nil {
		panic(err)
	}

	//return Cfg
}
