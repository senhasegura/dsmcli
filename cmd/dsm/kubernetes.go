/*
Copyright © 2021 Matheus Rolim <mrolim@senhasegura.com>

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
package dsm

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"git.mt4.com.br/dev/senhasegura/devops/dsm/cmd/dsm/kubernetes"
)

var KubernetesCmd = &cobra.Command{
	Use:   "kubernetes",
	Short: "Functions to credentials sharing with kubernets",
	Long:  ``,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		err := ValidadeEnvVars()
		if err != nil {
			return err
		}
		return nil;
	},
}

func init() {
	KubernetesCmd.AddCommand(kubernetes.InitContainerCmd)
	KubernetesCmd.AddCommand(kubernetes.SidecarCmd)
}

func ValidadeEnvVars() error {
	if (viper.GetString("SENHASEGURA_SECRETS_FOLDER") == "") {
		os.Setenv("SENHASEGURA_SECRETS_FOLDER", "/var/run/secrets")
	}
	v("Using secret folder: %s", viper.GetString("SENHASEGURA_SECRETS_FOLDER"))
	

	err := os.MkdirAll(viper.GetString("SENHASEGURA_SECRETS_FOLDER"), os.ModePerm)
	if err != nil {
		return fmt.Errorf("invalid config 'SENHASEGURA_SECRETS_FOLDER': %s - %s", viper.GetString("SENHASEGURA_SECRETS_FOLDER"), err.Error())
	}
	
	url := viper.GetString("SENHASEGURA_URL")
	if (url == "") {
		return fmt.Errorf("'SENHASEGURA_URL' must be defined as env var or in config file %s", viper.ConfigFileUsed())
	}

	clientId := viper.GetString("SENHASEGURA_CLIENT_ID")
	if (clientId == "") {
		return fmt.Errorf("'SENHASEGURA_CLIENT_ID' must be defined as env var or in config file %s", viper.ConfigFileUsed())
	}


	clientSecret := viper.GetString("SENHASEGURA_CLIENT_SECRET")
	if (clientSecret == "") {
		return fmt.Errorf("'SENHASEGURA_CLIENT_SECRET' must be defined as env var or in config file %s", viper.ConfigFileUsed())
	}

	return nil;
}
