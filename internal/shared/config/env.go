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
	AUTH_SECRET         ENV_VARIABLE = "AUTH_SECRET"
	GO_ENV              ENV_VARIABLE = "GO_ENV"
	PORT                ENV_VARIABLE = "PORT"
	MONGODB_URI         ENV_VARIABLE = "MONGODB_URI"
	MONGODB_NAME        ENV_VARIABLE = "MONGODB_NAME"
	MONGODB_REPLICA_SET ENV_VARIABLE = "MONGODB_REPLICA_SET"
)

type env struct {
	authSecret        string
	goEnv             Environment
	port              string
	mongoDbURI        string
	mongoDbName       string
	mongoDbReplicaSet string
}

func (e *env) AuthSecret() string {
	return e.authSecret
}

func (e *env) GoEnv() Environment {
	return e.goEnv
}

func (e *env) MongoDbURI() string {
	return e.mongoDbURI
}
func (e *env) MongoDbName() string {
	return e.mongoDbName
}

func (e *env) Port() string {
	return e.port
}
func (e *env) MongoDbReplicaSet() string {
	return e.mongoDbReplicaSet
}

func init() {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(b), "../../../")
	envPath := filepath.Join(basePath, ".env")

	if err := godotenv.Load(envPath); err != nil {
		log.Println(".env file not found or failed to load")
	}
}

const DEFAULT_PORT = "8080"

func Env() *env {
	authSecret := getEnv(AUTH_SECRET)
	if strings.TrimSpace(authSecret) == "" {
		panic(fmt.Sprintf("%s environment variable is not set", AUTH_SECRET))
	}
	goEnv := getEnv(GO_ENV)
	if strings.TrimSpace(goEnv) == "" {
		panic(fmt.Sprintf("%s environment variable is not set", GO_ENV))
	}
	port := getEnv(PORT)
	mongoURI := getEnv(MONGODB_URI)
	if strings.TrimSpace(mongoURI) == "" {
		panic(fmt.Sprintf("%s environment variable is not set", MONGODB_URI))
	}
	mongoName := getEnv(MONGODB_NAME)
	if strings.TrimSpace(mongoURI) == "" {
		panic(fmt.Sprintf("%s environment variable is not set", MONGODB_NAME))
	}

	if strings.TrimSpace(port) == "" {
		port = DEFAULT_PORT
	}
	mongoReplicaSet := getEnv(MONGODB_REPLICA_SET)
	if strings.TrimSpace(mongoReplicaSet) == "" {
		panic(fmt.Sprintf("%s environment variable is not set", MONGODB_REPLICA_SET))
	}

	return &env{
		authSecret:        authSecret,
		goEnv:             Environment(goEnv),
		port:              port,
		mongoDbURI:        mongoURI,
		mongoDbName:       mongoName,
		mongoDbReplicaSet: mongoReplicaSet,
	}
}

func getEnv(key ENV_VARIABLE) string {
	return os.Getenv(string(key))
}
