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
package kubernetes

import (
	"fmt"
	"time"

	dsmSdk "github.com/senhasegura/dsmcli/senseg-sdk/dsm"
	isoSdk "github.com/senhasegura/dsmcli/senseg-sdk/iso"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var SidecarCmd = &cobra.Command{
	Use:   "sidecar",
	Short: "Periodically according to the secrets ttl a loop, with the information from /etc/senhasegura/ request the application secrets and save it in the folder /etc/run/secrets/sechasegura/[app_name] to keep updated",
	Long:  `Periodically according to the secrets ttl a loop, with the information from /etc/senhasegura/ request the application secrets and save it in the folder /etc/run/secrets/sechasegura/[app_name] to keep updated`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _ := isoSdk.NewClient(
			viper.GetString("SENHASEGURA_URL"),
			viper.GetString("SENHASEGURA_CLIENT_ID"),
			viper.GetString("SENHASEGURA_CLIENT_SECRET"),
			Verbose);
		appClient := dsmSdk.NewApplicationClient(&client, ApplicationName, Environment, System);
		
		for {
			secrets, err := appClient.GetSecrets()
			if err != nil {
				return err
			}
			secrets.SaveToFile()
			ttl := secrets.GetMinTTL()
			fmt.Printf("Next update in %d seconds...\n", ttl)
			time.Sleep(time.Duration(ttl) * time.Second)
		}
	},
}

func init() {
	SidecarCmd.Flags().BoolVarP(&Verbose, "verbose", "v", false, "verbose mode")
	SidecarCmd.Flags().StringVarP(&Environment, "environment", "e", "", "Application environment (required)")
	SidecarCmd.Flags().StringVarP(&System, "system", "s", "", "Application system (required)")
	SidecarCmd.Flags().StringVar(&ApplicationName, "app-name", "", "Application name (required)")
	SidecarCmd.MarkFlagRequired("environment")
	SidecarCmd.MarkFlagRequired("system")
	SidecarCmd.MarkFlagRequired("app-name")
}
