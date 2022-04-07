package iso

import (
	"C"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Client struct {
	Verbose      bool
	url          string
	clientID     string
	clientSecret string
	accessToken  string
}

/**
 * Contructor for client object
 */
func NewClient(senhaseguraUrl string, clientID string, clientSecret string, verbose bool) (Client, error) {
	url := strings.Trim(string(senhaseguraUrl), "\n ")
	if url == "" {
		return Client{}, fmt.Errorf("senhaseguraUrl cannot be null")
	}

	clientID = strings.Trim(string(clientID), "\n ")
	if clientID == "" {
		return Client{}, fmt.Errorf("clientID cannot be null")
	}

	clientSecret = strings.Trim(string(clientSecret), "\n ")
	if clientSecret == "" {
		return Client{}, fmt.Errorf("clientSecret cannot be null")
	}

	c := Client{
		Verbose:      verbose,
		url:          url,
		clientID:     clientID,
		clientSecret: clientSecret,
	}

	return c, nil
}

func (c *Client) DefineNewCredentials(clientID string, clientSecret string) error {
	clientID = strings.Trim(string(clientID), "\n ")
	if clientID == "" {
		return fmt.Errorf("clientID cannot be null")
	}

	clientSecret = strings.Trim(string(clientSecret), "\n ")
	if clientSecret == "" {
		return fmt.Errorf("clientSecret cannot be null")
	}

	c.clientID = clientID
	c.clientSecret = clientSecret
	return nil
}

/**
 * Performs authetication on senhasegura DevSecOps API
 */
func (c *Client) Authenticate() {
	c.V("Trying to authenticate on senhasegura DevSecOps API\n")

	resource := "/iso/oauth2/token"

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", c.clientID)
	data.Set("client_secret", c.clientSecret)

	var oauth2Resp Oauth2Response
	err := c.Post(resource, data, &oauth2Resp)
	if err != nil {
		log.Fatal("Error trying to authenticate: " + err.Error())
	}

	c.accessToken = "Bearer " + oauth2Resp.GetAccessToken()

	c.V("Authenticated successfully\n")
}

func (c *Client) V(format string, a ...interface{}) {
	if c.Verbose {
		fmt.Printf(format, a...)
	}
}

/**
 * Performs a post request on senhasegura server
 */
func (c Client) Post(resource string, data url.Values, responseObj ResponseInterface) error {
	return c.call(http.MethodPost, resource, data, responseObj)
}

/**
 * Performs a get request on senhasegura server
 */
func (c Client) Get(resource string, data url.Values, responseObj ResponseInterface) error {
	return c.call(http.MethodGet, resource, data, responseObj)
}

/**
 * Performs a request on senhasegura server
 */
func (c Client) call(method string, resource string, data url.Values, responseObj ResponseInterface) error {
	headers := make(map[string]string)
	if c.accessToken != "" {
		headers["Authorization"] = c.accessToken
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	headers["Content-Length"] = strconv.Itoa(len(data.Encode()))

	responseData, err := DoRequest(c.url, resource, data, headers, method)
	if err != nil {
		return err
	}

	err = responseObj.Unmarshal(responseData)
	if err != nil {
		return err
	}

	err = responseObj.Validate()
	if err != nil {
		return err
	}

	return nil
}


func DoRequest(host string, resource string, data url.Values, headers map[string]string, method string) ([]byte, error) {
	u, err := url.ParseRequestURI(host)
	if err != nil {
		return nil, err
	}
	u.Path = resource
	urlStr := u.String()

	tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }	
	httpClient := &http.Client{Transport: tr}
	r, err := http.NewRequest(method, urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		r.Header.Add(k, v)
	}

	resp, err := httpClient.Do(r)
	if err != nil {
		return nil, err
	}

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}
