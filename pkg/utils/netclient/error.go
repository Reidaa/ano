package netclient

import "fmt"

type UnsuccessfulRequestError struct {
	Url        string
	StatusCode int
}

func (e *UnsuccessfulRequestError) Error() string {
	return fmt.Sprintf("http response status code from %s is not successful: %d", e.Url, e.StatusCode)
}
