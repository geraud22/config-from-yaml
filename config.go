package cfy

import (
	"fmt"
	"io"

	validatorModule "github.com/go-playground/validator"
	"github.com/spf13/viper"
)

type Validator interface {
	Struct(any) error
}

func LoadConfig[T any](in io.Reader, filetype string, validator Validator) (T, error) {
	var result T
	v := viper.New()
	v.SetConfigType(filetype)
	if err := v.ReadConfig(in); err != nil {
		return result, fmt.Errorf("error reading config: %v", err)
	}
	if err := v.Unmarshal(&result); err != nil {
		return result, fmt.Errorf("error unmarshalling config into provided struct: %v", err)
	}
	if validator == nil {
		validator = validatorModule.New()
	}
	if err := validator.Struct(&result); err != nil {
		return result, fmt.Errorf("error validating config: %v", err)
	}
	return result, nil
}
