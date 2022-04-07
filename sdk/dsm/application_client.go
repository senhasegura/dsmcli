package dsm

import (
	"fmt"
	"net/url"
	"os"

	sdk "github.com/senhasegura/dsmcli/sdk/iso"
)

type ApplicationClient struct {
	client      *sdk.Client
	name        string
	system      string
	environment string
}

/**
 * Constructor for ApplicationClient
 */
func NewApplicationClient(client *sdk.Client, name string, environment string, system string) ApplicationClient {
	if string(name) == "" {
		fmt.Println("Error: Application name must be defined")
		os.Exit(1)
	}

	if string(environment) == "" {
		fmt.Println("Error: Environment must be defined")
		os.Exit(1)
	}

	if string(system) == "" {
		fmt.Println("Error: System must be defined")
		os.Exit(1)
	}

	a := ApplicationClient{
		name:        name,
		environment: environment,
		system:      system,
		client:      client,
	}

	return a
}

/**
 * Register a new authorization for this application on senhasegura using the iso endpoint
 * "POST /iso/dapp/Application"
 */
func (a *ApplicationClient) Register() (ApplicationResponse, error) {
	a.client.V("Registering Application on DevSecOps\n")

	a.client.Authenticate()

	data := url.Values{
		"application": {a.name},
		"environment": {a.environment},
		"system":      {a.system},
	}

	var appResp ApplicationResponse
	err := a.client.Post("/iso/dapp/Application", data, &appResp)
	if err != nil {
		return ApplicationResponse{}, err
	}
	a.client.V("Application register success\n")

	return appResp, nil
}

/**
 * Authenticate on senhasegura with credentials provided from application api response
 */
func (a *ApplicationClient) DefineCredentialsByApplication(application ApplicationResponse) error {
	return a.client.DefineNewCredentials(application.ID, application.Signature)
}

/**
 * Makes requests for /iso/dapp/Application
 * to get Application
 */
func (a ApplicationClient) GetApplication() (ApplicationResponse, error) {
	a.client.Authenticate()

	var appResp ApplicationResponse
	err := a.client.Get("/iso/dapp/Application", url.Values{}, &appResp)
	if err != nil {
		return ApplicationResponse{}, err
	}

	return appResp, nil
}

/**
 * Makes requests for /iso/dapp/Application
 * to get secrets of Application
 */
func (a ApplicationClient) GetSecrets() (secrets, error) {
	a.client.V("Finding secrets from application\n")

	app, err := a.GetApplication()
	if err != nil {
		return nil, err
	}

	return app.Application.Secrets, nil
}

func (a ApplicationClient) GetClient() *sdk.Client {
	return a.client
}