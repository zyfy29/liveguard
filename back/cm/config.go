package cm

import (
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"log"
	"reflect"
	"testing"
)

type Config struct {
	Pocket PocketConfig `json:"pocket"`
	AAI    AAIConfig    `json:"aai"`
	//Medium MediumConfig `json:"medium"`
}

type PocketConfig struct {
	Token    string `json:"token"`
	AppInfo  string `mapstructure:"app_info" json:"app_info"`
	Interval int    `json:"interval"`
}

type AAIConfig struct {
	Token string `json:"token"`
}

type MediumConfig struct {
	Token string `json:"token"`
}

func init() {
	configFile := "config"
	if testing.Testing() {
		configFile = "config.test"
	}

	viper.SetConfigName(configFile)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(GetProjectRoot())

	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("設定ファイルの読み込みエラー: %v", err)
	}

	if lo.IsEmpty(viper.GetInt("app.back_port")) ||
		lo.IsEmpty(viper.GetString("app.data_dir")) {
		log.Panicf("config of 'app' section can't be empty")
	}

	_ = GetConfig() // ensure the structs is correct
}

func GetConfig() *Config {
	config := &Config{}
	if err := viper.Unmarshal(&config); err != nil {
		log.Panic("failed to decode config!", err)
	}
	return config
}

func SetConfig(config *Config) int {
	var changedCount int
	setConfigRecursive(reflect.ValueOf(config), "", &changedCount)
	if err := viper.WriteConfig(); err != nil {
		log.Panic("failed to write config!", err)
	}
	return changedCount
}

func setConfigRecursive(v reflect.Value, prefix string, changedCount *int) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		key := jsonTag
		if prefix != "" {
			key = prefix + "." + jsonTag
		}

		fieldValue := v.Field(i)
		if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
			continue
		}

		if fieldValue.Kind() == reflect.Struct {
			setConfigRecursive(fieldValue, key, changedCount)
		} else {
			value := fieldValue.Interface()
			if !isEmptyValue(fieldValue) {
				viper.Set(key, value)
				(*changedCount)++
			}
		}
	}
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
