package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DBConfig       DBConfig       `yaml:"dbConfig"`
	JWTSecret      string         `yaml:"jwtSecret"`
	ServerConfig   ServerConfig   `yaml:"serverConfig"`
	RabbitMQConfig RabbitMQConfig `yaml:"rabbitMQConfig"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port"`
}

type ServerConfig struct {
	Port         int   `yaml:"port"`
	ReadTimeout  int64 `yaml:"readTimeout"`
	WriteTimeout int64 `yaml:"writeTimeout"`
	IdleTimeout  int64 `yaml:"idleTimeout"`
}

type RabbitMQConfig struct {
	Host                  string `yaml:"host"`
	Port                  int    `yaml:"port"`
	User                  string `yaml:"user"`
	Password              string `yaml:"password"`
	StockDecreaseExchange string `yaml:"stockDecreaseExchange"`
	StockDecreaseQueue    string `yaml:"stockDecreaseQueue"`
	StockStatusExchange   string `yaml:"stockStatusExchange"`
	StockStatusQueue      string `yaml:"stockStatusQueue"`
}

func Read() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$PWD/config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/config")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("fatal error unmarshalling config: %w", err))
	}

	return &config
}
