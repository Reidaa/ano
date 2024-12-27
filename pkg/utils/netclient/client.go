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
	unsuccessfulHTTPResponseThreshold = 300
)

type NetClient struct {
	http http.Client
}

func New() *NetClient {
	return &NetClient{
		http: http.Client{
			Timeout: time.Second * timeout,
		},
	}
}

func (client *NetClient) Get(url string) ([]byte, error) {
	var err error
	var response *http.Response
	var responseData []byte

	utils.Debug.Printf("GET %s", url)

	response, err = http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request to %s failed: %w", url, err)
	}

	defer response.Body.Close()

	if response.StatusCode >= unsuccessfulHTTPResponseThreshold {
		return nil, &UnsuccessfulRequestError{
			StatusCode: response.StatusCode,
			Url:        url,
		}
	}

	responseData, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from %s: %w", url, err)
	}

	return responseData, nil
}
