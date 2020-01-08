package outputs

/**
 * Copyright 2020 Panther Labs Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import (
	"bytes"
	"io/ioutil"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

// post sends a JSON body to an endpoint.
func (client *HTTPWrapper) post(input *PostInput) *AlertDeliveryError {
	payload, err := jsoniter.Marshal(input.body)
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
