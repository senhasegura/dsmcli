
# senhasegura DSM CLI
​
DSM CLI is an unified tool to manage senhasegura services. With this tool, you'll be able to use senhasegura DSM services from the command line and automate them using scripts. The main purpose of this tool is to be an agnostic plugin for intercepting environment variables and injecting secrets into systems and CI/CD pipelines.
​
Using this plugin, DevOps teams have an easy way to centralize application and secret data through senhasegura DSM, providing a secure way for the application to consume sensible variables during the build and deployment steps.
​
## Using DSM CLI as Running Belt
​
The CLI can be executed in two main modes: RunB and Kubernetes. In this section we are going to explain the usage through the RunB option.
​
As an executable binary, its installation is quite simple. Before deploying the plugin it is important to have a configured application using OAuth 2.0 and an authorization on senhasegura DSM. For more information on how to register applications and authorizations, please check the [Applications](../applications.md) and [Authorizations](../authorizations.md) guide.
​
The first thing needed is to the executable into a directory of your environment or CI/CD tool together with a configuration file for authentication on senhasegura DSM. After that, DSM CLI need information from the configured application such as its name, system and environment so it can retrieve the secrets.
​
For the configuration file, it should be a .yaml file containing the following information from senhasegura DSM:
​
-   **_SENHASEGURA_URL:_** The URL of your senhasegura where DSM is enabled;
-   **_SENHASEGURA_CLIENT_ID:_** An authorization Client ID for authentication.
-   **_SENHASEGURA_CLIENT_SECRET:_** An authorization Client Secret for authentication.
​
Example of a **.config.yaml** file:
​
``` yaml title=".config.yaml"
SENHASEGURA_URL: "<senhasegura URL>"
SENHASEGURA_CLIENT_ID: "<senhasegura Client ID>"
SENHASEGURA_CLIENT_SECRET: "<senhasegura Client Secret>"
```
​
:::tip Using Environment Variables
Instead of using a configuration file, DSM CLI can use authentication information through CI/CD environment variables, making the configuration file optional.
:::
​
To execute the binary you can run the following command line providing the needed information:
​
``` bash
dsm runb \
    --app-name <application name> \
    --system <system name> \
    --environment <environment name> \
    --config <path to config file>
```
​
Being agnostic means that it can run in any environment or CI/CD tool, but DSM CLI already comes with some additional configuration allowing you to integrate more seamlessly with your tool.
​
After executing the plugin with the necessary informations, it will collect all the environment variables running on that pipeline execution and send them to senhasegura DSM.
​
Then, it will query for all the application secrets registered, injecting them in a file called **.runb.vars**, which can be sourcered on the system to update the environment variables with the new values through the command bellow:
​
```bash
source .runb.vars
```
​
This way, developers will not have to worry about injecting secrets during pipelines, for example. They can be managed directly via API or through senhasegura DSM interface by any developer or security team member.
​
:::caution Security Best Practice
Make sure to delete the variables file from the environment to prevent secret leakage.
:::
​
:::tip CI/CD Solutions
By default DSM CLI can parse the secrets and inject it on tools like GitHub, Azure DevOps, Bamboo, BitBucket, CircleCI, TeamCity and Linux (default option). You can change the default option with the --tool-name argument during its execution.
:::
​
## Using DSM CLI as Kubernetes Sidecar
​
The DSM CLI also have an option to run similarly to the Kubernetes Sidecar plugin, where it fetches the secrets from senhasegura DSM and inject them as files in a user defined folder (usually /var/run/secrets/senhasegura).
​
This method also allows you to run it as sidecar or init-container. As a sidecar, DSM CLI will run continuously, updating the secrets every 120 seconds, while as init-container it will run only once during its execution.
​
You can use the following commands to execute it as sidecar:
​
``` bash
dsm kubernetes sidecar \
    --app-name <application name>
    --system <system name> \
    --environment <environment name> \
    --config <path to config file>
```
​
Or the following to execute it as init-container:
​
``` bash
dsm kubernetes init-container \
    --app-name <application name>
    --system <system name> \
    --environment <environment name> \
    --config <path to config file>
```
​
:::caution Inject Secrets on the Default Folder
The default folder for secret injection is `/var/run/secrets/senhasegura/<secret identifier>`. To inject secrets in the default folder make sure you run it with administrative privileges.
:::
​
:::tip Change Secrets Default Folder
Additionally, in the config file you can define the **SENHASEGURA_SECRETS_FOLDER** with a path where you want the plugin to make the secret data available, as in the example:
​
``` yaml title=".config.yaml"
SENHASEGURA_URL: "<senhasegura URL>"
SENHASEGURA_CLIENT_ID: "<senhasegura Client ID>"
SENHASEGURA_CLIENT_SECRET: "<senhasegura Client Secret>"
SENHASEGURA_SECRETS_FOLDER: "<senhasegura Secrets Destination Folder>"
```
:::
​
## Using DSM CLI to Register and Update Secrets
​
Using DSM CLI also allows developers to create or update secret values directly from the pipeline using a mapping file called **senhasegura-mapping.json**. This file makes it easy to identify secret variables through their names and automatically register them as secrets on senhasegura DSM.
​
To do that, the only additional configuration needed is actually to provide the mapping file together with the executable and the configuration file. Here is an example of mapping file's content:
​
``` json title="senhasegura-mapping.json"
{
  "access_keys": [
    {
      "name": "AWS_VARIABLES",
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
      "name": "SECRET_VARIABLES",
      "fields": ["KEY_VALUE_VARIABLE"]
    }
  ]
}
```
​
This file can be broken down in 3 main blocks:
​
-   **_access_keys:_** An array of objects composed by a `name` attribute, `type` and a sub-object `fields`, where this one is composed by an `access_key_id` and `secret_access_key`. These attribute values should be the name of the variable holding the values, so senhasegura DSM will validate if the provided data exists on the Cloud IAM module and if it does it will register it as a secret for that provided authorization.
-   **_credentials:_** An array of objects composed by a `name` and a sub-object `fields`, where this one is composed by `user`, `password` and `host`. The values of those attributes should be the name of the variables holding that information so senhasegura DSM will validate if the provided data exists on the PAM module and if it does it will register it as a Secret for that provided authorization.
-   **_key_value:_** An array of objects composed by `name` and a sub-array of `fields`, where the values of the array should be the name of the variables to be registered as secrets on senhasegura DSM.
​
:::caution Mapping File Name
This file should be named exactly as senhasegura-mapping.json and should be on the same directory level as the executable.
:::
​
:::caution Type Values
Currently senhasegura DSM only supports access keys through integration with **AWS**, **Azure** or **GCP**, so the **_type_** attribute informed should be one of the supported.
:::
