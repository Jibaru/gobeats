package config

import (
	"github.com/spf13/viper"
)

type AppConfig struct {
	GoogleDriveRootFolderKey string
	GoogleDriveDriveApiKey   string
	UseAutoPlay              bool
	InitialVolume            int
}

func NewAppConfig() (*AppConfig, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return &AppConfig{
		GoogleDriveRootFolderKey: viper.GetString("google_drive.root_folder_key"),
		GoogleDriveDriveApiKey:   viper.GetString("google_drive.api_key"),
		UseAutoPlay:              viper.GetBool("player.autoplay"),
		InitialVolume:            viper.GetInt("player.initial_volume"),
	}, nil
}
