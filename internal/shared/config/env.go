package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

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
	AUTH_SECRET          ENV_VARIABLE = "AUTH_SECRET"
	GO_ENV               ENV_VARIABLE = "GO_ENV"
	PORT                 ENV_VARIABLE = "PORT"
	MONGODB_URI          ENV_VARIABLE = "MONGODB_URI"
	MONGODB_NAME         ENV_VARIABLE = "MONGODB_NAME"
	MONGODB_REPLICA_SET  ENV_VARIABLE = "MONGODB_REPLICA_SET"
	TIMEOUT              ENV_VARIABLE = "TIMEOUT"
	MAX_POOL_SIZE        ENV_VARIABLE = "MAX_POOL_SIZE"
	MIN_POOL_SIZE        ENV_VARIABLE = "MIN_POOL_SIZE"
	CONN_IDLE_TIME       ENV_VARIABLE = "CONN_IDLE_TIME"
	MONGODB_RETRY_WRITES ENV_VARIABLE = "MONGODB_RETRY_WRITES"
	MONGODB_RETRY_READS  ENV_VARIABLE = "MONGODB_RETRY_READS"
	POSTGRES_URI         ENV_VARIABLE = "POSTGRES_URI"
)

type env struct {
	authSecret         string
	goEnv              Environment
	port               string
	mongoDbURI         string
	mongoDbName        string
	mongoDbReplicaSet  string
	maxPoolSize        int32
	minPoolSize        int32
	connIdleTime       time.Duration
	mongoDbRetryWrites bool
	mongoDbRetryReads  bool
	timeout            time.Duration
	postgresURI        string
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

func NewEnv() *env {
	return &env{
		authSecret:         mustGetEnv(AUTH_SECRET),
		goEnv:              Environment(mustGetEnv(GO_ENV)),
		port:               getEnvWithDefault(PORT, DEFAULT_PORT),
		mongoDbURI:         getEnv(MONGODB_URI),
		mongoDbName:        getEnv(MONGODB_NAME),
		mongoDbReplicaSet:  getEnvWithDefault(MONGODB_REPLICA_SET, "rs0"),
		timeout:            time.Duration(getEnvInt(TIMEOUT, 10)) * time.Second,
		maxPoolSize:        int32(getEnvInt(MAX_POOL_SIZE, 100)),
		minPoolSize:        int32(getEnvInt(MIN_POOL_SIZE, 5)),
		connIdleTime:       time.Duration(getEnvInt(CONN_IDLE_TIME, 30)) * time.Minute,
		mongoDbRetryWrites: getEnvBool(MONGODB_RETRY_WRITES, true),
		mongoDbRetryReads:  getEnvBool(MONGODB_RETRY_READS, true),
		postgresURI:        getEnv(POSTGRES_URI),
	}
}

func (e *env) AuthSecret() string {
	return e.authSecret
}
func (e *env) Port() string {
	return e.port
}

func (e *env) PostgresURI() string {
	return e.postgresURI
}

func (e *env) MongoDbURI() string {
	return e.mongoDbURI
}

func (e *env) MongoDbReplicaSet() string {
	return e.mongoDbReplicaSet
}

func (e *env) MongoDbName() string {
	return e.mongoDbName
}
func (e *env) Timeout() time.Duration {
	return e.timeout
}

func (e *env) GoEnv() Environment {
	return e.goEnv
}

func (e *env) MaxPoolSize() int32 {
	return e.maxPoolSize
}

func (e *env) MinPoolSize() int32 {
	return e.minPoolSize
}

func (e *env) ConnIdleTime() time.Duration {
	return e.connIdleTime
}

func (e *env) MongoDbRetryReads() bool {
	return e.mongoDbRetryReads
}

func (e *env) MongoDbRetryWrites() bool {
	return e.mongoDbRetryWrites
}

func getEnv(key ENV_VARIABLE) string {
	return os.Getenv(strings.TrimSpace(string(key)))
}

func getEnvWithDefault(key ENV_VARIABLE, defaultValue string) string {
	value := getEnv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func mustGetEnv(key ENV_VARIABLE) string {
	value := getEnv(key)
	if value == "" {
		panic(fmt.Sprintf("%s environment variable is not set", key))
	}
	return value
}

func getEnvInt(key ENV_VARIABLE, defaultValue int) int {
	value := getEnv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Sprintf("%s environment variable is not a valid integer", key))
	}
	return intValue
}

func getEnvBool(key ENV_VARIABLE, defaultValue bool) bool {
	value := getEnv(key)
	if value == "" {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		panic(fmt.Sprintf("%s environment variable is not a valid boolean", key))
	}
	return boolValue
}
