package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Env                string
	Port               int
	DSN                string
	TestDSN            string
	AccessTokenSecret  string
	RefreshTokenSecret string
}

func LoadConfig() (*Config, error) {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)

	root := filepath.Join(basePath, "../..")
	targetPath := filepath.Join(root, ".env")

	if err := godotenv.Load(targetPath); err != nil {
		return nil, errors.New("no .env file found")
	}

	port, err := GetEnvInt("PORT", 4000)
	if err != nil {
		return nil, fmt.Errorf("LoadConfig: invalid port number")
	}

	cfg := &Config{
		Env:                GetEnv("ENV", "development"),
		Port:               port,
		DSN:                GetEnv("MERKATO_STD_DB_DSN", ""),
		TestDSN:            GetEnv("TEST_MERKATO_STD_DB_DSN", ""),
		AccessTokenSecret:  GetEnv("ACCESS_TOKEN_SECRET", ""),
		RefreshTokenSecret: GetEnv("REFRESH_TOKEN_SECRET", ""),
	}

	if cfg.DSN == "" {
		return nil, fmt.Errorf("LoadConfig: env var DSN is required")
	}

	if cfg.AccessTokenSecret == "" {
		return nil, fmt.Errorf("LoadConfig: env var ACCESS_TOKEN_SECRET is required")
	}

	return cfg, nil
}

func GetEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}

	return defaultVal
}

func GetEnvInt(key string, defaultVal int) (int, error) {
	if val, exists := os.LookupEnv(key); exists {
		res, err := strconv.Atoi(val)
		if err != nil {
			return 0, fmt.Errorf("invalid value for env var %s", key)
		}

		return res, nil
	}

	return defaultVal, nil
}
