package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/mitchellh/go-homedir"
	flag "github.com/spf13/pflag"

	"github.com/spf13/viper"

	"github.com/kelseyhightower/envconfig"
	"github.com/shakahl/gqlassist/internal/utils"
	"github.com/shakahl/gqlassist/meta"
)

const (
	EnvPrefix = "GQLASSIST"
)

var (
	config *ConfigSchema
)

func init() {

}

func Get() *ConfigSchema {
	if config == nil {
		panic("config is not yet initialized")
	}
	return config
}

func GetViper() *viper.Viper {
	return viper.GetViper()
}

func BindViperPersistentFlags(f *flag.FlagSet, m map[string]string) {
	for viperKey, flagName := range m {
		utils.Must(viper.BindPFlag(viperKey, f.Lookup(flagName)))
	}
}

func SetDefault(key string, value interface{}) {
	GetViper().SetDefault(key, value)
}

// initConfig reads in config file and ENV variables if set.
func Initialize(cfgFile string) {
	viperCfg := GetViper()

	// Loading dotenv (.env) file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Could not find .env file")
	}

	// Create schema
	config = NewSchema()

	// Loading envconfig
	err = envconfig.Process(EnvPrefix, config)
	if err != nil {
		log.Fatal(err.Error())
	}

	viperCfg.SetEnvPrefix(EnvPrefix)
	viperCfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viperCfg.SetTypeByDefaultValue(true)
	viperCfg.AllowEmptyEnv(false)
	viperCfg.AutomaticEnv() // read in environment variables that match

	// Set viper defaults
	// for k, v := range structToMap(config) {
	// 	utils.Must(GetViper().BindEnv(k))
	// 	viperCfg.SetDefault(k, v)
	// }

	utils.Must(viperCfg.MergeConfigMap(structToMap(config)))

	if cfgFile != "" {
		// Use config file from the flag.
		viperCfg.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".smart-parking-core" (without extension).
		viperCfg.AddConfigPath(".")
		// viper.AddConfigPath(filepath.FromSlash(fmt.Sprintf("/etc/%s/", meta.ConfigFileName)))
		viperCfg.AddConfigPath(home)
		viperCfg.SetConfigName(meta.ConfigFileNameBase)
	}

	// If a config file is found, read it in.
	if err := viperCfg.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	err = viperCfg.Unmarshal(config)
	if err != nil {
		log.Fatalf("unable to decode configuration, %v", err)
	}
}

func structToMap(in interface{}) map[string]interface{} {
	var res map[string]interface{}
	tmp, _ := json.Marshal(in)
	_ = json.Unmarshal(tmp, &res)
	return res
}
