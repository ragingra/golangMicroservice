package application

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	RedisAddress string
	ServerPort   uint16
}

func LoadConfig() Config {
	cfg := Config{
		RedisAddress: "localhost:6379",
		ServerPort:   8080,
	}

	if redisAddr := os.Getenv("REDIS_ADDR"); redisAddr != "" {
		cfg.RedisAddress = redisAddr
	}

	if serverPort := os.Getenv("SERVER_PORT"); serverPort != "" {
		port, err := strconv.ParseUint(serverPort, 10, 16)
		if err != nil {
			fmt.Println("failed to parse SERVER_PORT, using default")
		} else {
			cfg.ServerPort = uint16(port)
		}
	}

	return cfg
}
