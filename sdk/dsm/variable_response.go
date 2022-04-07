package dsm

import (
	"encoding/json"
	"fmt"
)

type VariableResponse struct {
	Error        string   `json:"error"`
	Message      string   `json:"message"`
	Response     struct {
		Status    int    `json:"status"`
		Message   string `json:"message"`
		Error     bool   `json:"error"`
		ErrorCode int    `json:"error_code"`
	} `json:"response"`
}


func (r *VariableResponse) Unmarshal(msg []byte) error {
	err := json.Unmarshal(msg, r)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Validate the response of senhasegura server
 */
func (r *VariableResponse) Validate() error {
	if r.Error != "" {
		return fmt.Errorf(r.Message)
	}

	if r.Response.Error {
		return fmt.Errorf(r.Response.Message)
	}

	return nil
}

func (r *VariableResponse) GetError() string {
	return r.Error
}

func (r *VariableResponse) GetMessage() string {
	return r.Message
}

func (r *VariableResponse) GetAccessToken() string {
	return r.Message
}

func (r *VariableResponse) GetResponse() interface{} {
	return r.Response
}

func (r *VariableResponse) GetEntity() interface{} {
	return r.Response
}