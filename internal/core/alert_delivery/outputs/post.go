package outputs

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// post sends a JSON body to an endpoint.
func (client *HTTPWrapper) post(input *PostInput) *AlertDeliveryError {
	payload, err := json.Marshal(input.body)
	if err != nil {
		return &AlertDeliveryError{Message: "json marshal error: " + err.Error(), Permanent: true}
	}

	request, err := http.NewRequest("POST", *input.url, bytes.NewBuffer(payload))
	if err != nil {
		return &AlertDeliveryError{Message: "http request error: " + err.Error(), Permanent: true}
	}

	request.Header.Set("Content-Type", "application/json")

	//Adding dynamic headers
	for key, value := range input.headers {
		request.Header.Set(key, *value)
	}

	response, err := client.httpClient.Do(request)
	if err != nil {
		return &AlertDeliveryError{Message: "network error: " + err.Error()}
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		body, _ := ioutil.ReadAll(response.Body)
		return &AlertDeliveryError{
			Message: "request failed: " + response.Status + ": " + string(body)}
	}

	return nil
}
