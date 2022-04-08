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

	dsmSdk "github.com/senhasegura/dsmcli/sdk/dsm"
	isoSdk "github.com/senhasegura/dsmcli/sdk/iso"
	"github.com/spf13/cobra"
)

var Verbose bool
var ExecutionType string
var Environment string
var System string
var ApplicationName string

var InitContainerCmd = &cobra.Command{
	Use:   "init-container",
	Short: "Collect the application's secrets in senhasegura and insert into the application container in the directory **/var/run/secrets/senhasegura**",
	Long: `
    In a traditional container structure, it is very common for an application to be replicated for use in different
    parts of the world. As a result, it is necessary that the secrets used by her also be replicated, generating a high risk
    attack point. Since they are replicas, they all tend to use the same secret to use external applications and a leak of
    this information would cause damage that cannot be traced.
    
    The sidecar senhasegura comes to mediate this communication and dynamically provision different secrets for each of these
    replicas and keep all these environments safe and remote. In this way, the leakage of one of these secrets would not have
    a direct impact on other applications, it would allow quick tracking and an efficient counter attack when updating this secret.
    
    Sidecars and init containers are used to mediate communication between the container and the safe for the search for secrets and
    storage in a file with its contents on the server
    
    Init container is a special type of **auxiliary container** for applications running inside Kubernetes pods, normally performing
    bootstrap tasks like application settings and **secrets query**
    
    The init container is executed during the creation of a pod, **and only allows applications to start after its successful execution**
    
    For more information see the [official Kubernetes documentation on init containers] (https://kubernetes.io/docs/concepts/workloads/pods/init-containers/)
    
    How does it work?
    There are three types of execution for the module, which is defined by the --type flag:
    
    iso:
    Register a new authorization in senhasegura using the iso endpoint  POST /iso/dapp/application", with the information from
    /etc/senhasegura/ and save the credentials of authorizantion into path /etc/run/secrets/iso
    
    inject_template:
    Register a new authorization in ISO using the iso endpoint "POST /iso/dapp/application", using the information from
    /etc/senhasegura/ with the acquired authorization, request the application secrets and save it in the folder /etc/run/secrets/sechasegura/[app_name]
    
    inject:
    With the information from /etc/senhasegura/ request the application secrets and save it in the folder /etc/run/secrets/sechasegura/[app_name]
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, _ := isoSdk.NewClient(getConfig())
		appClient := dsmSdk.NewApplicationClient(&client, ApplicationName, Environment, System)

		switch ExecutionType {
		case "iso":
			return iso(appClient)

		case "inject":
			return inject(appClient)

		case "inject-template":
			return injectTemplate(appClient)

		default:
			return fmt.Errorf(
				"execution type '%s' is invalid, it must be one of the following values: iso, inject or inject-template",
				ExecutionType,
			)
		}
	},
}

func init() {
	InitContainerCmd.Flags().BoolVarP(&Verbose, "verbose", "v", false, "verbose mode")
	InitContainerCmd.Flags().StringVarP(&ExecutionType, "type", "t", "iso", "Execution type [iso, inject-template, inject]")
	InitContainerCmd.Flags().StringVarP(&Environment, "environment", "e", "", "Application environment (required)")
	InitContainerCmd.Flags().StringVarP(&System, "system", "s", "", "Application system (required)")
	InitContainerCmd.Flags().StringVarP(&ApplicationName, "app-name", "a", "", "Application name (required)")
	InitContainerCmd.MarkFlagRequired("environment")
	InitContainerCmd.MarkFlagRequired("system")
	InitContainerCmd.MarkFlagRequired("app-name")
}

/**
 * Register a new authorization in ISO using the iso endpoint
 * POST /iso/dapp/application", with the information from
 * /etc/senhasegura/ and save credentials of authorizantion
 * into path /var/run/secrets/iso
 */
func iso(appClient dsmSdk.ApplicationClient) error {
	// Register a new authorization in ISO
	app, err := appClient.Register()
	if err != nil {
		return err
	}
	// Save the authorizantion credentials
	err = app.SaveToFile()
	if err != nil {
		return err
	}

	return nil
}

func inject(appClient dsmSdk.ApplicationClient) error {
	// Register a new authorization in ISO
	app, err := appClient.Register()
	if err != nil {
		return err
	}

	// Uses the acquired authorization to next requests
	err = appClient.DefineCredentialsByApplication(app)
	if err != nil {
		return err
	}

	// Finds application of the acquired authorization
	secrets, err := appClient.GetSecrets()
	if err != nil {
		return err
	}

	// Save credentials of application
	err = secrets.SaveToFile()
	if err != nil {
		return err
	}

	return nil
}

func injectTemplate(appClient dsmSdk.ApplicationClient) error {
	// Finds application from te configured credentials
	secrets, err := appClient.GetSecrets()
	if err != nil {
		return err
	}

	// Save credentials
	err = secrets.SaveToFile()
	if err != nil {
		return err
	}

	return nil
}
