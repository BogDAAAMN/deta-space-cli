package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/deta/pc-cli/internal/api"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "space",
		Short: "Deta Space CLI for mananging Deta Space projects",
		Long: fmt.Sprintf(`Deta Space command line interface for managing Deta Space projects. 
Complete documentation available at %s`, docsUrl),
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
		// no usage shown on errors
		SilenceUsage: false,
	}

	client = api.NewDetaClient()

	logger = log.New(os.Stderr, "", 0)
)

// Execute xx
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
