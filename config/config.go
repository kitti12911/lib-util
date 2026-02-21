package config

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/kitti12911/lib-util/validator"
)

// Load reads a config file into T, overrides fields that have an `env` struct tag
//
// Example:
//
//	type Config struct {
//	    Port     int           `mapstructure:"port"     env:"PORT"      validate:"required,gte=1"`
//	    Host     string        `mapstructure:"host"     env:"HOST"      validate:"required"`
//	    LogLevel string        `mapstructure:"logLevel" env:"LOG_LEVEL" validate:"oneof=debug info warn error"`
//	    Timeout  time.Duration `mapstructure:"timeout"  env:"TIMEOUT"`
//	}
func Load[T any](path string) (*T, error) {
	v := viper.New()
	if path != "" {
		v.SetConfigFile(path)

		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("config: read file: %w", err)
		}
	}

	bindEnvs(v, reflect.TypeFor[T](), "")
	cfg := new(T)

	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("config: unmarshal: %w", err)
	}

	val := validator.New()
	if err := val.Validate(cfg); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	return cfg, nil
}

func bindEnvs(v *viper.Viper, t reflect.Type, prefix string) {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return
	}

	for field := range t.Fields() {
		tag := strings.SplitN(field.Tag.Get("mapstructure"), ",", 2)[0]

		if tag == "" || tag == "-" {
			tag = strings.ToLower(field.Name)
		}

		key := tag
		if prefix != "" {
			key = prefix + "." + tag
		}

		fieldType := field.Type
		if fieldType.Kind() == reflect.Pointer {
			fieldType = fieldType.Elem()
		}

		if fieldType.Kind() == reflect.Struct && fieldType != reflect.TypeFor[time.Time]() {
			bindEnvs(v, fieldType, key)
			continue
		}

		if envKey := field.Tag.Get("env"); envKey != "" {
			v.BindEnv(key, envKey)
		}
	}
}
