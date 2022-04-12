/*
Copyright © 2021 Matheus Rolim mrolim@senhasegrua.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/senhasegura/dsmcli/cmd/dsm"
)

var Config string

var rootCmd = &cobra.Command{
	Use:   "dsm",
	Short: "A command line interface to interact with senhasegura DSM API.",
	Long: `DSM CLI is an unified tool to manage senhasegura services. With this tool, you'll be able to use senhasegura DSM services from the command line and automate them using scripts. 

The main purpose of this tool is to be an agnostic plugin for intercepting environment variables and injecting secrets into systems and CI/CD pipelines.

Using this plugin, DevOps teams have an easy way to centralize application and secret data through senhasegura DSM, providing a secure way for the application to consume sensible variables during the build and deployment steps.`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&Config, "config", "c", "", "Configuration file (default is $HOME/.config.yaml)")

	rootCmd.AddCommand(dsm.RunbCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Read in environment variables
	viper.AutomaticEnv()

	if Config != "" {
		// Use config file from the flag.
		viper.SetConfigFile(Config)
	} else if envConfig := viper.GetString("SENHASEGURA_CONFIG_FILE"); envConfig != "" {
		// Use config file from the environment variable.
		viper.SetConfigFile(envConfig)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".config" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".config")
		viper.SetConfigType("yaml")
		Config = home + "/.config"
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		if strings.Contains(err.Error(), "unmarshal") {
			log.Fatalf(`Invalid yaml syntax on config file '%s'`, viper.ConfigFileUsed())
		}
	}
}
