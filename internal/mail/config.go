package mail

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	FromName string
	UseTLS   bool
}

func LoadConfigFromEnv() (Config, error) {
	cfg := Config{
		Host:     os.Getenv("SMTP_HOST"),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     os.Getenv("SMTP_FROM"),
		FromName: os.Getenv("SMTP_FROM_NAME"),
	}

	portRaw := os.Getenv("SMTP_PORT")
	if portRaw == "" {
		return Config{}, fmt.Errorf("email: load config: SMTP_PORT is required")
	}

	port, err := strconv.Atoi(portRaw)
	if err != nil {
		return Config{}, fmt.Errorf("email: load config: invalid SMTP_PORT %q %w", portRaw, err)
	}

	cfg.Port = port

	tlsRaw := os.Getenv("SMTP_TLS")
	if tlsRaw == "" {
		return Config{}, fmt.Errorf("email: load config: SMTP_TLS is required")
	}

	useTLS, err := strconv.ParseBool(tlsRaw)
	if err != nil {
		return Config{}, fmt.Errorf("email: load config: invalid SMTP_TLS %q %w", tlsRaw, err)
	}

	cfg.UseTLS = useTLS

	if err := cfg.validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) validate() error {
	var missing []string

	if c.Host == "" {
		missing = append(missing, "SMTP_HOST")
	}

	if c.Port <= 0 || c.Port > 65535 {
		missing = append(missing, "SMTP_PORT")
	}

	if c.Password == "" {
		missing = append(missing, "SMTP_PASSWORD")
	}

	if c.From == "" {
		missing = append(missing, "SMTP_FROM")
	}

	if len(missing) > 0 {
		return fmt.Errorf("email: validate config: missing or invalid %v", missing)
	}

	return nil
}
