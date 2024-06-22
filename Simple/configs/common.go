package configs

import (
	"encoding/json"
	"os"

	com "github.com/PlayerR9/LyneCmL/Simple/common"
)

///////////////////////////////////////////////////////

const (
	// DisplayConfig is the display configuration.
	DisplayConfig string = "display"

	// InternalsConfig is the internals configuration.
	InternalsConfig string = "internals"
)

// Configer is a configuration.
type Configer interface {
	// Default gets the default configuration.
	//
	// Returns:
	//   - Configer: The default configuration.
	Default() Configer

	json.Marshaler
	json.Unmarshaler

	com.Fixer
}

// Config is a configuration.
type Config struct {
	// loc is the location of the configuration file.
	loc string

	// filePerm is the file permissions of the configuration file.
	filePerm os.FileMode

	// configTable is the table of configurations.
	configTable map[string]Configer
}

// Fix implements Configer interface.
//
// This never errors.
func (c *Config) Fix() error {
	conf, ok := c.configTable[DisplayConfig]
	if !ok {
		def := (*new(DisplayConfigs)).Default()
		c.configTable[DisplayConfig] = def
	} else {
		conf.Fix()
	}

	conf, ok = c.configTable[InternalsConfig]
	if !ok {
		def := (*new(ProgConfig)).Default()
		c.configTable[InternalsConfig] = def
	} else {
		conf.Fix()
	}

	for key, config := range c.configTable {
		if config == nil {
			delete(c.configTable, key)
		} else {
			config.Fix()
		}
	}

	return nil
}

// MarshalJSON implements Configer interface.
func (c *Config) MarshalJSON() ([]byte, error) {
	data, err := json.MarshalIndent(c.configTable, "", "  ")
	if err != nil {
		return nil, err
	}

	return data, nil
}

// UnmarshalJSON implements Configer interface.
func (c *Config) UnmarshalJSON(data []byte) error {
	configMap := make(map[string]Configer)

	err := json.Unmarshal(data, &configMap)
	if err != nil {
		return err
	}

	c.configTable = configMap

	return nil
}

// fileperm = 0644
func NewConfig(loc string, filePerm os.FileMode) *Config {
	config := &Config{
		loc:         loc,
		filePerm:    filePerm,
		configTable: make(map[string]Configer),
	}

	config.configTable[DisplayConfig] = (*new(DisplayConfigs)).Default()
	config.configTable[InternalsConfig] = (*new(ProgConfig)).Default()

	return config
}

// AddConfig adds a configuration to the configuration table.
//
// Parameters:
//   - key: The key of the configuration.
//   - config: The configuration to add.
//
// Returns:
//   - bool: True if the configuration was added, false if it already exists.
func (c *Config) AddConfig(key string, config Configer) bool {
	if config == nil {
		return true
	}

	_, ok := c.configTable[key]
	if ok {
		return false
	}

	c.configTable[key] = config

	return true
}

// Load loads the configuration from the file.
//
// Returns:
//   - error: An error if the configuration failed to load.
func (c *Config) Load() error {
	data, err := os.ReadFile(c.loc)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, c)
	if err != nil {
		return err
	}

	c.Fix()

	return nil
}

// Save saves the configuration to the file.
//
// Returns:
//   - error: An error if the configuration failed to save.
func (c *Config) Save() error {
	c.Fix()

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(c.loc, data, c.filePerm)
	if err != nil {
		return err
	}

	return nil
}

// GetConfigs gets the configuration.
//
// Parameters:
//   - key: The key of the configuration.
//
// Returns:
//   - Configer: The configuration. Nil if it does not exist.
func (c *Config) GetConfigs(key string) Configer {
	config, ok := c.configTable[key]
	if !ok {
		return nil
	}

	return config
}
