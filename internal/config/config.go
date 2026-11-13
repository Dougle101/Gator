package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const cfgFP = ".gatorconfig.json"

func get_fp() (fp string, e error) {
	fp, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Could not find home directory: %v", err)
	}
	return filepath.Join(fp, cfgFP), nil
}

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func write(c Config) error {
	data, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("Error encoding config into JSON: %v", err)
	}

	fp, err := get_fp()
	if err != nil {
		return err
	}

	err = os.WriteFile(fp, data, 0o600)
	if err != nil {
		return fmt.Errorf("Error writing file: %v\nThis may have caused a partial write, please check file integrity.", err)
	}
	return nil
}

func (c Config) SetUser(user string) error {
	c.CurrentUserName = user
	err := write(c)
	if err != nil {
		return err
	}
	return nil
}

func Read() (Config, error) {
	fp, err := get_fp()
	if err != nil {
		return Config{}, err
	}

	conf, err := os.ReadFile(fp)
	if err != nil {
		return Config{}, fmt.Errorf("Error reading config file: %v", err)
	}

	var dec Config
	err = json.Unmarshal(conf, &dec)
	if err != nil {
		return Config{}, fmt.Errorf("Error decoding config file into struct: %v", err)
	}
	return dec, nil
}
