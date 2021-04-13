package config

import "github.com/spf13/viper"

type Config struct {
	MongoURI               string `mapstructure:"MONGO_URI"`
	ServerPort             uint16 `mapstructure:"SERVER_PORT"`
	GoogleClientID         string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret     string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GoogleOauthRedirectURL string `mapstructure:"GOOGLE_OAUTH_REDIRECT_URL"`
}

func LoadConfig(path, configName, configType string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)

	// viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
