package configs

// ProgConfig are the configurations for a program.
type ProgConfig struct {
	// ProgName is the name of the program.
	ProgName string `json:"prog_name"`
}

// Fix implements Configer interface.
//
// This never returns an error.
func (po *ProgConfig) Fix() error {
	return nil
}

// Default implements Configer interface.
func (po *ProgConfig) Default() {
	po.ProgName = ""
}
