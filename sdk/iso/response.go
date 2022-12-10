// go:build linux,amd64,!cgo

package iso

import (
	"C"
	"encoding/json"
	"fmt"
)

type ResponseInterface interface {
	Unmarshal(msg []byte) error
	Validate() error
	GetError() string
	GetMessage() string
	GetAccessToken() string
	GetResponse() interface{}
	GetEntity() interface{}
}

type Oauth2Response struct {
	ID          string `json:"id,omitempty"`
	Error       string `json:"error,omitempty"`
	Message     string `json:"message,omitempty"`
	Reason      string `json:"reason,omitempty"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
	Signature   string `json:"signature,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
	AccessToken string `json:"access_token"`
	Response    struct {
		Status    int    `json:"status,omitempty"`
		Message   string `json:"message,omitempty"`
		Error     bool   `json:"error,omitempty"`
		ErrorCode int    `json:"error_code,omitempty"`
	} `json:"response,omitempty"`
}

func (r *Oauth2Response) Unmarshal(msg []byte) error {
	err := json.Unmarshal(msg, &r)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Validate the response of senhasegura server
 */
func (r *Oauth2Response) Validate() error {
	if r.Error != "" {
		return fmt.Errorf(r.Message)
	}

	if r.Response.Error {
		return fmt.Errorf(r.Response.Message)
	}

	return nil
}

func (r *Oauth2Response) GetError() string {
	return r.Error
}

func (r *Oauth2Response) GetMessage() string {
	return r.Message
}

func (r *Oauth2Response) GetAccessToken() string {
	return r.AccessToken
}

func (r *Oauth2Response) GetResponse() interface{} {
	return r.Response
}

func (r *Oauth2Response) GetEntity() interface{} {
	return r.Response
}
