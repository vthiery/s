package cmd

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"

	"github.com/zquestz/go-ucl"
	"github.com/zquestz/s/providers"
)

// Config stores all the application configuration.
type Config struct {
	Blacklist       []string                    `json:"blacklist"`
	Binary          string                      `json:"binary"`
	Cert            string                      `json:"cert"`
	CustomProviders []*providers.CustomProvider `json:"customProviders"`
	DisplayVersion  bool                        `json:"-"`
	Key             string                      `json:"key"`
	ListProviders   bool                        `json:"-"`
	Port            int                         `json:"port,string"`
	Provider        string                      `json:"provider"`
	ServerMode      bool                        `json:"-"`
	Verbose         bool                        `json:"verbose,string"`
	Whitelist       []string                    `json:"whitelist"`
}

// Load reads the configuration from ~/.s/config and loads it into the Config struct.
// The config is in UCL format.
func (c *Config) Load() error {
	conf, err := c.loadConfig()
	if err != nil {
		return err
	}

	// There are cases when we don't have a configuration.
	if conf != nil {
		err = c.applyConf(conf)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) loadConfig() ([]byte, error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(filepath.Join(u.HomeDir, ".s", "config"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}
	defer f.Close()

	ucl.Ucldebug = false
	data, err := ucl.NewParser(f).Ucl()
	if err != nil {
		return nil, err
	}

	conf, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func (c *Config) applyConf(conf []byte) error {
	err := json.Unmarshal(conf, c)
	if err != nil {
		return err
	}

	return nil
}