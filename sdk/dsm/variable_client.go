package dsm

import (
	"net/url"

	sdk "github.com/senhasegura/dsmcli/sdk/iso"
)

type VariableClient struct {
	client *sdk.Client
}

/**
 * Constructor for VariableClient
 */
func NewVariableClient(client *sdk.Client) VariableClient {
	a := VariableClient{
		client: client,
	}

	return a
}

/**
 * Post variables on senhasegura using the cicd endpoint
 * "POST /iso/cicd/variables"
 */
func (a *VariableClient) Register(envVars string, mapVars string) (VariableResponse, error) {
	a.client.V("Posting variables in senhasegura...\n")

	a.client.Authenticate()

	data := url.Values{
		"env": {envVars},
		"map": {mapVars},
	}

	var varResp VariableResponse
	err := a.client.Post("/iso/cicd/variables", data, &varResp)
	if err != nil {
		return VariableResponse{}, err
	}

	a.client.V("Posting variables successfully\n")

	return varResp, nil
}
