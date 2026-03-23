package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(ConfigDir())

	viper.SetEnvPrefix("PORTAINER")
	viper.AutomaticEnv()

	viper.SetDefault("url", "")
	viper.SetDefault("endpoint_id", 2)

	_ = viper.ReadInConfig()
}

func ConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "portainer-cli")
}

func URL() string {
	return viper.GetString("url")
}

func EndpointID() int {
	return viper.GetInt("endpoint_id")
}

func Save() error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	return viper.WriteConfigAs(filepath.Join(dir, "config.yaml"))
}

// Token represents a stored JWT token.
type Token struct {
	JWT       string    `json:"jwt"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (t *Token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt.Add(-60 * time.Second))
}

func tokenPath() string {
	return filepath.Join(ConfigDir(), "token.json")
}

func LoadToken() (*Token, error) {
	data, err := os.ReadFile(tokenPath())
	if err != nil {
		return nil, fmt.Errorf("not authenticated, run 'portainer login'")
	}
	var t Token
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, fmt.Errorf("corrupt token file: %w", err)
	}
	return &t, nil
}

func SaveToken(t *Token) error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(tokenPath(), data, 0600)
}
