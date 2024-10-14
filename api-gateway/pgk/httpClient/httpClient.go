package httpClient

import (
	"io"
	"net/http"
	"strings"
)

func PostRequest(url string, body io.Reader) ([]byte, error) {
	resp, err := http.Post(url, "application/json", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// TODO: FIX THESE FUNCTIONS
func PostRequestLogin(url string, body io.Reader) (*http.Response, error) {
	resp, err := http.Post(url, "application/json", body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func PutRequest(url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func DeleteRequest(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodDelete, url, strings.NewReader(""))
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
