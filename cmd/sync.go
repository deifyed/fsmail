package cmd

import (
	"fmt"
	"os"

	"github.com/deifyed/fssmtp/cmd/sync"
	"github.com/spf13/cobra"
)

var targetDir string

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "synchronizes directory with server",
	RunE:  sync.RunE(fs, &targetDir),
}

func init() {
	rootCmd.AddCommand(syncCmd)

	workDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("getting work directory: %w", err))
	}

	syncCmd.Flags().StringVarP(&targetDir, "directory", "d", workDir, "target directory")
}
