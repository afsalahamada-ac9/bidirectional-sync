package Common

import (
	"log"

	"github.com/spf13/viper"
)

func GetFromEnv(key string) string {
	viper.SetConfigFile(".env")
	log.Println("key:", key)
	err := viper.ReadInConfig()
	if err != nil {
		// todo: clean architecture based implementation
		log.Println("there is an error in the config file", err)
	}
	value, ok := viper.Get(key).(string)
	log.Println(value, "received")
	if !ok {
		log.Println("error getting the key")
	}
	return value
}

func WriteToEnv(key, value string) error {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	viper.Set(key, value)
	err = viper.WriteConfigAs(".env")
	if err != nil {
		return err
	}
	return nil
}
