package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	defaultExt       = ".state"
	defaultTrysCount = 10
	defaultChunkSize = 1000
	minChunkSize     = 100
	maxChunkSize     = 10000
	importFn         = "importChunkKV"
	exportFn         = "exportChunkKV"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "migration-manager",
	Short: "HLF migration-manager tool",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "migration.yaml", "HLF connection configuration file")
}
