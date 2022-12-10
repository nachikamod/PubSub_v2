package utils

import (
	"log"
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
func mustParseEnv(key string, def string) *string {

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
func mustParseEnvToInt(key string, def int) *int {

	intStr, exists := parseEnv(key, "")

	if !exists {
		return &def
	}

	env, err := strconv.Atoi(*intStr)

	if err != nil {
		log.Fatalf("Error converting env to int : %v", err)
	}

	return &env
}

// Must parse env to bool method parse if env exists then convert to boolean else returns default
func mustParseEnvToBool(key string, def int) bool {
	env := *mustParseEnvToInt(key, def)

	if env == 1 {
		return true
	} else if env == 0 {
		return false
	} else {
		log.Fatalf("Error converting env to bool : parsing \"%v\" : invalid env %v", key, env)
	}
	return false
}

// Parse environment variable
func ParseConfig() *Config {
	conf := Config{}

	// Each methods accepts two variables - key, default value
	conf.Mode = mustParseEnvToBool("MODE", 1) // 1 stands for development mode
	conf.GRPC_PORT = *mustParseEnv("GRPC_PORT", ":6000")
	conf.WS_PORT = *mustParseEnv("WS_PORT", ":6001")
	conf.WS_ReadBufferSize = *mustParseEnvToInt("WS_READ_BUFFER", 1024)
	conf.WS_WriteBufferSize = *mustParseEnvToInt("WS_WRITE_BUFFER", 1024)
	conf.AllowedOrigins = *parseEnvToArr("ALLOWED_ORIGINS", []string{"*"})

	return &conf
}
