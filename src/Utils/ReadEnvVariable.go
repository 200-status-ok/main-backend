package Utils

import "github.com/spf13/viper"

func ReadFromEnvFile(path string, key string) (string, string) {
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()
	return viper.Get(key).(string), ""
}
