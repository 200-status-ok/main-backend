package utils

import (
	"github.com/spf13/viper"
)

func ReadFromEnvFile(path string, key string) string {
	viper.SetConfigFile(path)
	_ = viper.ReadInConfig()
	if viper.Get(key) == nil {
		return ""
	}
	return viper.Get(key).(string)
}
