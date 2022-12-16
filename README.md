# senhasegura DSM CLI

DSM CLI is an unified tool to manage senhasegura services. With this tool, you'll be able to use senhasegura DSM services from the command line and automate them using scripts. The main purpose of this tool is to be an agnostic plugin for intercepting environment variables and injecting secrets into systems and CI/CD pipelines.

Using this plugin, DevOps teams have an easy way to centralize application and secret data through senhasegura DSM, providing a secure way for the application to consume sensible variables during the build and deployment steps.

## Using DSM CLI as Running Belt

As for today, senhasegura DSM CLI can only be executed as Running Belt, which reads Secrets from senhasegura DSM module and inject them as environment variables.

As an executable binary, the installation is quite simple. Before deploying the plugin it is important to have a configured application using OAuth 2.0 and an authorization on senhasegura DSM. For more information on how to register applications and authorizations, please check the DSM manual in [Help Center](https://docs.senhasegura.io/?utm_source=Github&utm_medium=Link&utm_campaign=dsm_cli).

The first thing needed is to place the executable into a directory of your environment or CI/CD tool together with a configuration file for authentication on senhasegura DSM. After that, DSM CLI need information from the configured application such as its name, system and environment so it can retrieve the secrets.

For the configuration file, it should be a .yaml file containing the following information from senhasegura DSM:

- **_SENHASEGURA_URL:_** The URL of your senhasegura where DSM is enabled;
- **_SENHASEGURA_CLIENT_ID:_** An authorization Client ID for authentication.
- **_SENHASEGURA_CLIENT_SECRET:_** An authorization Client Secret for authentication.

DSM CLI accepts extra parameters. Here is an axample of a **full .config.yaml** file:

```yaml title=".config.yaml"
# Default properties needed for execution
SENHASEGURA_URL: "<senhasegura URL>"
SENHASEGURA_CLIENT_ID: "<senhasegura Client ID>"
SENHASEGURA_CLIENT_SECRET: "<senhasegura Client Secret>"
SENHASEGURA_MAPPING_FILE: "<Secrets variable name mapping file with path>"
SENHASEGURA_SECRETS_FILE: "<File name with path to inject Secret>"
SENHASEGURA_DISABLE_RUNB: 0

# Properties needed to delete GitLab variables
GITLAB_ACCESS_TOKEN: "<Your GitLab Access Token>"
CI_API_V4_URL: "<Your GitLab API URL as for V4>"
CI_PROJECT_ID: "<Your GitLab Project ID>"
```

> **Using Environment Variables**
> 
> Instead of using a configuration file, DSM CLI can use authentication information through CI/CD environment variables, making the configuration file optional.

To execute the binary you can run the following command line providing the needed information:

```bash
dsm runb \
    --application <application name> \
    --system <system name> \
    --environment <environment name> \
    --config <path to config file>
```
> **Using Environment Variables**
> 
> It is possible to use a **SENHASEGURA_CONFIG_FILE** environment variable to define the configuration file location.

Being agnostic means that it can run in any environment or CI/CD tool, but DSM CLI already comes with some additional configuration allowing you to integrate more seamlessly with your tool.

After executing the plugin with the necessary informations, it will collect all the environment variables running on that pipeline execution and send them to senhasegura DSM.

Then, it will query for all the application secrets registered, injecting them in a file called **.runb.vars** by default or whatever is set on **SENHASEGURA_SECRETS_FILE** if provided, which can be sourcered on the system to update the environment variables with the new values through the command bellow:

```bash
source .runb.vars
```

This way, developers will not have to worry about injecting secrets during pipelines, for example. They can be managed directly via API or through senhasegura DSM interface by any developer or security team member.

> **Security Best Practice**
> 
> Make sure to delete the variables file from the environment to prevent secret leakage.

> **CI/CD Solutions**
> 
> By default DSM CLI can parse the secrets and inject it on tools like GitHub, Azure DevOps, Bamboo, BitBucket, CircleCI, TeamCity and Linux (default option). You can change the default option with the --tool-name argument during its execution.

## Using DSM CLI to Register and Update Secrets

Using DSM CLI also allows developers to create or update secret values directly from the pipeline using a mapping file. This file makes it easy to identify secret variables through their names and automatically register them as secrets on senhasegura DSM.

To do that, the only additional configuration needed is actually to provide the mapping file together with the executable and the configuration file. Here is an example of mapping file's content:

``` json
{
  "access_keys": [
    {
      "name": "ACCESS_KEY_VARIABLES",
      "type": "aws",
      "fields": {
        "access_key_id": "AWS_ACCESS_KEY_ID_VARIABLE",
        "secret_access_key": "AWS_SECRET_ACCESS_KEY_VARIABLE"
      }
    }
  ],
  "credentials": [
    {
      "name": "CREDENCIAL_VARIABLES",
      "fields": {
        "user": "USER_VARIABLE",
        "password": "PASSWORD_VARIABLE",
        "host": "HOST_VARIABLE"
      }
    }
  ],
  "key_value": [
    {
      "name": "GENERIC_VARIABLES",
      "fields": ["KEY_VALUE_VARIABLE"]
    }
  ]
}
```

This file can be broken down in 3 main blocks:

- **_access_keys:_** An array of objects composed by a `name` attribute, `type` and a sub-object `fields`, where this one is composed by an `access_key_id` and `secret_access_key`. These attribute values should be the name of the variable holding the values, so senhasegura DSM will validate if the provided data exists on the Cloud IAM module and if it does it will register it as a secret for that provided authorization.
- **_credentials:_** An array of objects composed by a `name` and a sub-object `fields`, where this one is composed by `user`, `password` and `host`. The values of those attributes should be the name of the variables holding that information so senhasegura DSM will validate if the provided data exists on the PAM module and if it does it will register it as a Secret for that provided authorization.
- **_key_value:_** An array of objects composed by `name` and a sub-array of `fields`, where the values of the array should be the name of the variables to be registered as secrets on senhasegura DSM.

> **Mapping File Name**
> 
> To set the name of the file so DSM CLI can read it use the **SENHASEGURA_MAPPING_FILE** option in the configuration file or set its value as an environment variable pointing to the file full path.

> **Type Values**
> 
> Currently senhasegura DSM only supports access keys through integration with **AWS**, **Azure** or **GCP**, so the **_type_** attribute informed should be one of the supported.
