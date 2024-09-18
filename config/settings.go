package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config holds the database and Redis connection pool configurations
type Config struct {
	Type        string      // Database type (e.g., "postgres", "mysql", "redis")
	DBConfig    DBConfig    // Database configuration
	RedisConfig RedisConfig // Redis configuration
}

// DBConfig holds the database connection pool configuration
type DBConfig struct {
	DSN             string        // Data Source Name
	MaxIdleConns    int           // Maximum number of idle connections
	MaxOpenConns    int           // Maximum number of open connections
	ConnMaxIdleTime time.Duration // Maximum idle time for a connection
	ConnMaxLifetime time.Duration // Maximum lifetime of a connection
}

// RedisConfig holds the Redis connection configuration
type RedisConfig struct {
	Addr     string // Redis server address
	Password string // Redis password (leave empty if no password)
	DB       int    // Redis database index
}

// initConfig initializes Viper to read configuration from environment variables
func InitConfig() *Config {
	// Set Viper to read environment variables
	viper.AutomaticEnv() // Enable automatic environment variable reading

	viper.SetDefault("DB_TYPE", "redis") // Default database type

	// Set default values for database configuration
	viper.SetDefault("DB_DSN", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	viper.SetDefault("DB_MAX_IDLE_CONNS", 10)
	viper.SetDefault("DB_MAX_OPEN_CONNS", 20)
	viper.SetDefault("DB_CONN_MAX_IDLE_TIME", "10m")
	viper.SetDefault("DB_CONN_MAX_LIFETIME", "1h")

	// Set default values for Redis configuration
	viper.SetDefault("REDIS_ADDR", "localhost:6379")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", 0)

	// Create and return Config with values from environment variables or defaults
	return &Config{
		Type: viper.GetString("DB_TYPE"),
		DBConfig: DBConfig{
			DSN:             viper.GetString("DB_DSN"),
			MaxIdleConns:    viper.GetInt("DB_MAX_IDLE_CONNS"),
			MaxOpenConns:    viper.GetInt("DB_MAX_OPEN_CONNS"),
			ConnMaxIdleTime: viper.GetDuration("DB_CONN_MAX_IDLE_TIME"),
			ConnMaxLifetime: viper.GetDuration("DB_CONN_MAX_LIFETIME"),
		},
		RedisConfig: RedisConfig{
			Addr:     viper.GetString("REDIS_ADDR"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
	}
}
