package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io.ml.challenges/io.ml.challenges.books-wishlist/pkg/mlbwlist/models/config"
)

const (
	appName            = "MercadoLibre - Books Wishlist App"
	appCLIName         = "mlbwlist"
	configFilepathFlag = "configs"
	secretFilepathFlag = "secrets"
)

var cmd = &cobra.Command{
	Use:     appCLIName,
	Long:    fmt.Sprintf("Daemon for %s [%s]", appName, appCLIName),
	PreRunE: VerifyRequiredFlags,
	Run:     StartServer,
}

func main() {

	cmd.Flags().String(configFilepathFlag, fmt.Sprintf("/etc/%s/configs.yml", appCLIName), "configurations to start service and use connectors")
	_ = viper.BindPFlag(configFilepathFlag, cmd.Flags().Lookup(configFilepathFlag))

	cmd.Flags().String(secretFilepathFlag, fmt.Sprintf("/etc/%s/secrets.yml", appCLIName), "configurations to start service and use connectors")
	_ = viper.BindPFlag(secretFilepathFlag, cmd.Flags().Lookup(secretFilepathFlag))

	viper.AutomaticEnv()
	err := cmd.Execute()
	if err != nil {
		fmt.Println("ERROR: cannot run server with the provided configuration")
		os.Exit(1)
	}
}

func ReadConfig() config.Config {

	configfilepath := viper.GetString(configFilepathFlag)
	configfilename := getFilenameFromFilepath(configfilepath)
	configfiledir := getFiledireFromFilepath(configfilepath)

	fmt.Printf("reading configuration file from %s\n", configfilepath)

	viper.SetConfigName(configfilename)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configfiledir)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("ERROR: cannot read configuration file from %s\n", configfilepath)
		os.Exit(2)
	}

	secretfilepath := viper.GetString(secretFilepathFlag)
	secretfilename := getFilenameFromFilepath(secretfilepath)
	secretfiledir := getFiledireFromFilepath(secretfilepath)

	fmt.Printf("reading secrets file from %s\n", secretfilepath)

	viper.SetConfigName(secretfilename)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(secretfiledir)
	err = viper.MergeInConfig()
	if err != nil {
		fmt.Printf("ERROR: cannot read secrets file from %s\n", secretfilepath)
		os.Exit(2)
	}

	appconfig := config.Config{}
	_ = viper.Unmarshal(&appconfig)

	appconfig.UUID = uuid.New().String()

	return appconfig
}

func VerifyRequiredFlags(cmd *cobra.Command, _ []string) error {

	var requiredFlagsNotFound []string
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		required := len(flag.Annotations[cobra.BashCompOneRequiredFlag]) > 0 && flag.Annotations[cobra.BashCompOneRequiredFlag][0] == "true"
		_, envVarDefined := os.LookupEnv(flag.Name)
		if required && (!flag.Changed && !envVarDefined) {
			requiredFlagsNotFound = append(requiredFlagsNotFound, flag.Name)
		}
	})
	if len(requiredFlagsNotFound) > 0 {
		return fmt.Errorf("The following flags are required and were not present: %s ", strings.Join(requiredFlagsNotFound, ", "))
	}
	return nil

}

func StartServer(_ *cobra.Command, _ []string) {

	config := ReadConfig()
	_, defined := os.LookupEnv("BOOKS_WISHLIST_DEBUG_MODE")
	if defined {
		configJSON, _ := json.MarshalIndent(config, "", "    ")
		fmt.Println(string(configJSON))
	}
}

func getFilenameFromFilepath(filepath string) string {

	paths := strings.Split(filepath, "/")
	return paths[len(paths)-1]
}

func getFiledireFromFilepath(filepath string) string {

	paths := strings.Split(filepath, "/")
	return strings.Join(paths[:len(paths)-1], "/")
}
