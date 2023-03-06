package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
)

const (
	configDir = "./config"
	dbEnvFile = "db.env"
)

type StanConfig struct {
	ClusterID string
	ClientID  string
	ChannelID string
}

type PostgresConfig struct {
	Host     string
	Port     string
	DbName   string
	Username string
	Password string
}

type App struct {
	Host string
	Port string
	Pg   PostgresConfig
	Stan StanConfig
}

func InitConfig() (App, error) {
	viper.AddConfigPath(configDir)
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		return App{}, err
	}

	err = godotenv.Load(dbEnvFile)
	if err != nil {
		return App{}, err
	}

	pg := PostgresConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		DbName:   viper.GetString("db.name"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
	}

	stan := StanConfig{
		ClusterID: viper.GetString("stan.cluster_id"),
		ClientID:  viper.GetString("stan.client_id"),
		ChannelID: viper.GetString("stan.channel_id"),
	}

	app := App{
		Host: viper.GetString("app.host"),
		Port: viper.GetString("app.port"),
		Pg:   pg,
		Stan: stan,
	}

	return app, nil
}
