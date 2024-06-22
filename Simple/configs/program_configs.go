package configs

import (
	"encoding/json"
)

///////////////////////////////////////////////////////

type ProgConfig struct{}

// Fix implements Configer interface.
//
// This never returns an error.
func (po *ProgConfig) Fix() error {
	return nil
}

// Default implements Configer interface.
func (po *ProgConfig) Default() Configer {
	config := &ProgConfig{}
	return config
}

// MarshalJSON implements Configer interface.
func (po *ProgConfig) MarshalJSON() ([]byte, error) {
	type Alias struct{}

	a := &Alias{}

	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return nil, err
	}

	return data, nil
}

// UnmarshalJSON implements Configer interface.
func (po *ProgConfig) UnmarshalJSON(data []byte) error {
	type Alias struct{}

	err := json.Unmarshal(data, &Alias{})
	if err != nil {
		return err
	}

	return nil
}
