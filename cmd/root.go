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
var Verbose bool

var rootCmd = &cobra.Command{
	Use:   "dsm",
	Short: "A helper to interact with senhasegura appications.",
	Long: `The senhasegura dsm is a unified tool to management senhasegura devops services.
With this tool, you'll be able to use senhasegura dsm's services from the command line and automate
them using scripts.

The senhasegura CLI offers features for dvops including init_container, sidecar and runb support to
help you strengthen the security of your shared credentials with containers and ephemeral machines.`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose mode")
	rootCmd.PersistentFlags().StringVarP(&Config, "config", "c", "", "config file (default is $HOME/dsm/.config.yaml)")

	rootCmd.AddCommand(dsm.KubernetesCmd)
	rootCmd.AddCommand(dsm.RunbCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match

	if Config != "" {
		// Use config file from the flag.
		viper.SetConfigFile(Config)
	} else if envConfig := viper.GetString("SENHASEGURA_CONFIG_FILE"); envConfig != "" {
		// Use config file from the environment.
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
