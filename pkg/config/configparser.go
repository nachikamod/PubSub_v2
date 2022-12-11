package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Mode               bool
	GRPC_PORT          string
	WS_PORT            string
	AllowedOrigins     []string
	WS_ReadBufferSize  int
	WS_WriteBufferSize int
}

// Parse envrionment variable in different data types

// parse environment variable, if exists return the value and true else return default value and false. Here, true or false denotes if the env exists or not
func parseEnv(key string, def string) (*string, bool) {
	env, exists := os.LookupEnv(key)

	if !exists {
		return &def, false
	}

	return &env, true
}

// Must parse env method parse if env exists else returns default
func parseEnvOrDefault(key string, def string) *string {

	env, exists := os.LookupEnv(key)

	if !exists {
		return &def
	}

	return &env
}

// Parse env to array, checks if env exists then parse array from comma separated string else return default array
func parseEnvToArr(key string, def []string) *[]string {

	arrStr, exists := parseEnv(key, "")

	if !exists {
		return &def
	}

	envs := strings.Split(*arrStr, ",")

	for i, env := range envs {

		envs[i] = strings.TrimSpace(env) // Remove any leading and trailing white spaces

	}

	return &envs
}

// Must parse env to int method parse if env exists then convert to integer else returns default
func parseEnvToInt(key string, def int) (*int, error) {

	intStr, exists := parseEnv(key, "")

	if !exists {
		return &def, nil
	}

	var env int = 0

	env, err := strconv.Atoi(*intStr)

	if err != nil {
		return &env, err
	}

	return &env, nil
}

// Must parse env to bool method parse if env exists then convert to boolean else returns default
func parseEnvToBool(key string, def int) (bool, error) {

	env, err := parseEnvToInt(key, def)

	if err != nil {
		return false, err
	}

	if *env == 1 {
		return true, nil
	} else if *env == 0 {
		return false, nil
	}

	return false, errors.New(fmt.Sprintf("Error converting env to bool : parsing \"%v\" : invalid env %v", key, env))
}

// Parse environment variable
func ParseConfig() (*Config, error) {
	conf := Config{}
	// Each methods accepts two variables - key, default value
	mode, err := parseEnvToBool("MODE", 1) // 1 stands for development mode

	if err != nil {
		return &conf, err
	}

	conf.Mode = mode

	conf.GRPC_PORT = *parseEnvOrDefault("GRPC_PORT", ":6000")
	conf.WS_PORT = *parseEnvOrDefault("WS_PORT", ":6001")
	rbs, err := parseEnvToInt("WS_READ_BUFFER", 1024)

	if err != nil {
		return &conf, err
	}

	conf.WS_ReadBufferSize = *rbs

	wbs, err := parseEnvToInt("WS_WRITE_BUFFER", 1024)

	if err != nil {
		return &conf, err
	}

	conf.WS_WriteBufferSize = *wbs

	conf.AllowedOrigins = *parseEnvToArr("ALLOWED_ORIGINS", []string{"*"})

	return &conf, err
}
