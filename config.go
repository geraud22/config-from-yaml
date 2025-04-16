package main

import (
	"fmt"
	"io"
	"sync"

	"github.com/go-playground/validator"
	"github.com/spf13/viper"
)

var (
	once     sync.Once
	instance *viper.Viper
)

func LoadConfig[T any](in io.Reader, filetype string) (T, error) {
	var result T
	v := viper.New()
	v.SetConfigType(filetype)
	if err := v.ReadConfig(in); err != nil {
		return result, fmt.Errorf("error reading config: %v", err)
	}
	if err := v.Unmarshal(&result); err != nil {
		return result, fmt.Errorf("error unmarshalling config into provided struct: %v", err)
	}
	validate := validator.New()
	if err := validate.Struct(&result); err != nil {
		return result, fmt.Errorf("error validating config: %v", err)
	}
	return result, nil
}

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
