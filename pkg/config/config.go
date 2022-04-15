package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	ConfigPathKey = "CONFIG_PATH"
)

type Config struct {
	Hostname       string `mapstructure:"hostname"`
	ServerURL      string `mapstructure:"server_url"`
	ServerWSURL    string `mapstructure:"server_ws_url"`
	TeamName       string `mapstructure:"team_name"`
	ChannelLogName string `mapstructure:"channel_log_name"`
	AuthToken      string `mapstructure:"auth_token"`
}

func LoadConfig() (config Config, err error) {

	_ = viper.BindEnv("hostname", "HOSTNAME")
	_ = viper.BindEnv("server_url", "SERVER_URL")
	_ = viper.BindEnv("server_ws_url", "SERVER_WS_URL")
	_ = viper.BindEnv("team_name", "TEAM_NAME")
	_ = viper.BindEnv("channel_log_name", "CHANNEL_LOG_NAME")
	_ = viper.BindEnv("auth_token", "BOT_AUTH_TOKEN")

	configPath := os.Getenv(ConfigPathKey)
	if configPath == "" {
		zap.S().Info("no configuration file provided, defaulting to current directory")
		configPath = "."
	}
	viper.AddConfigPath(configPath)
	files, _ := os.ReadDir(configPath)
	for _, file := range files {
		fileName := file.Name()
		lastDotIndex := strings.LastIndex(fileName, ".")
		if lastDotIndex == -1 {
			zap.S().Info("File without extension will be ignored", "filename", fileName)
			continue
		}
		extFile := filepath.Ext(file.Name())
		if extFile != ".yaml" && extFile != ".yml" {
			zap.S().Info("File not in a yaml format, will be ignored", "filename", fileName)
			continue
		}
		viper.SetConfigName(fileName[:lastDotIndex])
		viper.SetConfigType("yaml")
		err = viper.MergeInConfig()
		if err != nil {
			return
		}
	}
	viper.AutomaticEnv()
	err = viper.Unmarshal(&config)
	if err != nil {
		return config, errors.Wrap(err, "failed to unmarshal config")
	}
	if config.ServerURL == "" {
		config.ServerURL = "https://" + config.Hostname
	}
	if config.ServerWSURL == "" {
		config.ServerWSURL = "wss://" + config.Hostname
	}
	return config, nil
}
