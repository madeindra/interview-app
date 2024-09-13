package openai

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func getResponseBody(resp *http.Response) (io.ReadCloser, error) {
	if resp == nil || resp.Body == nil {
		log.Default().Println("response is nil")

		return nil, fmt.Errorf("response is nil")
	}

	if resp.StatusCode != http.StatusOK {
		log.Default().Printf("unexpected status code: %d", resp.StatusCode)
		log.Default().Printf("response: %v", resp)

		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}

func unmarshalJSONResponse(resp *http.Response, v interface{}) error {
	respBody, err := getResponseBody(resp)
	if err != nil {
		log.Default().Printf("error getting response body: %v", err)

		return err
	}
	if respBody == nil {
		log.Default().Println("response body is nil")

		return fmt.Errorf("response body is nil")
	}
	defer respBody.Close()

	return json.NewDecoder(respBody).Decode(&v)
}
