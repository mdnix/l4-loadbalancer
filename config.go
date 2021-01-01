package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type Config struct {
	Global struct {
		Output string `mapstructure:"output"`
	} `mapstructure:"global"`
	Services []struct {
		Name      string   `mapstructure:"name"`
		Bind      string   `mapstructure:"bind"`
		Backends  []string `mapstructure:"backends"`
		Algorithm string   `mapstructure:"algorithm"`
	} `mapstructure:"services"`
}

var config Config

func initConfig() {
	if configFile != "" {
		// Use config file from flag.
		viper.SetConfigType("yaml")
		viper.SetConfigFile(configFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		cfgPath := fmt.Sprintf("%s/%s", home, ".services")

		// Search config in home directory with name ".clusters".
		viper.SetConfigType("yaml")
		viper.SetConfigFile(cfgPath)
	}
	// Read the config file
	if err := viper.ReadInConfig(); err == nil {
		err = viper.Unmarshal(&config)
	} else {
		fmt.Println("Unable to read config file: ", err)
		os.Exit(1)
	}
}
