package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/shakahl/gqlassist/internal/config"
)

// initCmd represents the "generate" command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration file",
	Long:  `Initialize configuration file`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("INIT: %+v", config.Get())
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
