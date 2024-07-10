package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

// Config represents the structure of the configuration file.
type Config struct {
	Env      string `yaml:"env" env-default:"prod"`
	Database Database
	Server   Server
}

// Database configuration details.
type Database struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"5432"`
	User     string `yaml:"userName" env-default:"postgres"`
	Password string `yaml:"userPassword" env-default:"postgres"`
	DBName   string `yaml:"dbname" env-required:"true"`
	SSLMode  string `yaml:"sslMode" env-default:"disable"`
}

// Server configuration details.
// TODO: It may be necessary to supplement the parameters.
type Server struct {
	Port string `yaml:"port" env-default:"8080"`
}

// MustLoadConfig loads the configuration from the specified path.
func MustLoadConfig() *Config {
	configPath := searchPathConfig()
	if configPath == "" {
		panic("empty config path")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file not found: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

// searchPathConfig retrieves the configuration path from either
// a command-line flag or an environment variable.
func searchPathConfig() string {
	var way string
	flag.StringVar(&way, "config", "", "config path")
	flag.Parse()

	if way == "" {
		way = os.Getenv("CFG_authREST")
	}

	return way
}
