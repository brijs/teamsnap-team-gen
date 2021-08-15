package ts

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	cj "github.com/brijs/teamsnap-team-gen/collectionjson"
)

type Client struct {
	accessToken string
	baseURL     string
	HTTPClient  *http.Client
}

func NewClient(accessToken string) *Client {
	return &Client{
		accessToken: accessToken,
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		baseURL: "https://api.teamsnap.com/v3/",
	}
}

// TODO
type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type successResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

// Content-type and body should be already added to req
func (c *Client) sendRequest(req *http.Request) (resCj cj.CollectionJsonType, err error) {
	// common headers
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return resCj, err
	}

	defer res.Body.Close()

	// handle response
	if res.StatusCode != http.StatusOK {
		return resCj, fmt.Errorf("Received Bad HTTP status code: %d\n", res.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return resCj, err
	}

	resCj, err = cj.ReadCollectionJson(bodyBytes)
	if err != nil {
		return resCj, err
	}

	// TODO: check resCj.Collection.Error

	return resCj, err
}
