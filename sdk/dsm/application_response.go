package dsm

import (
	// "C"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/spf13/viper"
)

type ApplicationResponse struct {
	ID          string      `json:"id"`
	Signature   string      `json:"signature"`
	Error       string      `json:"error"`
	Message     string      `json:"message"`
	Application Application `json:"application"`
	Response    struct {
		Status    int    `json:"status"`
		Message   string `json:"message"`
		Error     bool   `json:"error"`
		ErrorCode int    `json:"error_code"`
	} `json:"response"`
}

type Application struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	System      string   `json:"system"`
	Environment string   `json:"Environment"`
	Secrets     secrets  `json:"secrets"`
}

type secrets []Secret

type Secret struct {
	SecretID       string              `json:"secret_id"`
	SecretName     string              `json:"secret_name"`
	Identity       string              `json:"identity"`
	Version        string              `json:"version"`
	ExpirationDate string              `json:"expiration_date"`
	Engine         string              `json:"engine"`
	Data           []map[string]string `json:"data"`
}

func (r *ApplicationResponse) Unmarshal(msg []byte) error {
	err := json.Unmarshal(msg, r)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Validate the response of senhasegura server
 */
func (r *ApplicationResponse) Validate() error {
	if r.Error != "" {
		return fmt.Errorf(r.Message)
	}

	if r.Response.Error {
		return fmt.Errorf(r.Response.Message)
	}

	return nil
}

func (r *ApplicationResponse) GetError() string {
	return r.Error
}

func (r *ApplicationResponse) GetMessage() string {
	return r.Message
}

func (r *ApplicationResponse) GetAccessToken() string {
	return r.Message
}

func (r *ApplicationResponse) GetResponse() interface{} {
	return r.Response
}

func (r *ApplicationResponse) GetEntity() interface{} {
	return r.Response
}

/**
 * Save the current client info to
 * files at /var/run/secrets/senhasegura/iso
 */
func (a *ApplicationResponse) SaveToFile() error {
	fmt.Println("Adding credentials to system...")

	secretDirectory := viper.GetString("SENHASEGURA_SECRETS_FOLDER") + "/senhasegura/iso"
	err := os.MkdirAll(secretDirectory, os.ModePerm)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(
		secretDirectory+"/SENHASEGURA_URL",
		[]byte(viper.GetString("SENHASEGURA_URL")),
		os.ModePerm,
	)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(
		secretDirectory+"/SENHASEGURA_CLIENT_ID",
		[]byte(a.ID),
		os.ModePerm,
	)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(
		secretDirectory+"/SENHASEGURA_CLIENT_SECRET",
		[]byte(a.Signature),
		os.ModePerm,
	)
	if err != nil {
		return err
	}

	fmt.Println("Complete.")
	return nil
}

/**
 * Save multiple credentials on folders
 * "/var/run/secrets/senhasegura/[application_name]"
 */
func (s secrets) SaveToFile() error {
	fmt.Println("Adding credentials to system...")

	secretDirectory := viper.GetString("SENHASEGURA_SECRETS_FOLDER") + "/senhasegura"

	err := os.MkdirAll(secretDirectory, os.ModePerm)
	if err != nil {
		return err
	}

	err = RemoveContents(secretDirectory)
	if err != nil {
		return err
	}

	// Read dirs and files inside /var/run/secrets/senhasegura
	dir, err := ioutil.ReadDir(secretDirectory)
	if err != nil {
		return err
	}

	// Remove all dirs and files inside /var/run/secrets/senhasegura
	for _, d := range dir {
		err := os.RemoveAll(path.Join([]string{secretDirectory, d.Name()}...))

		if err != nil {
			return err
		}
	}

	err = os.MkdirAll(secretDirectory, os.ModePerm)
	if err != nil {
		return err
	}

	for _, secret := range s {
		folder := fmt.Sprintf(
			"%s/%s",
			secretDirectory,
			secret.Identity,
		)

		err = os.MkdirAll(folder, os.ModePerm)
		if err != nil {
			return err
		}

		secret.saveToFile(folder)
	}

	fmt.Println("Complete.")
	return nil
}

func (s secrets) GetMinTTL() int64 {
	var ttl int64 = 120

	for _, secret := range s {
		ttl = secret.getMinTTL(ttl)
	}

	return ttl
}

/**
 * Save the data of secret on folders
 * "/etc/run/secrets/senhasegura/[application_name]/[secret_identifier]"
 */
func (s Secret) saveToFile(folder string) error {
	for _, data := range s.Data {
		for filename, content := range data {
			err := ioutil.WriteFile(
				fmt.Sprintf("%s/%s", folder, filename),
				[]byte(content),
				os.ModePerm,
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s Secret) getMinTTL(current int64) int64 {
	newTTL := current

	for _, data := range s.Data {
		for key, value := range data {
			if key == "TTL" && value != "" {
				ttl, err := strconv.ParseInt(value, 10, 64)
				if err != nil && ttl > 10 && ttl < newTTL {
					newTTL = ttl
				}
			}
		}
	}

	return newTTL
}

// Remove o conteudo de um diretorio
func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
