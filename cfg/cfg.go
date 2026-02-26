package cfg

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// EnvPrefix environment prefix.
const EnvPrefix = "MIGRATION"

// HLF contains hlf connection settings.
type HLF struct {
	Config      string        `mapstructure:"config"`
	Org         string        `mapstructure:"org"`
	User        string        `mapstructure:"user"`
	Channel     string        `mapstructure:"channel"`
	Chaincode   string        `mapstructure:"chaincode"`
	UseSmartBFT bool          `mapstructure:"usebft"`
	ExecTimeout time.Duration `mapstructure:"exectimeout"`
}

// Config contains listen and HLF paramethers.
type Config struct {
	SnapshotDir string `mapstructure:"snapshot"`
	HLF         *HLF   `mapstructure:"hlf"`
}

// ReadFromFile reads config from disk by using config viper.
func ReadFromFile(envPrefix, filename string, out any) error {
	viper.AutomaticEnv()
	viper.SetTypeByDefaultValue(true)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix(envPrefix)
	viper.SetConfigFile(filename)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed viper.ReadInConfig: %w", err)
	}

	if err := viper.Unmarshal(out); err != nil {
		return fmt.Errorf("failed viper.Unmarshal: %w", err)
	}

	return nil
}
