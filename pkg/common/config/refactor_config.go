package config

import (
	"fmt"
	"github.com/spf13/viper"
	"reflect"
	"strings"
)

func Init() {
	viper.SetConfigFile("./config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	err := viper.Unmarshal(&Config)
	if err != nil {
		panic(err)
	}

	viper.SetConfigFile("./config/notification.yaml")

	if err = viper.ReadInConfig(); err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&Config.Notification)
	if err != nil {
		panic(err)
	}
	fmt.Println(Config)
}

func GetRedisAddress() []string {
	return viper.GetStringSlice("redis.address")
}

func GetRedisUsername() string {
	return viper.GetString("redis.username")
}

func GetRedisPassword() string {
	return viper.GetString("redis.password")
}

func customDecodeHook(from, to reflect.Type, data interface{}) (interface{}, error) {
	if to == reflect.TypeOf([]string{}) {
		// 解析为字符串
		str, ok := data.(string)
		if !ok {
			return nil, fmt.Errorf("invalid data type for []string")
		}
		// 切割字符串并去除空白项
		items := strings.Split(str, ",")
		var result []string
		for _, item := range items {
			item = strings.TrimSpace(item)
			if item != "" {
				result = append(result, item)
			}
		}
		return result, nil
	}
	// 使用默认解析器解析其他类型
	return data, nil
}
