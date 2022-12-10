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

func parseEnv(key string, def string) (*string, bool) {
	env, exists := os.LookupEnv(key)

	if !exists {
		return &def, false
	}

	return &env, true
}

func mustParseEnv(key string, def string) *string {

	env, exists := os.LookupEnv(key)

	if !exists {
		return &def
	}

	return &env
}

func parseEnvToArr(key string, def []string) *[]string {

	arrStr, exists := parseEnv(key, "")

	if !exists {
		return &def
	}

	envs := strings.Split(*arrStr, ",")

	for i, env := range envs {

		envs[i] = strings.TrimSpace(env)

	}

	return &envs
}

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

func ParseConfig() *Config {
	conf := Config{}

	conf.Mode = mustParseEnvToBool("MODE", 0) // 0 stands for false
	conf.GRPC_PORT = *mustParseEnv("GRPC_PORT", ":6000")
	conf.WS_PORT = *mustParseEnv("WS_PORT", ":6001")
	conf.WS_ReadBufferSize = *mustParseEnvToInt("WS_READ_BUFFER", 1024)
	conf.WS_WriteBufferSize = *mustParseEnvToInt("WS_WRITE_BUFFER", 1024)
	conf.AllowedOrigins = *parseEnvToArr("ALLOWED_ORIGINS", []string{"*"})

	return &conf
}
