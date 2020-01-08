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
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockHTTPClient struct {
	HTTPiface
	statusCode   int
	requestError bool
	requestBody  string // Request body is saved here for tests to verify
}

var requestEndpoint = "https://runpanther.io"

func (m *mockHTTPClient) Do(request *http.Request) (*http.Response, error) {
	if m.requestError {
		return nil, errors.New("endpoint unreachable")
	}
	requestBytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}
	m.requestBody = string(requestBytes)

	responseBody := ioutil.NopCloser(bytes.NewReader([]byte("response")))
	return &http.Response{Body: responseBody, StatusCode: m.statusCode}, nil
}

func TestPostInvalidJSON(t *testing.T) {
	body := map[string]interface{}{"func": TestPostInvalidJSON}
	postInput := &PostInput{
		url:  &requestEndpoint,
		body: body,
	}
	c := &HTTPWrapper{httpClient: &mockHTTPClient{}}
	assert.NotNil(t, c.post(postInput))
}

func TestPostErrorSubmittingRequest(t *testing.T) {
	c := &HTTPWrapper{httpClient: &mockHTTPClient{requestError: true}}
	postInput := &PostInput{
		url:  &requestEndpoint,
		body: map[string]interface{}{"abc": 123},
	}
	assert.NotNil(t, c.post(postInput))
}

func TestPostNotOk(t *testing.T) {
	c := &HTTPWrapper{httpClient: &mockHTTPClient{statusCode: http.StatusBadRequest}}
	postInput := &PostInput{
		url:  &requestEndpoint,
		body: map[string]interface{}{"abc": 123},
	}
	assert.NotNil(t, c.post(postInput))
}

func TestPostOk(t *testing.T) {
	c := &HTTPWrapper{httpClient: &mockHTTPClient{statusCode: http.StatusOK}}
	postInput := &PostInput{
		url:  &requestEndpoint,
		body: map[string]interface{}{"abc": 123},
	}
	assert.Nil(t, c.post(postInput))
}

func TestPostCreated(t *testing.T) {
	c := &HTTPWrapper{httpClient: &mockHTTPClient{statusCode: http.StatusCreated}}
	postInput := &PostInput{
		url:  &requestEndpoint,
		body: map[string]interface{}{"abc": 123},
	}
	assert.Nil(t, c.post(postInput))
}
