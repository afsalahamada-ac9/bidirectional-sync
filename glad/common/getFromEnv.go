package commontools

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

func GetFromEnv(key string) string {
	viper.SetConfigFile("/Users/afsalahamada/Desktop/SF to AWS/glad/common/.env")
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

func UpdateEnv(key, value string) error {
	log.Println("Updating the .env file")

	// Set the config file path
	viper.SetConfigFile("/Users/afsalahamada/Desktop/SF to AWS/glad/common/.env")

	// Read the existing config
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("Error reading the config file:", err)
		return err
	}

	// Update the value in Viper
	viper.Set(key, value)

	// Open the .env file for writing
	file, err := os.OpenFile(viper.ConfigFileUsed(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Println("Error opening the .env file:", err)
		return err
	}
	defer file.Close()

	// Write the updated key-value pairs to the file
	for _, k := range viper.AllKeys() {
		line := k + "=" + viper.GetString(k) + "\n"
		if _, err := file.WriteString(line); err != nil {
			log.Println("Error writing to the .env file:", err)
			return err
		}
	}

	log.Println("Successfully updated the .env file")
	return nil
}
