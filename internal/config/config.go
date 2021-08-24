package config

import (
	"log"
	"reflect"

	"github.com/spf13/viper"
)

type Config struct {
	MongoURI               string `mapstructure:"MONGO_URI"`
	ServerPort             uint16 `mapstructure:"PORT"`
	Server                 string `mapstructure:"SERVER"`
	SessionSecret          string `mapstructure:"SESSION_SECRET"`
	GoogleClientID         string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret     string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GoogleOauthRedirectURL string `mapstructure:"GOOGLE_OAUTH_REDIRECT_URL"`
	LineLoginClientID      string `mapstructure:"LINE_LOGIN_CLIENT_ID"`
	LineLoginClientSecret  string `mapstructure:"LINE_LOGIN_CLIENT_SECRET"`
	LineLoginRedirectURL   string `mapstructure:"LINE_LOGIN_REDIRECT_URL"`
	JwtSecret              string `mapstructure:"JWT_SECRET"`
	RabbitMQUserName       string `mapstructure:"RABBIT_MQ_USER_NAME"`
	RabbitMQPassword       string `mapstructure:"RABBIT_MQ_PASSWORD"`
	RabbitMQHost           string `mapstructure:"RABBIT_MQ_HOST"`
	RabbitMQPort           string `mapstructure:"RABBIT_MQ_PORT"`
}

func LoadConfig(path, configName, configType string) (*Config, error) {
	config := &Config{}

	viper.AddConfigPath(path)
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)

	// environment variable will override config if exists
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("Read config file failed. Use environment variables instead.")
			t := reflect.TypeOf(config)
			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}
			for i := 0; i < t.NumField(); i++ {
				key := t.Field(i).Tag.Get("mapstructure")
				viper.BindEnv(key)
			}
		} else {
			return nil, err
		}
	}

	err := viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
