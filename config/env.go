package config

import (
	"checkin/logging"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var EnvConfig = envConfigSchema{}

func init() {
	envInit()
}

var defaultConfig = envConfigSchema{
	ENV: "dev",

	MAX_REQUEST_BODY_SIZE: 200 * 1024 * 1024,
}

type envConfigSchema struct {
	ENV string `env:"ENV,DREAM_ENV"`

	MAX_REQUEST_BODY_SIZE int
}

func (s *envConfigSchema) IsDev() bool {
	return s.ENV == "dev" || s.ENV == "TESTING"
}

// envInit Reads .env as environment variables and fill corresponding fields into EnvConfig.
// To use a value from EnvConfig , simply call EnvConfig.FIELD like EnvConfig.ENV
// Note: Please keep Env as the first field of envConfigSchema for better logging.
func envInit() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file, ignored")
	}
	v := reflect.ValueOf(defaultConfig)
	typeOfV := v.Type()

	for i := 0; i < v.NumField(); i++ {
		envNameAlt := make([]string, 0)
		fieldName := typeOfV.Field(i).Name
		fieldType := typeOfV.Field(i).Type
		fieldValue := v.Field(i).Interface()

		envNameAlt = append(envNameAlt, fieldName)
		if fieldTag, ok := typeOfV.Field(i).Tag.Lookup("env"); ok && len(fieldTag) > 0 {
			tags := strings.Split(fieldTag, ",")
			envNameAlt = append(envNameAlt, tags...)
		}

		switch fieldType {
		case reflect.TypeOf(0):
			{
				configDefaultValue, ok := fieldValue.(int)
				if !ok {
					logging.Logger.WithFields(logrus.Fields{
						"field": fieldName,
						"type":  "int",
						"value": fieldValue,
						"env":   envNameAlt,
					}).Warningf("Failed to parse default value")
					continue
				}
				envValue := resolveEnv(envNameAlt, fmt.Sprintf("%d", configDefaultValue))
				if EnvConfig.IsDev() {
					fmt.Printf("Reading field[ %s ] default: %v env: %s\n", fieldName, configDefaultValue, envValue)
				}
				if len(envValue) > 0 {
					envValueInteger, err := strconv.ParseInt(envValue, 10, 64)
					if err != nil {
						logging.Logger.WithFields(logrus.Fields{
							"field": fieldName,
							"type":  "int",
							"value": fieldValue,
							"env":   envNameAlt,
						}).Warningf("Failed to parse env value, ignored")
						continue
					}
					reflect.ValueOf(&EnvConfig).Elem().Field(i).SetInt(envValueInteger)
				}
				continue
			}
		case reflect.TypeOf(""):
			{
				configDefaultValue, ok := fieldValue.(string)
				if !ok {
					logging.Logger.WithFields(logrus.Fields{
						"field": fieldName,
						"type":  "int",
						"value": fieldValue,
						"env":   envNameAlt,
					}).Warningf("Failed to parse default value")
					continue
				}
				envValue := resolveEnv(envNameAlt, configDefaultValue)

				if EnvConfig.IsDev() {
					fmt.Printf("Reading field[ %s ] default: %v env: %s\n", fieldName, configDefaultValue, envValue)
				}
				if len(envValue) > 0 {
					reflect.ValueOf(&EnvConfig).Elem().Field(i).SetString(envValue)
				}
			}
		}

	}
}

func resolveEnv(configKeys []string, defaultValue string) string {
	for _, item := range configKeys {
		envValue := os.Getenv(item)
		if envValue != "" {
			return envValue
		}
	}
	return defaultValue
}
