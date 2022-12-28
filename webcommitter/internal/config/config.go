package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Build information -ldflags .
const (
	version    string = "dev"
	commitHash string = "-"
)

var cfg *Config

// GetConfigInstance returns service config
func GetConfigInstance() Config {
	if cfg != nil {
		return *cfg
	}

	return Config{}
}

// Database - contains all parameters database connection.
type Database struct {
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	Migrations string `yaml:"migrations"`
	Name       string `yaml:"name"`
	SslMode    string `yaml:"sslmode"`
	Driver     string `yaml:"driver"`
}

// Project - contains all parameters project information.
type Project struct {
	Debug       bool   `yaml:"debug"`
	Name        string `yaml:"name"`
	Environment string `yaml:"environment"`
	Version     string
	CommitHash  string
}

// Status config for service.
type Status struct {
	Port          int    `yaml:"port"`
	Host          string `yaml:"host"`
	VersionPath   string `yaml:"versionPath"`
	LivenessPath  string `yaml:"livenessPath"`
	ReadinessPath string `yaml:"readinessPath"`
}

// Config - contains all configuration parameters in config package.
type Config struct {
	Project   Project   `yaml:"project"`
	Database  Database  `yaml:"database"`
	Status    Status    `yaml:"status"`
	Telegramm Telegramm `yaml:"telegramm"`
	Gitlab    Gitlab    `yaml:"gitlab"`
	Rest      Rest      `yaml:"rest"`
}

type Rest struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Pathfile string `yaml:"pathfile"`
}

type Telegramm struct {
	Token         string
	Chat_id_error uint64
}

type Gitlab struct {
	Url         string `yaml:"url"`
	Project_url string `yaml:"project_url"`
	CurrPath    string `yaml:"path_temp_project"`
}

// ReadConfigYML - read configurations from file and init instance Config.
func ReadConfigYML(filePath string) error {
	if cfg != nil {
		return nil
	}

	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return err
	}

	cfg.Project.Version = version
	cfg.Project.CommitHash = commitHash

	return nil
}
