package config

import (
	"log"

	"github.com/spf13/viper"
)

// AppConfig holds the application's configurations
type AppConfig struct {
	KafkaConfig    KafkaConfig
	RabbitMQConfig RabbitMQConfig
	TestConfig     TestConfig `mapstructure:"test_config"`
}

// KafkaConfig defines Kafka-specific configurations
type KafkaConfig struct {
	Brokers []string
	Topic   string
}

// RabbitMQConfig defines RabbitMQ-specific configurations
type RabbitMQConfig struct {
	URL       string
	QueueName string
}

// TestConfig defines parameters for the tests to be conducted
type TestConfig struct {
	NumberOfEvents int `mapstructure:"number_of_events"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(configPath string) (*AppConfig, error) {
	// Load configurations from file
	viper.SetConfigName("config")   // name of config file (without extension)
	viper.SetConfigType("yaml")     // or "json", "toml"
	viper.AddConfigPath(configPath) // path to look for the config file in

	viper.AutomaticEnv() // read in environment variables that match

	var config AppConfig
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
		return nil, err
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %s", err)
		return nil, err
	}

	return &config, nil
}
