/*
Copyright © 2021 NAME HERE mrolim@senhasegura.com

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
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	dsmSdk "github.com/senhasegura/dsmcli/sdk/dsm"
	isoSdk "github.com/senhasegura/dsmcli/sdk/iso"
)

var PreparedData map[string]string

var Verbose bool
var ToolName string
var Environment string
var System string
var ApplicationName string

var RunbCmd = &cobra.Command{
	Use:   "runb",
	Short: "Running Belt plugin to insert/get/replace environment variables in most CI/CD process.",
	Long:  `Running Belt plugin to insert/get/replace environment variables in most CI/CD process.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if isDisabled() {
			return errors.Errorf("RUNB_DISABLED is set - Plugin is disabled")
		}

		client, appClient, err := registerApplication()
		if err != nil {
			return err
		}

		envVars := loadEnvVars()
		mapVars := loadMapVars()

		varClient := dsmSdk.NewVariableClient(&client)

		_, err = varClient.Register(envVars, mapVars)
		if err != nil {
			return errors.Errorf("error when posting variables in senhasegura: " + err.Error())
		}

		secrets, err := appClient.GetSecrets()
		if err != nil {
			return err
		}

		err = injectEnvironmentVariables(secrets)
		if err != nil {
			return err
		}

		return deleteCICDVariables()
	},
}

func init() {
	RunbCmd.Flags().BoolVarP(&Verbose, "verbose", "v", false, "Verbose mode")
	RunbCmd.Flags().StringVarP(&ToolName, "tool-name", "t", "linux", "Tool name [github, azure-devops, bamboo, bitbucket, circleci, teamcity, linux]")
	RunbCmd.Flags().StringVarP(&Environment, "environment", "e", "", "Application environment (required)")
	RunbCmd.Flags().StringVarP(&System, "system", "s", "", "Application system (required)")
	RunbCmd.Flags().StringVarP(&ApplicationName, "app-name", "a", "", "Application name (required)")
	RunbCmd.MarkFlagRequired("environment")
	RunbCmd.MarkFlagRequired("system")
	RunbCmd.MarkFlagRequired("app-name")
}

func isDisabled() bool {
	return viper.GetString("RUNB_DISABLED") == "1"
}

func injectEnvironmentVariables(secrets []dsmSdk.Secret) error {
	switch ToolName {
	case "github":
		return injectGithub(secrets)
	case "azure-devops":
		return injectAzureDevops(secrets)
	case "bamboo":
		return injectBamboo(secrets)
	case "bitbucket":
		return injectBitbucket(secrets)
	case "circleci":
		return injectCircleci(secrets)
	case "teamcity":
		return injectTeamcity(secrets)
	case "linux":
		return injectLinux(secrets)

	default:
		return errors.Errorf(
			"tool-name '%s' is invalid, it must be one of the following values: github, azure-devops, bamboo, bitbucket, circleci, teamcity or linux",
			ToolName,
		)
	}
}

func injectGithub(secrets []dsmSdk.Secret) error {
	return inject(secrets, "echo '%s=%s' >> $GITHUB_ENV\n")
}

func injectAzureDevops(secrets []dsmSdk.Secret) error {
	return inject(secrets, "##vso[task.setvariable variable=(%s);issecret=true;](.[%s])\n")
}

func injectBamboo(secrets []dsmSdk.Secret) error {
	return inject(secrets, "(%s)=(.[%s])\n")
}

func injectBitbucket(secrets []dsmSdk.Secret) error {
	return inject(secrets, "export (%s)=\"(.[%s])\"\n")
}

func injectCircleci(secrets []dsmSdk.Secret) error {
	return inject(secrets, "echo '\"'\"'export (%s)=\"(.[%s])\"'\"'\"' >> $BASH_ENV\n")
}

func injectTeamcity(secrets []dsmSdk.Secret) error {
	return inject(secrets, "echo '\"'\"'##teamcity[setParameter name=\"(%s)\" value=\"(.[%s])\"]'\"'\"'\"\n")
}

func injectLinux(secrets []dsmSdk.Secret) error {
	return inject(secrets, "declare -x %s='%s'\n")
}

func inject(secrets []dsmSdk.Secret, format string) error {
	v("Injecting secrets!\n")

	secretsFile := viper.GetString("SENHASEGURA_SECRETS_FILE")

	if secretsFile == "" {
		secretsFile = ".runb.vars"
	}

	file, err := os.OpenFile(secretsFile, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}

	PreparedData = prepareData(secrets)

	if len(PreparedData) == 0 {
		v("No secrets to be injected!\n")
		return nil
	}

	for key, value := range PreparedData {
		v("Injecting secret into %s: %s...", secretsFile, key)
		_, err = file.WriteString(fmt.Sprintf(format, key, value))
		if err != nil {
			return err
		}
		v(" Sucess\n")
	}

	file.Close()
	v("Secrets injected!\n")
	return nil
}

func prepareData(secrets []dsmSdk.Secret) map[string]string {
	preparedData := make(map[string]string)
	for _, secret := range secrets {
		for _, data := range secret.Data {
			for k, v := range data {
				preparedData[k] = v
			}

		}
	}
	return preparedData
}

func deleteCICDVariables() error {
	v("Deleting %s variables...\n", ToolName)

	if len(PreparedData) == 0 {
		v("No variables to be deleted!\n")
		return nil
	}

	switch ToolName {
	case "github":
		err := deleteGitLabVars()
		if err != nil {
			return err
		}

	case "azure-devops":
		v("Is not possible to delete %s variables!\n", ToolName)

	case "bamboo":
		v("Is not possible to delete %s variables!\n", ToolName)

	case "bitbucket":
		v("Is not possible to delete %s variables!\n", ToolName)

	case "circleci":
		v("Is not possible to delete %s variables!\n", ToolName)

	case "teamcity":
		v("Is not possible to delete %s variables!\n", ToolName)

	case "linux":
		v("Is not possible to delete %s variables!\n", ToolName)

	default:
		return errors.Errorf(
			"tool-name '%s' is invalid, it must be one of the following values: github, azure-devops, bamboo, bitbucket, circleci, teamcity or linux",
			ToolName,
		)
	}

	v("Finish\n")

	return nil
}

func deleteGitLabVars() error {
	if !IsSet("GITLAB_ACCESS_TOKEN", "CI_API_V4_URL", "CI_PROJECT_ID") {
		v("Deletion failed\n")
		v("To delete github variables, you need to define the configs GITLAB_ACCESS_TOKEN, CI_API_V4_URL and CI_PROJECT_ID\n")
		return nil
	}

	if len(PreparedData) == 0 {
		v("Deletion failed\n")
		v("Has no credentials to exclude variables on 'github' tool ...\n")
		return nil
	}

	headers := map[string]string{"PRIVATE-TOKEN": viper.GetString("GITLAB_ACCESS_TOKEN")}

	for key := range PreparedData {
		v("Delelting %s variable\n", key)

		resource := fmt.Sprintf(
			"%s/projects/%s/variables/%s",
			viper.GetString("CI_API_V4_URL"),
			viper.GetString("CI_PROJECT_ID"),
			key,
		)

		_, err := isoSdk.DoRequest(
			viper.GetString("GITLAB_ACCESS_TOKEN"),
			resource,
			url.Values{},
			headers,
			http.MethodDelete,
		)

		if err != nil {
			v("Failed trying to delete '%s' variable\n", err.Error())
			continue
		}

		v("Deleted\n")
	}
	return nil
}

func registerApplication() (isoSdk.Client, dsmSdk.ApplicationClient, error) {
	client, _ := isoSdk.NewClient(getConfig())
	appClient := dsmSdk.NewApplicationClient(&client, ApplicationName, Environment, System)

	appResponse, err := appClient.Register()
	if err != nil {
		return client, appClient, err
	}

	client.DefineNewCredentials(appResponse.ID, appResponse.Signature)

	return client, appClient, nil
}

func loadEnvVars() string {
	envVars := strings.Join(os.Environ(), "\n")
	envVars = base64.StdEncoding.EncodeToString([]byte(envVars))
	envVars = replaceSpecials(envVars)
	return envVars
}

func loadMapVars() string {
	if !IsSet("SENHASEGURA_MAPPING_FILE") {
		v("Mapping file not found, proceeding\n")
	} else {
		v("Using mapping file: %s\n", viper.GetString("SENHASEGURA_MAPPING_FILE"))
	}

	content, err := os.ReadFile(viper.GetString("SENHASEGURA_MAPPING_FILE"))
	if err != nil {
		return ""
	}

	mapVars := string(content)
	mapVars = base64.StdEncoding.EncodeToString([]byte(mapVars))
	mapVars = replaceSpecials(mapVars)
	return mapVars
}

func replaceSpecials(value string) string {
	value = strings.Replace(value, "+", "-", -1)
	value = strings.Replace(value, "/", "_", -1)
	value = strings.Replace(value, "=", ",", -1)
	return value
}
