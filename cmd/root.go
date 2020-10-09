package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"

	"github.com/shakahl/graphql-typedef-go/internal/config"
	"github.com/shakahl/graphql-typedef-go/internal/utils"
	"github.com/shakahl/graphql-typedef-go/meta"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   meta.BinaryName,
	Short: meta.ShortDescription,
	Long:  meta.LongDescription,
	// 	Long: `A longer description that spans multiple lines and likely contains
	// examples and usage of using your application. For example:
	//
	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	// Run: func(cmd *cobra.Command, args []string) {},
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

	cobra.OnInitialize(onCobraInit)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is $HOME/%s)", meta.ConfigFileName))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.PersistentFlags().StringP("output-directory", "o", config.DefaultOutputDirectory, "Output directory")
	rootCmd.PersistentFlags().StringP("output-package", "p", config.DefaultOutputPackage, "Output package name")
	rootCmd.PersistentFlags().StringP("graphql-endpoint", "e", config.DefaultGraphQLEndpoint, "GraphQL server api endpoint")
	rootCmd.PersistentFlags().StringP("graphql-auth-header", "a", config.DefaultGraphQLAuthHeader, "Header name to be used for authentication")
	rootCmd.PersistentFlags().StringP("graphql-auth-token", "t", config.DefaultGraphQLAuthToken, "Token name to be used for authentication")
	// rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")

	pflags := rootCmd.PersistentFlags()

	utils.Must(viper.BindPFlag("config", pflags.Lookup("config")))
	utils.Must(viper.BindPFlag("output_directory", pflags.Lookup("output-directory")))
	utils.Must(viper.BindPFlag("output_package", pflags.Lookup("output-package")))
	utils.Must(viper.BindPFlag("graphql_endpoint", pflags.Lookup("graphql-endpoint")))
	utils.Must(viper.BindPFlag("graphql_auth_header", pflags.Lookup("graphql-auth-header")))
	utils.Must(viper.BindPFlag("graphql_auth_token", rootCmd.PersistentFlags().Lookup("graphql-auth-token")))

	// bindViperPersistentFlags(config.GetViper(), rootCmd.PersistentFlags(), map[string]string{
	// 	"config": "config",
	// 	"output_directory": "output-directory",
	// 	"output_package": "output-package",
	// 	"graphql_endpoint": "graphql-endpoint",
	// 	"graphql_auth_header": "graphql-auth-header",
	// 	"graphql_auth_token": "graphql-auth-token",
	// })
}

// onCobraInit starts the initialization process
func onCobraInit() {
	initConfig()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	log.SetFlags(log.Flags() | log.LUTC)
	log.SetPrefix(meta.ProjectID + " ")

	config.Initialize(cfgFile)

	if config.Get().Debug {
		config.GetViper().Debug()
	}

	// environment := viper.GetString("environment")
	// log.Printf("RUNTIME: environment=%s", environment)
	// if environment == "development" {
	// 	log.Printf("RUNTIME: config={%v}", viper.AllSettings())
	// }

	// err := godotenv.Load()
	// if err != nil {
	// 	log.Printf("Could not find .env file")
	// }
	//
	// if cfgFile != "" {
	// 	// Use config file from the flag.
	// 	viper.SetConfigFile(cfgFile)
	// } else {
	// 	// Find home directory.
	// 	home, err := homedir.Dir()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		os.Exit(1)
	// 	}
	//
	// 	// Search config in home directory with name ".smart-parking-core" (without extension).
	// 	viper.AddConfigPath(".")
	// 	// viper.AddConfigPath(filepath.FromSlash(fmt.Sprintf("/etc/%s/", meta.ConfigFileName)))
	// 	viper.AddConfigPath(home)
	// 	viper.SetConfigName(meta.ConfigFileNameBase)
	// }
	//
	// viper.SetEnvPrefix("GRAPHQL_TYPEDEF")
	// viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// viper.AllowEmptyEnv(false)
	// viper.AutomaticEnv() // read in environment variables that match
	//
	// // If a config file is found, read it in.
	// if err := viper.ReadInConfig(); err == nil {
	// 	fmt.Println("Using config file:", viper.ConfigFileUsed())
	// }
	//
	// err = viper.Unmarshal(config.Get())
	// if err != nil {
	// 	log.Fatalf("unable to decode configuration, %v", err)
	// }
}
