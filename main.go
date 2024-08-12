package main

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

var (
	once     sync.Once
	instance *viper.Viper
)

func NewConfig(filename string) *viper.Viper {
	once.Do(func() {
		v := viper.New()
		v.SetConfigName(filename)
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		err := v.ReadInConfig()
		if err != nil {
			fmt.Printf("Error reading config file: %s\n", err)
		}
		instance = v
	})
	return instance
}

func Get(filename string) *viper.Viper {
	if instance == nil {
		return NewConfig(filename)
	}
	return instance
}
