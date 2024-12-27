package netclient

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/reidaa/ano/pkg/utils"
)

const (
	timeout                           = 10
	unsuccessfulHTTPResponseThreshold = 299
)

type NetClient struct {
	client http.Client
}

func New() *NetClient {
	return &NetClient{
		client: http.Client{
			Timeout: time.Second * timeout,
		},
	}
}

func (nc *NetClient) Get(url string) ([]byte, error) {
	var err error
	var response *http.Response
	var responseData []byte

	utils.Debug.Printf("GET %s", url)

	response, err = nc.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request to %s failed -> %w", url, err)
	}
	utils.Debug.Printf("Received a %d", response.StatusCode)

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotModified {
		return nil, &UnsuccessfulRequestError{
			StatusCode: response.StatusCode,
			Url:        url,
		}
	}

	responseData, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from %s -> %w", url, err)
	}

	return responseData, nil
}
