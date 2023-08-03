// Package config contains logic pertaining to application configuration.
package config

import "github.com/spf13/viper"

const defaultPort int = 8080

// InitConfig will initialize the Viper configuration and bind environment variables for the application.
func InitConfig() error {
	viper.AutomaticEnv()

	// Bind general variables
	viper.SetDefault("service.PORT", defaultPort)
	err := viper.BindEnv("service.port", "PORT")
	if err != nil {
		return err
	}

	viper.SetDefault("files.tmp", "/tmp")

	// Bind postgres environment variables
	err = viper.BindEnv("postgres.username", "POSTGRES_USERNAME")
	if err != nil {
		return err
	}

	err = viper.BindEnv("postgres.password", "POSTGRES_PASSWORD")
	if err != nil {
		return err
	}

	err = viper.BindEnv("postgres.db", "POSTGRES_DB")
	if err != nil {
		return err
	}

	viper.SetDefault("postgres.host", "localhost")
	err = viper.BindEnv("postgres.host", "POSTGRES_HOST")

	return err
}
