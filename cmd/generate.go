package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/shakahl/graphql-typedef-go/internal/config"
	"github.com/shakahl/graphql-typedef-go/internal/generator"
)

// generateCmd represents the "generate" command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Go package for GraphQL schema types",
	Long:  `Generate Go package for GraphQL schema types based on GraphQL server schema introspection query`,
	Run: func(cmd *cobra.Command, args []string) {

		logger := log.New(os.Stdout, "graphql_type_def_generator", log.LstdFlags|log.LUTC)
		cfg := config.Get()

		gopts := generator.GraphQLTypeDefGeneratorOptions{
			Endpoint:        cfg.GraphQLEndpoint,
			AuthHeader:      cfg.GraphQLAuthHeader,
			AuthToken:       cfg.GraphQLAuthToken,
			OutputDirectory: cfg.OutputDirectory,
			OutputPackage:   cfg.OutputPackage,
		}

		gen := generator.New(gopts, logger)

		if err := gen.Generate(); err != nil {
			logger.Fatalln(err)
		}
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
