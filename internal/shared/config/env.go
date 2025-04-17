package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/joho/godotenv"
)

type Environment string

const (
	Development Environment = "development"
	Production  Environment = "production"
	Test        Environment = "test"
)

type ENV_VARIABLE string

const (
	AUTH_SECRET ENV_VARIABLE = "AUTH_SECRET"
	GO_ENV      ENV_VARIABLE = "GO_ENV"
)

type env struct {
	authSecret string
	goEnv      Environment
}

func (e *env) AuthSecret() string {
	return e.authSecret
}

func (e *env) GoEnv() Environment {
	return e.goEnv
}

func init() {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(b), "../../../")
	envPath := filepath.Join(basePath, ".env")

	if err := godotenv.Load(envPath); err != nil {
		log.Println(".env file not found or failed to load")
	}
}

func Env() *env {
	authSecret := getEnv(AUTH_SECRET)
	if strings.TrimSpace(authSecret) == "" {
		panic(fmt.Sprintf("%s environment variable is not set", AUTH_SECRET))
	}
	goEnv := getEnv(GO_ENV)
	if strings.TrimSpace(goEnv) == "" {
		panic(fmt.Sprintf("%s environment variable is not set", GO_ENV))
	}

	return &env{
		authSecret: authSecret,
		goEnv:      Environment(goEnv),
	}
}

func getEnv(key ENV_VARIABLE) string {
	return os.Getenv(string(key))
}
