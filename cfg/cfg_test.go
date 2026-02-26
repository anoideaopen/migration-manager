package cfg_test

import (
	"testing"

	"github.com/anoideaopen/migration-manager/cfg"
)

func TestConfigDefaults(t *testing.T) {
	t.Parallel()

	cfg1 := new(cfg.Config)
	_ = cfg.ReadFromFile(
		"MIGRATION",
		"test_config.yaml",
		cfg1,
	)

	t.Log(cfg1)
}
