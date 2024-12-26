package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
	"time"
)

type Config struct {
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	Server    ServerConfig    `mapstructure:"server"`
	App       AppConfig       `mapstructure:"app"`
	ShortCode ShortCodeConfig `mapstructure:"shortcode"`
}

func LoadConfig(filePath string) (*Config, error) {
	viper.SetConfigFile(filePath)
	viper.SetEnvPrefix("URL_SHORTENER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

type DatabaseConfig struct {
	Driver     string `mapstructure:"driver"`
	Host       string `mapstructure:"host"`
	Port       int    `mapstructure:"port"`
	User       string `mapstructure:"user"`
	Password   string `mapstructure:"password"`
	DBName     string `mapstructure:"dbname"`
	SSLMode    string `mapstructure:"ssl_mode"`
	MaxIdleCon int    `mapstructure:"max_idle_con"`
	MaxOpenCon int    `mapstructure:"max_open_con"`
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s", d.Driver, d.User, d.Password, d.Host, d.Port, d.DBName, d.SSLMode)
}

type RedisConfig struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type ServerConfig struct {
	Address      string        `mapstructure:"address"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
}

type AppConfig struct {
	BaseHost        string        `mapstructure:"base_host"`
	BasePort        string        `mapstructure:"base_port"`
	DefaultDuration time.Duration `mapstructure:"default_duration"`
	CleanUpInterval time.Duration `mapstructure:"cleanup_interval"`
}

type ShortCodeConfig struct {
	Length int `mapstructure:"length"`
}
