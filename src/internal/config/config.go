package config

import (
	"fmt"
	"reflect"

	"github.com/spf13/viper"
)

type APPConfig struct {
	APP_NAME    string `mapstructure:"APP_NAME"`
	API_VERSION string `mapstructure:"API_VERSION"`
	APP_ENV     string `mapstructure:"APP_ENV"`
	APP_PORT    string `mapstructure:"APP_PORT"`
	APP_HOST    string `mapstructure:"APP_HOST"`
	JWT_SECRET  string `mapstructure:"JWT_SECRET"`
}

type DatabaseConfig struct {
	DATABASE_URL      string `mapstructure:"DATABASE_URL"`
	DATABASE_PORT     string `mapstructure:"DATABASE_PORT"`
	DATABASE_HOST     string `mapstructure:"DATABASE_HOST"` // Added Host
	DATABASE_NAME     string `mapstructure:"DATABASE_NAME"`
	DATABASE_USER     string `mapstructure:"DATABASE_USER"`
	DATABASE_PASSWORD string `mapstructure:"DATABASE_PASSWORD"`
	USE_SSL           string `mapstructure:"USE_SSL"`
}

type BaseConfig struct {
	APPConfig      `mapstructure:",squash"`
	DatabaseConfig `mapstructure:",squash"`
}

func LoadConfig() (*BaseConfig, error) {
	v := viper.New()

	v.AutomaticEnv()

	v.SetConfigFile(".env")
	v.SetConfigType("env")
	_ = v.ReadInConfig() // Ignore error if .env doesn't exist

	BindStructConfigs(v, BaseConfig{})

	var config BaseConfig

	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	// If DATABASE_URL is empty, we build it from components
	if config.DATABASE_URL == "" {
		sslMode := "disable"
		if config.USE_SSL == "true" {
			sslMode = "require"
		}

		// Format: postgres://user:password@host:port/dbname?sslmode=...
		config.DATABASE_URL = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			config.DATABASE_USER,
			config.DATABASE_PASSWORD,
			config.DATABASE_HOST,
			config.DATABASE_PORT,
			config.DATABASE_NAME,
			sslMode,
		)
	}

	return &config, nil
}

func BindStructConfigs(v *viper.Viper, i interface{}) {
	val := reflect.ValueOf(i)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("mapstructure")

		// Handle nested structs (like DatabaseConfig inside BaseConfig)
		if field.Type.Kind() == reflect.Struct {
			BindStructConfigs(v, reflect.New(field.Type).Interface())
			continue
		}

		if tag != "" && tag != ",squash" {
			v.BindEnv(tag)
		}
	}
}
