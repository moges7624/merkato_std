package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

type Config struct {
	Env     string
	Port    string
	DSN     string
	TestDSN string
}

func LoadConfig() (*Config, error) {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)

	root := filepath.Join(basePath, "../..")
	targetPath := filepath.Join(root, ".env")

	if err := godotenv.Load(targetPath); err != nil {
		return nil, errors.New("no .env file found")
	}

	cfg := &Config{
		Env:     GetEnv("ENV", "development"),
		Port:    GetEnv("PORT", "4000"),
		DSN:     GetEnv("MERKATO_STD_DB_DSN", ""),
		TestDSN: GetEnv("TEST_MERKATO_STD_DB_DSN", ""),
	}

	if cfg.DSN == "" {
		return nil, fmt.Errorf("LoadConfig: DSN is required")
	}

	return cfg, nil
}

func GetEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}

	return defaultVal
}
